package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic().Err(err).Msgf("cannot connect with service db [%s]", dsn)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
