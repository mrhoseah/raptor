package raptor

import (
	"fmt"
	"sort"
	"strings"
)

// Migration interface defines the contract for all database migrations.
type Migration interface {
	Name() string
	Up(s Schema) error
	Down(s Schema) error
}

// Schema defines the Domain Specific Language (DSL) for modifying the database structure.
// This interface is the "plug" point for different database drivers (SQLite, MySQL, Postgres).
type Schema interface {
	CreateTable(name string, columns []string) error
	DropTable(name string) error
	// Future methods like AddColumn, RenameTable, etc., would go here.
}

// --- Simulated Components (for Testing/Demonstration) ---

// SimulatedSchema implements the Schema interface without a real database connection.
type SimulatedSchema struct{}

// CreateTable simulates table creation.
func (s *SimulatedSchema) CreateTable(name string, columns []string) error {
	fmt.Printf("  [Simulated Schema] Creating table '%s' with columns: %s\n", name, strings.Join(columns, ", "))
	return nil
}

// DropTable simulates table dropping.
func (s *SimulatedSchema) DropTable(name string) error {
	fmt.Printf("  [Simulated Schema] Dropping table '%s'\n", name)
	return nil
}

// SimulatedDB replaces a real database connection for tracking history.
type SimulatedDB struct {
	History   map[string]int // MigrationName -> Batch Number
	NextBatch int
}

// NewSimulatedDB creates a new SimulatedDB instance.
func NewSimulatedDB() *SimulatedDB {
	return &SimulatedDB{
		History:   make(map[string]int),
		NextBatch: 1,
	}
}

// GetRanMigrations returns a list of migration names that have already been run.
func (s *SimulatedDB) GetRanMigrations() []string {
	names := make([]string, 0, len(s.History))
	for name := range s.History {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetLastBatch returns the highest batch number recorded.
func (s *SimulatedDB) GetLastBatch() int {
	maxBatch := 0
	for _, batch := range s.History {
		if batch > maxBatch {
			maxBatch = batch
		}
	}
	return maxBatch
}

// --- Migrator Core ---

// Migrator handles the execution and tracking of migrations.
type Migrator struct {
	migrations []Migration
	schema     Schema
	db         *SimulatedDB // In a real app, this would be a DB interface
}

// NewMigrator creates a new Migrator instance, accepting a concrete Schema implementation.
// This is the primary, plug-and-play constructor.
func NewMigrator(migrations []Migration, schema Schema) *Migrator {
	// Sort the migrations to ensure chronological order based on name (e.g., timestamp prefixes).
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	return &Migrator{
		migrations: migrations,
		schema:     schema,
		db:         NewSimulatedDB(), // Assuming we still use the simulated history tracking for simplicity
	}
}

// NewSimulatedMigrator creates a Migrator using the SimulatedSchema.
// This function is now properly capitalized and exported.
func NewSimulatedMigrator(migrations []Migration) *Migrator {
	return NewMigrator(migrations, &SimulatedSchema{})
}

// Migrate applies all pending migrations.
func (m *Migrator) Migrate() error {
	ranMigrations := m.db.GetRanMigrations()
	ranSet := make(map[string]bool)
	for _, name := range ranMigrations {
		ranSet[name] = true
	}

	pending := []Migration{}
	for _, mig := range m.migrations {
		if !ranSet[mig.Name()] {
			pending = append(pending, mig)
		}
	}

	if len(pending) == 0 {
		fmt.Println("Database is up to date. Nothing to migrate.")
		return nil
	}

	batch := m.db.NextBatch
	fmt.Printf("\n--- Running Migrations (Batch %d) ---\n", batch)

	for _, mig := range pending {
		fmt.Printf("-> Applying %s...\n", mig.Name())
		if err := mig.Up(m.schema); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", mig.Name(), err)
		}
		m.db.History[mig.Name()] = batch
	}

	m.db.NextBatch++
	fmt.Println("--- Migration Complete ---")
	return nil
}

// Rollback reverses the migrations from the most recent batch.
func (m *Migrator) Rollback() error {
	batchToRollback := m.db.GetLastBatch()

	if batchToRollback == 0 {
		fmt.Println("\nNo migrations have been run. Nothing to rollback.")
		return nil
	}

	migrationsToRevert := []Migration{}
	migrationsInBatch := make(map[string]bool)

	// 1. Identify migrations belonging to the batch to rollback
	for name, batch := range m.db.History {
		if batch == batchToRollback {
			migrationsInBatch[name] = true
		}
	}

	// 2. Map names back to the Migration instances and sort them in reverse chronological order
	// We iterate the full list backwards to get a list of migrations in the last batch,
	// ordered correctly for reversal.
	for i := len(m.migrations) - 1; i >= 0; i-- {
		mig := m.migrations[i]
		if migrationsInBatch[mig.Name()] {
			migrationsToRevert = append(migrationsToRevert, mig)
		}
	}

	if len(migrationsToRevert) == 0 {
		fmt.Printf("\nBatch %d has no recorded migrations. Nothing to rollback.\n", batchToRollback)
		return nil
	}

	fmt.Printf("\n--- Rolling back Batch %d ---\n", batchToRollback)

	for _, mig := range migrationsToRevert {
		fmt.Printf("<- Reverting %s...\n", mig.Name())
		if err := mig.Down(m.schema); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", mig.Name(), err)
		}
		// Remove from history
		delete(m.db.History, mig.Name())
	}

	// Decrement the next batch number so the next run uses the number of the rolled back batch
	m.db.NextBatch = batchToRollback

	fmt.Println("--- Rollback Complete ---")
	return nil
}

// Status prints the current migration status.
func (m *Migrator) Status() {
	ranMigrations := m.db.GetRanMigrations()
	ranSet := make(map[string]bool)
	for _, name := range ranMigrations {
		ranSet[name] = true
	}

	fmt.Println("\n--- Migration Status ---")
	fmt.Printf("%-30s %s\n", "Migration Name", "Status (Batch)")
	fmt.Println(strings.Repeat("-", 45))

	pendingCount := 0
	for _, mig := range m.migrations {
		name := mig.Name()
		if batch, ok := m.db.History[name]; ok {
			fmt.Printf("%-30s %s (%d)\n", name, "Ran", batch)
		} else {
			fmt.Printf("%-30s %s\n", name, "Pending")
			pendingCount++
		}
	}

	if pendingCount == 0 && len(ranMigrations) > 0 {
		fmt.Println("\nDatabase is fully up to date.")
	} else if len(ranMigrations) == 0 {
		fmt.Println("\nNo migrations have been run yet.")
	}
	fmt.Println(strings.Repeat("-", 45))
}
