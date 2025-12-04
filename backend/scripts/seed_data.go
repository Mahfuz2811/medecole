package main

import (
	"log"
	"github.com/Mahfuz2811/medecole/backend/internal/config"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/scripts/seeders"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run auto migration
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize and run seeder
	seeder := seeders.NewDatabaseSeeder(db)
	seeder.Run()
}
