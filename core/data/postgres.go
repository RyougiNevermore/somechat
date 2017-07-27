package data

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/niflheims-io/qb"
)

var postgres *sql.DB
var dal *qb.QB

func InitPostgres(url string, maxIdle, maxOpen int) error {
	var openErr error
	postgres, openErr = sql.Open("postgres", url)
	if openErr != nil {
		return openErr
	}
	pingErr := postgres.Ping()
	if pingErr != nil {

		return pingErr
	}
	postgres.SetMaxIdleConns(maxIdle)
	postgres.SetMaxOpenConns(maxOpen)
	dal = qb.New(postgres, qb.POSTGRES)
	return nil
}

func DAL() *qb.QB {
	return dal
}

func Postgres() *sql.DB {
	return postgres
}
