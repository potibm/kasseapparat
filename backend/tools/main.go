package main

import (
	"encoding/csv"
	"flag"
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
	userImportCsvFile string
)

func init() {
	// Register command-line flags
	flag.BoolVar(&purgeDB, "purge", false, "Purge the database before initializing")
	flag.BoolVar(&seedData, "seed", false, "Seed initial data to the database")
	flag.StringVar(&userImportCsvFile, "import-users", "", "CSV file to import users from")
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
	repo := repository.NewRepository()
	mailer := initializer.InitializeMailer()

	file, err := os.Open(filename)
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

		user := models.User{
			Username: record[0],
			Email:    record[1],
			Password: "",
			Admin:    record[2] == "true",
		}
		user.GenerateRandomPassword()
		validity := 24 * time.Hour
		user.GenerateChangePasswordToken(&validity)

		user, err := repo.CreateUser(user)
		if err != nil {
			log.Println("Error creating a user", err)
			continue
		}

		err = mailer.SendNewUserTokenMail(user.Email, user.ID, user.Username, *user.ChangePasswordToken)
		if err != nil {
			log.Println("Error sending email", err)
		}
	}
}
