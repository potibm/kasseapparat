package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/potibm/kasseapparat/internal/app/config"
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

	cfg := config.Load()
	cfg.OutputVersion()

	db := utils.ConnectToDatabase()

	Repo = sqliteRepo.NewRepository(db, int32(cfg.FormatConfig.FractionDigitsMax))
	Mailer = initializer.InitializeMailer(cfg.MailerConfig)

	if userImportCsvFile != "" {
		log.Println("Importing users from CSV file...")
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

		log.Println("Creating user", createUserName)

		err := createUser(createUserName, createUserEmail, createUserIsAdmin)
		if err != nil {
			log.Println("Failed to create user:", err)

			return
		}

		log.Println("User created")

		return
	}

	// Purge the database if requested
	if purgeDB {
		log.Println("Purging database...")
		utils.PurgeDatabase(db)
	}

	log.Println("Starting database migration...")
	utils.MigrateDatabase(db)

	// Seed initial data if requested
	if seedData {
		log.Println("Start seeding DB...")
		utils.SeedDatabase(db, false)
	} else if seedDataWithTest {
		log.Println("Start seeding DB with test data...")
		utils.SeedDatabase(db, true)
	}
}

func importUsers(filename string) {
	file, err := os.Open(filename) // #nosec G304
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file: %v", err)

		return
	}

	for _, record := range records {
		const expectedFields = 3
		if len(record) != expectedFields {
			log.Printf("Skipping malformed record: %v", record)

			continue
		}

		log.Println("Creating user", record[0])

		isAdmin := false
		if record[2] == "true" || record[2] == "false" {
			isAdmin = record[2] == "true"
		} else {
			log.Printf("Invalid admin value '%s' for user %s, defaulting to false", record[2], record[0])
		}

		if !validEmail(record[1]) {
			log.Printf("Invalid email format for user %s, skipping", record[0])

			continue
		}

		err := createUser(record[0], record[1], isAdmin)
		if err != nil {
			log.Println("Failed to create user:", err)

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
