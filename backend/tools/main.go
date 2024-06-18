package main

import (
	"flag"
	"log"

	"github.com/potibm/kasseapparat/internal/app/utils"
	"github.com/potibm/kasseapparat/internal/app/initializer"
)

var (
	purgeDB  bool // Flag to indicate whether to purge the database
	seedData bool // Flag to indicate whether to seed initial data
)

func init() {
	// Register command-line flags
	flag.BoolVar(&purgeDB, "purge", false, "Purge the database before initializing")
	flag.BoolVar(&seedData, "seed", false, "Seed initial data to the database")
	flag.Parse()
}

func main() {
	initializer.InitializeVersion()
	initializer.OutputVersion()

	db := utils.ConnectToDatabase()

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
