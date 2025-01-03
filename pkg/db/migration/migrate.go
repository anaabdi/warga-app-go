package migration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 2
	_defaultTimeout  = time.Second
)

func Exec(migrationType, databaseURL, migrationEnv string) {
	migrationDir := fmt.Sprintf("file://%s", resolvePath("scripts/db/migrations"))
	seedDir := fmt.Sprintf("file://%s", resolvePath(fmt.Sprintf("scripts/db/seeds/%s", migrationEnv)))

	if migrationType == "up" {
		exec(migrationType, databaseURL, migrationDir, "Migrate")
		exec(migrationType, databaseURL+"&x-migrations-table=schema_seeder", seedDir, "Seeding")
	} else {
		exec(migrationType, databaseURL+"&x-migrations-table=schema_seeder", seedDir, "Seeding")
		exec(migrationType, databaseURL, migrationDir, "Migrate")
	}
}

func exec(migrationType, databaseURL, migrationDir, tag string) {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New(migrationDir, databaseURL)
		if err == nil {
			break
		}
		log.Printf(tag+": postgres is trying to connect %s, attempts left: %d", databaseURL, attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf(tag+": postgres connect error: %s", err)
	}

	if migrationType == "up" {
		err = m.Up()
	} else {
		err = m.Down()
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf(tag+": %v error: %s", migrationType, err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Print(tag + ": no change")
		return
	}

	log.Printf(tag+": %v success", migrationType)
}

func resolvePath(relativePath string) string {
	baseDir, err := os.Getwd() // Get the current working directory
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err)
	}

	baseDir = strings.Replace(baseDir, "cmd/server", "", 1)

	return filepath.Join(baseDir, relativePath)
}
