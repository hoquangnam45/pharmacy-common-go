package db

import (
	"database/sql"
	"errors"
	"fmt"

	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgresDb(postgresHost, username, password, databaseName string, port int) (*sql.DB, error) {
	return h.FlatMap(
		h.Just(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, postgresHost, port, databaseName)),
		h.Lift(func(dsn string) (*sql.DB, error) {
			return sql.Open("postgres", dsn)
		})).Eval()
}

func WrapPostgresDbGorm(db *sql.DB, gormConfig *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{Conn: db}), gormConfig)
}

func IsDuplicatedError(err error) bool {
	if err == nil {
		return false
	}
	var perr *pq.Error
	if errors.As(err, &perr) && perr.Code == "23505" {
		return true
	}
	return false
}
