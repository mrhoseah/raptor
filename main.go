package main

import (
	"fmt"
	"log"
	"os"

	// 1. Import the core raptor package
	raptor "github.com/mrhoseah/raptor/core"
	// 2. Import the package containing your migration files
	"github.com/mrhoseah/raptor/migrations"
)

// GetMigrations registers all available migration structs.
func GetMigrations() []raptor.Migration {
	return []raptor.Migration{
		// Add all your concrete migration structs here
		&migrations.CreateUsersTable{},
	}
}

func main() {
	// --- Migrator Initialization (The Plug-and-Play Step) ---

	allMigrations := GetMigrations()

	// A. Simulated Usage (for demonstration and testing)
	mgr := raptor.NewSimulatedMigrator(allMigrations)

	// B. Real Database Usage (for SQLite, MySQL, Postgres)
	/*
		// Example for SQLite (requires a connection and a concrete Schema implementation)
		dbConn, err := sql.Open("sqlite3", "./migrations.db")
		if err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}
		// Assuming you have a concrete 'SQLiteSchema' struct that implements raptor.Schema
		sqliteSchema := drivers.NewSQLiteSchema(dbConn)
		mgr := raptor.NewMigrator(allMigrations, sqliteSchema)
	*/

	// --- Command Execution ---

	// Determine command from arguments (simulate artisan/cli tool)
	command := "status"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	log.Printf("Running command: %s", command)

	var err error
	switch command {
	case "migrate":
		err = mgr.Migrate()
	case "rollback":
		err = mgr.Rollback()
	case "status":
		mgr.Status()
	default:
		fmt.Printf("Unknown command: %s. Available commands: migrate, rollback, status\n", command)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
