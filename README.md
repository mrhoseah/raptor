# 🦖 Raptor: Go Database Migration Tool

**Raptor** is a powerful and flexible database migration tool for Go, inspired by Laravel's elegant migration system. It enables developers to manage database schema changes using clean, reversible Go structs and a simple CLI interface.

Whether you're working with PostgreSQL, MySQL, SQLite, or any other SQL database, Raptor makes schema evolution seamless and maintainable.

---

## 👤 Background & Motivation

Built in **2023**, Raptor was born out of necessity. Coming from a Laravel background, I was used to its intuitive migration system and wanted something similar in Go—clean, reversible, and easy to manage. After months of development and refinement, I’ve finalized the tool and made it publicly available for others who want Laravel-style migrations in their Go projects.

---

## ✨ Key Features

- **Plug-and-Play Database Support**  
  Easily integrate with any SQL database by implementing the `raptor.Schema` interface. Supports PostgreSQL, MySQL, SQLite, and more.

- **Reversible Migrations**  
  Each migration includes `Up()` and `Down()` methods for forward and backward compatibility.

- **Batch Execution & Rollback**  
  Migrations run in batches, allowing you to rollback the last set of changes with a single command.

- **Migration Status Reporting**  
  View which migrations are applied and which are pending with a simple status command.

---

## 🚀 Installation

To install Raptor in your Go project:

```bash
go get github.com/mrhoseah/raptor
```

---

## 🛠 Quick Start Guide

### 1. Define a Migration

Create a Go struct that implements the `raptor.Migration` interface:

```go
// migrations/001_create_users_table.go
package migrations

import raptor "github.com/mrhoseah/raptor/core"

type CreateUsersTable struct{}

func (m *CreateUsersTable) Name() string {
	return "001_create_users_table"
}

func (m *CreateUsersTable) Up(s raptor.Schema) error {
	return s.CreateTable("users", []string{"id", "email", "password", "created_at"})
}

func (m *CreateUsersTable) Down(s raptor.Schema) error {
	return s.DropTable("users")
}
```

---

### 2. Implement a Database Schema

Create a struct that implements the `raptor.Schema` interface for your database:

```go
// drivers/postgres_schema.go
package drivers

import (
    "database/sql"
    "github.com/mrhoseah/raptor"
)

type PostgresSchema struct {
    DB *sql.DB
}

var _ raptor.Schema = (*PostgresSchema)(nil)

func (s *PostgresSchema) CreateTable(name string, columns []string) error {
    // PostgreSQL-specific CREATE TABLE logic
    return nil
}

func (s *PostgresSchema) DropTable(name string) error {
    // PostgreSQL-specific DROP TABLE logic
    return nil
}
```

---

### 3. Run Migrations via CLI

Register migrations and execute them using Raptor's CLI:

```go
// main.go
package main

import (
    "os"
    "github.com/mrhoseah/raptor"
    // "your_project/drivers"
)

func main() {
    allMigrations := GetMigrations()

    // Example: Initialize your DB connection
    // dbConn, _ := sql.Open("postgres", "your_connection_string")
    // schema := &drivers.PostgresSchema{DB: dbConn}

    // Use simulated migrator for testing
    mgr := raptor.NewSimulatedMigrator(allMigrations)

    switch os.Args[1] {
    case "migrate":
        mgr.Migrate()
    case "rollback":
        mgr.Rollback()
    case "status":
        mgr.Status()
    }
}
```

---

## 📜 CLI Commands

| Command                   | Description                                                  |
|---------------------------|--------------------------------------------------------------|
| `go run main.go migrate`  | Runs all pending migrations and records them as a new batch. |
| `go run main.go rollback` | Reverts the most recent batch of migrations.                 |
| `go run main.go status`   | Displays applied and pending migrations.                     |

---

## 🔗 Useful Links

- 📦 [Raptor on GitHub](https://github.com/mrhoseah/raptor)
- 📚 Go Documentation: [golang.org/doc](https://golang.org/doc)
- 🐘 PostgreSQL: [postgresql.org](https://www.postgresql.org)
- 🐬 MySQL: [mysql.com](https://www.mysql.com)

---

## 💡 Why Choose Raptor?

Raptor is ideal for Go developers who want a clean, testable, and extensible migration system. With its Laravel-inspired design and plug-and-play architecture, it fits naturally into modern Go projects and CI/CD pipelines.

---
