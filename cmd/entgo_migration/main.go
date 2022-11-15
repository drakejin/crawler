package main

import (
	edgemysql "github.com/drakejin/crawler/edge/mysql"
	storagedb "github.com/drakejin/crawler/internal/storage/db"
	"github.com/rs/zerolog/log"
)

func main() {

	db, err := edgemysql.New("admin:passwd@tcp(localhost:23306)/indexer?parseTime=True")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	serviceDb := storagedb.New(db, true)
	if err = serviceDb.Migrate(); err != nil {
		log.Panic().Err(err).Msg("failed migration at ENV=local")
	}
}
