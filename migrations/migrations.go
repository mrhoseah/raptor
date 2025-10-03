package migrations

import (
	"log"

	raptor "github.com/mrhoseah/raptor/core"
)

// CreateUsersTable MUST start with a capital letter (C) to be exported
// and accessible by the 'main' package.
type CreateUsersTable struct{}

// Name returns the unique identifier for this migration.
func (m *CreateUsersTable) Name() string {
	return "001_create_users_table"
}

// Up defines the steps to apply this migration.
func (m *CreateUsersTable) Up(s raptor.Schema) error {
	log.Println("Running UP for 001_create_users_table")
	// Using generic column names for demonstration.
	return s.CreateTable("users", []string{"id", "email", "password", "created_at"})
}

// Down defines the steps to reverse this migration (rollback).
func (m *CreateUsersTable) Down(s raptor.Schema) error {
	log.Println("Running DOWN for 001_create_users_table")
	return s.DropTable("users")
}
