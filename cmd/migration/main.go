package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "")
	flag.StringVar(&migrationsPath, "migrations-path", "", "")
	flag.StringVar(&migrationsTable, "migrations-table", "", "")
	flag.Parse()

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s/%s?sslmode=disable", storagePath, migrationsTable))
	if err != nil {
		fmt.Printf(err.Error())
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		migrationsTable, driver)
	m.Up()
}
