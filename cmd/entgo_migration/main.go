package main

import (
	"github.com/rs/zerolog/log"

	edgemysql "github.com/drakejin/crawler/edge/mysql"
	"github.com/drakejin/crawler/internal/config"
	storagedb "github.com/drakejin/crawler/internal/storage/db"
)

func main() {
	cfg := config.MustNew()
	db, err := edgemysql.New(cfg.DB.DSN)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	serviceDb := storagedb.New(db, true)
	if err = serviceDb.Migrate(); err != nil {
		log.Panic().Err(err).Msg("failed migration at ENV=local")
	}
}
