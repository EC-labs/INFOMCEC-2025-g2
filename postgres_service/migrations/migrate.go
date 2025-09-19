package migrations

import (
	"fmt"
	"log"
	"postgres_service/config"
	"postgres_service/models"
)

// RunMigrations performs auto-migration for all models
func RunMigrations() error {
	db := config.GetDB()
	if db == nil {
		return fmt.Errorf("database connection not established")
	}

	log.Println("Running database migrations...")

	// Auto-migrate all models
	allModels := models.GetAllModels()
	for _, model := range allModels {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
		log.Printf("Migrated model: %T", model)
	}

	log.Println("Database migrations completed successfully!")
	return nil
}