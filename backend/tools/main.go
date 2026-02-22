package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"time"

	"github.com/potibm/kasseapparat/internal/app/config"
	"github.com/potibm/kasseapparat/internal/app/exitcode"
	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/mailer"
	"github.com/potibm/kasseapparat/internal/app/models"
	sqliteRepo "github.com/potibm/kasseapparat/internal/app/repository/sqlite"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

var (
	purgeDB           bool // Flag to indicate whether to purge the database
	seedData          bool // Flag to indicate whether to seed initial data
	seedDataWithTest  bool // Flag to indicate whether to seed initial data with test data
	createUserName    string
	createUserEmail   string
	userImportCsvFile string
	createUserIsAdmin bool // Flag to indicate whether to create a user with admin rights
)

var (
	Repo   *sqliteRepo.Repository
	Mailer mailer.Mailer
)

var (
	version = "0.0.0"
)

func main() {
	// Register command-line flags
	flag.BoolVar(&purgeDB, "purge", false, "Purge the database before initializing")
	flag.BoolVar(&seedData, "seed", false, "Seed initial data to the database")
	flag.BoolVar(&seedDataWithTest, "seed-with-test", false, "Seed initial data plus testdata to the database")
	flag.StringVar(&userImportCsvFile, "import-users", "", "CSV file to import users from")
	flag.StringVar(&createUserName, "create-user", "", "Create a user with the given username")
	flag.StringVar(&createUserEmail, "create-user-email", "", "Email for the user to create")
	flag.BoolVar(&createUserIsAdmin, "create-user-admin", false, "Create a user with admin rights")
	flag.Parse()

	log := initializer.InitTxtLogger("debug")

	cfg := config.Load()
	cfg.SetVersion(version)
	cfg.OutputVersion()

	db := utils.ConnectToDatabase()

	Repo = sqliteRepo.NewRepository(db, int32(cfg.FormatConfig.FractionDigitsMax))
	Mailer = initializer.InitializeMailer(cfg.MailerConfig)

	if userImportCsvFile != "" {
		log.Info("Importing users from CSV file...")
		importUsers(userImportCsvFile)

		return
	}

	if createUserName != "" {
		if createUserEmail == "" {
			fmt.Println("Error: -create-user-email is required when creating a user.")
			os.Exit(1)
		}

		if !validEmail(createUserEmail) {
			fmt.Println("Error: Invalid email format")
			os.Exit(1)
		}

		log.Info("Creating user", "username", createUserName)

		err := createUser(createUserName, createUserEmail, createUserIsAdmin)
		if err != nil {
			log.Error("Failed to create user", "error", err)

			return
		}

		log.Info("User created")

		return
	}

	// Purge the database if requested
	if purgeDB {
		log.Info("Purging database...")
		utils.PurgeDatabase(db)
	}

	log.Info("Starting database migration...")
	utils.MigrateDatabase(db)

	// Seed initial data if requested
	if seedData {
		log.Info("Start seeding DB...")
		utils.SeedDatabase(db, false)
	} else if seedDataWithTest {
		log.Info("Start seeding DB with test data...")
		utils.SeedDatabase(db, true)
	}
}

func importUsers(filename string) {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		slog.Error("Failed to open CSV file", "error", err)
		os.Exit(int(exitcode.Usage))
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		slog.Error("Failed to read CSV file", "error", err)

		return
	}

	for _, record := range records {
		const expectedFields = 3
		if len(record) != expectedFields {
			slog.Warn("Skipping malformed record", "record", record)

			continue
		}

		slog.Info("Creating user", "username", record[0])

		isAdmin := false
		if record[2] == "true" || record[2] == "false" {
			isAdmin = record[2] == "true"
		} else {
			slog.Warn("Invalid admin value", "value", record[2], "username", record[0])
		}

		if !validEmail(record[1]) {
			slog.Warn("Invalid email format for user", "username", record[0])

			continue
		}

		err := createUser(record[0], record[1], isAdmin)
		if err != nil {
			slog.Error("Failed to create user", "error", err)

			continue
		}
	}
}

func createUser(username string, email string, isAdmin bool) error {
	user := models.User{
		Username: username,
		Email:    email,
		Password: "",
		Admin:    isAdmin,
	}
	user.GenerateRandomPassword()

	const tokenValidityHours = 24

	validity := tokenValidityHours * time.Hour
	user.GenerateChangePasswordToken(&validity)

	user, err := Repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = Mailer.SendNewUserTokenMail(user.Email, user.ID, user.Username, *user.ChangePasswordToken)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}
