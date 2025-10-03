Raptor: Go Migration Tool
Raptor is a powerful and flexible database migration tool for Go, inspired by the intuitive approach of Laravel's migrations. It allows you to define database schema changes in clean, reversible classes and manage their execution using a simple CLI.

✨ Features
Plug-and-Play Database Support: Database drivers are implemented via the raptor.Schema interface, allowing easy integration with any SQL database (SQLite, MySQL, PostgreSQL, etc.) by implementing the necessary SQL dialect logic.

Up and Down: Every migration is reversible using Up() and Down() methods.

Batch Tracking: Migrations are executed in batches, enabling simple rollback of the last set of changes.

Status Reporting: Easily view which migrations are applied and which are pending.

🚀 Installation
To use the raptor package in your Go project:

go get [github.com/mrhoseah/raptor](https://github.com/mrhoseah/raptor)

🛠 Quick Start
1. Define Migrations
Migrations are Go structs that implement the raptor.Migration interface.

// migrations/001_create_users_table.go
package migrations

import "[github.com/mrhoseah/raptor](https://github.com/mrhoseah/raptor)"

type CreateUsersTable struct{}

func (m *CreateUsersTable) Name() string {
    return "001_create_users_table"
}

// Up runs when migrating forward
func (m *CreateUsersTable) Up(s raptor.Schema) error {
    return s.CreateTable("users", []string{"id", "email", "password", "created_at"})
}

// Down runs when rolling back
func (m *CreateUsersTable) Down(s raptor.Schema) error {
    return s.DropTable("users")
}

2. Implement a Database Schema (Plug-and-Play)
You must create a struct that implements the raptor.Schema interface for your chosen database. This is where you connect to the database and run the specific SQL commands.

// In a separate file (e.g., drivers/postgres_schema.go)
package drivers

import "[github.com/mrhoseah/raptor](https://github.com/mrhoseah/raptor)"
import "database/sql"

type PostgresSchema struct {
    DB *sql.DB
}

// Ensure PostgresSchema implements raptor.Schema
var _ raptor.Schema = (*PostgresSchema)(nil) 

func (s *PostgresSchema) CreateTable(name string, columns []string) error {
    // In a real application, construct and execute PostgreSQL-specific CREATE TABLE SQL here.
    return nil 
}

func (s *PostgresSchema) DropTable(name string) error {
    // In a real application, construct and execute PostgreSQL-specific DROP TABLE SQL here.
    return nil
}

3. Run Migrations
In your main.go file, register the migrations and pass your concrete Schema implementation to the raptor.NewMigrator constructor.

// main.go
package main

import (
	"log"
    // Assuming you implemented drivers/postgres_schema.go
    // "your_project/drivers" 
    
    // ... imports ...
)

func main() {
    // 1. Register migrations
	allMigrations := GetMigrations()
    
    // 2. Initialize your database connection (Example)
    // dbConn, _ := sql.Open("postgres", "connection_string")
    // postgresSchema := &drivers.PostgresSchema{DB: dbConn}
    
    // 3. Instantiate the Migrator (PLUG-AND-PLAY)
    // mgr := raptor.NewMigrator(allMigrations, postgresSchema) 

    // For simulation/testing:
    mgr := raptor.NewSimulatedMigrator(allMigrations) 

    // ... CLI logic remains the same ...
	
	command := os.Args[1]
	
	switch command {
	case "migrate":
		// Executes pending migrations
	case "rollback":
		// Reverts the last batch
	case "status":
		// Shows migration history
	}
}

📜 CLI Commands
Command

Description

go run main.go migrate

Runs all pending migrations and records them as a new batch.

go run main.go rollback

Reverses the most recent batch of migrations by calling the Down() method.

go run main.go status

Displays a list of all migrations, showing which are APPLIED and which are PENDING.

