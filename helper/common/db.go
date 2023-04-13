package common

import (
	"database/sql"

	"github.com/hoquangnam45/pharmacy-common-go/helper/db"
	"github.com/hoquangnam45/pharmacy-common-go/helper/migrator"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"gorm.io/gorm"
)

func InitializePostgresDb(host string, username string, password string, database string, port int, gormConfig *gorm.Config, migrationPath string, migrateVersion int) *gorm.DB {
	db := h.FlatMap(
		h.FactoryM(func() (*sql.DB, error) {
			return db.OpenPostgresDb(host, username, password, database, port)
		}),
		h.Lift(func(connection *sql.DB) (*gorm.DB, error) {
			return db.WrapPostgresDbGorm(connection, gormConfig)
		}),
	).PanicEval()

	h.FactoryM(func() (any, error) {
		return nil, migrator.MigratePostgres(host, username, password, database, port, "file://"+migrationPath, migrateVersion)
	}).PanicEval()
	return db
}
