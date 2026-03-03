//go:build mage

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Default = Build

// Initialize creates a new SQLite database
func Init() error {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database initialized successfully: words.db")
	return nil
}

// Migrate runs all migration files in order
func Migrate() error {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	migrationsDir := "db/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		filepath := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		fmt.Printf("Applied migration: %s\n", filename)
	}

	fmt.Println("All migrations completed successfully")
	return nil
}

// Seed imports JSON seed files into the database
func Seed() error {
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	seedsDir := "db/seeds"
	files, err := ioutil.ReadDir(seedsDir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filename := file.Name()
			groupName := strings.TrimSuffix(filename, ".json")

			filepath := filepath.Join(seedsDir, filename)
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				return fmt.Errorf("failed to read seed file %s: %w", filename, err)
			}

			// Insert group first
			var groupID int64
			result, err := db.Exec("INSERT OR IGNORE INTO groups (name) VALUES (?)", groupName)
			if err != nil {
				return fmt.Errorf("failed to insert group %s: %w", groupName, err)
			}

			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get group ID: %w", err)
			}

			if id == 0 {
				// Group already exists, get its ID
				err = db.QueryRow("SELECT id FROM groups WHERE name = ?", groupName).Scan(&groupID)
				if err != nil {
					return fmt.Errorf("failed to get group ID for %s: %w", groupName, err)
				}
			} else {
				groupID = id
			}

			// Parse JSON and insert words
			var words []map[string]interface{}
			if err := json.Unmarshal(content, &words); err != nil {
				return fmt.Errorf("failed to parse JSON in %s: %w", filename, err)
			}

			for _, wordData := range words {
				bengali := wordData["bengali"].(string)
				partOfSpeech := wordData["part_of_speech"].(string)
				english := wordData["english"].(string)

				// Insert word
				var wordID int64
				result, err := db.Exec(`
					INSERT OR IGNORE INTO words (bengali, parts, english) 
					VALUES (?, ?, ?)
				`, bengali, partOfSpeech, english)

				if err != nil {
					return fmt.Errorf("failed to insert word %s: %w", bengali, err)
				}

				id, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("failed to get word ID: %w", err)
				}

				if id == 0 {
					// Word already exists, get its ID
					err = db.QueryRow("SELECT id FROM words WHERE bengali = ?", bengali).Scan(&wordID)
					if err != nil {
						return fmt.Errorf("failed to get word ID for %s: %w", bengali, err)
					}
				} else {
					wordID = id
				}

				// Link word to group
				_, err = db.Exec(`
					INSERT OR IGNORE INTO words_groups (word_id, group_id) 
					VALUES (?, ?)
				`, wordID, groupID)

				if err != nil {
					return fmt.Errorf("failed to link word %s to group %s: %w", bengali, groupName, err)
				}
			}

			fmt.Printf("Seeded group: %s (ID: %d) with %d words\n", groupName, groupID, len(words))
		}
	}

	fmt.Println("Database seeding completed successfully")
	return nil
}

// Build runs the full setup: init, migrate, and seed
func Build() error {
	fmt.Println("Starting database build process...")

	if err := Init(); err != nil {
		return err
	}

	if err := Migrate(); err != nil {
		return err
	}

	if err := Seed(); err != nil {
		return err
	}

	fmt.Println("Database build completed successfully!")
	return nil
}

// TestDB creates a test database for testing
func TestDB() error {
	fmt.Println("Creating test database...")

	// Use test database path from environment or default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "words.test.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open test database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping test database: %w", err)
	}

	fmt.Printf("Test database initialized successfully: %s\n", dbPath)

	// Run migrations on test database
	if err := Migrate(); err != nil {
		return fmt.Errorf("test database migration failed: %w", err)
	}

	// Seed test database
	if err := Seed(); err != nil {
		return fmt.Errorf("test database seeding failed: %w", err)
	}

	fmt.Println("Test database setup completed successfully!")
	return nil
}

// TestDBSQL creates a test database using SQL file (faster approach)
func TestDBSQL() error {
	fmt.Println("Creating test database with SQL...")

	// Use test database path from environment or default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "words.test.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open test database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping test database: %w", err)
	}

	fmt.Printf("Test database initialized successfully: %s\n", dbPath)

	// Run migrations directly on test database
	migrationsDir := "db/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		migrationPath := filepath.Join(migrationsDir, file)
		content, err := ioutil.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		fmt.Printf("Applied migration: %s\n", file)
	}

	fmt.Println("All migrations completed successfully")

	// Load and execute SQL test data
	sqlFile := "db/seeds/test_data.sql"
	content, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read SQL seed file %s: %w", sqlFile, err)
	}

	// Split SQL content into individual statements
	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement '%s': %w", stmt, err)
		}
	}

	fmt.Println("Test database setup with SQL completed successfully!")
	return nil
}

// Clean removes the database file
func Clean() error {
	if err := os.Remove("words.db"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Database file does not exist")
			return nil
		}
		return fmt.Errorf("failed to remove database file: %w", err)
	}

	fmt.Println("Database file removed successfully")
	return nil
}
