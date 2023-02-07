package migrator

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hoquangnam45/pharmacy-common-go/helper/db"
	_ "github.com/lib/pq"
)

func MigratePostgres(postgresHost, username, password, databaseName string, port int, migrationFilePath string) error {
	db, err := db.OpenPostgresDb(postgresHost, username, password, databaseName, port)
	if err != nil {
		return err
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()
	mi, err := migrate.NewWithDatabaseInstance(migrationFilePath, databaseName, driver)
	if err != nil {
		return err
	}
	defer mi.Close()

	err = mi.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
