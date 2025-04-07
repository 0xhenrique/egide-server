package migration

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	log.Printf("Running migrations from directory: %s", migrationsDir)
	
	// Get all .sql files from the migrations directory
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	
	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)
	
	// Execute each migration file in a transaction
	for _, fileName := range migrationFiles {
		log.Printf("Applying migration: %s", fileName)
		
		filePath := filepath.Join(migrationsDir, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			log.Printf("Migration failed: %v", err)
			// @TODO: this is a hack. Should strive for keeping track of the failing migrations
			continue
		}
		
		if err := tx.Commit(); err != nil {
			return err
		}
		
		log.Printf("Successfully applied migration: %s", fileName)
	}
	
	return nil
}
