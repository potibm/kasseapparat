package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/potibm/kasseapparat/internal/app/initializer"
	"github.com/potibm/kasseapparat/internal/app/models"
	"github.com/potibm/kasseapparat/internal/app/repository"
	"github.com/potibm/kasseapparat/internal/app/utils"
)

var (
	purgeDB           bool // Flag to indicate whether to purge the database
	seedData          bool // Flag to indicate whether to seed initial data
	createUserName    string
	createUserEmail   string
	userImportCsvFile string
	createUserIsAdmin bool // Flag to indicate whether to create a user with admin rights
)

func init() {
	// Register command-line flags
	flag.BoolVar(&purgeDB, "purge", false, "Purge the database before initializing")
	flag.BoolVar(&seedData, "seed", false, "Seed initial data to the database")
	flag.StringVar(&userImportCsvFile, "import-users", "", "CSV file to import users from")
	flag.StringVar(&createUserName, "create-user", "", "Create a user with the given username")
	flag.StringVar(&createUserEmail, "create-user-email", "", "Email for the user to create")
	flag.BoolVar(&createUserIsAdmin, "create-user-admin", false, "Create a user with admin rights")
	flag.Parse()
}

func main() {
	initializer.InitializeVersion()
	initializer.OutputVersion()

	db := utils.ConnectToDatabase()

	if userImportCsvFile != "" {
		log.Println("Importing users from CSV file...")
		importUsers(userImportCsvFile)
		return
	}

	if createUserName != "" {
		if createUserEmail == "" {
			fmt.Println("Error: --createUserEmail is required when creating a user.")
			os.Exit(1)
		}
		log.Println("Creating user", createUserName)
		createUser(createUserName, createUserEmail, createUserIsAdmin)
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
		utils.SeedDatabase(db)
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
		log.Fatal(err)
	}

	for _, record := range records {
		log.Println("Creating user", record[0])

		createUserAndSendToken(record[0], record[1], record[2] == "true")
	}
}

func createUser(username string, email string, isAdmin bool) {
	createUserAndSendToken(username, email, isAdmin)

	log.Println("User created")
}

func createUserAndSendToken(username string, email string, isAdmin bool) {
	repo := repository.NewRepository()
	mailer := initializer.InitializeMailer()

	user := models.User{
		Username: username,
		Email:    email,
		Password: "",
		Admin:    isAdmin,
	}
	user.GenerateRandomPassword()

	validity := 24 * time.Hour
	user.GenerateChangePasswordToken(&validity)

	user, err := repo.CreateUser(user)
	if err != nil {
		log.Println("Error creating a user", err)
		return
	}

	err = mailer.SendNewUserTokenMail(user.Email, user.ID, user.Username, *user.ChangePasswordToken)
	if err != nil {
		log.Println("Error sending email", err)
	}
}
