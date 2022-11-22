package main

import (
	"context"
	edgemysql "github.com/drakejin/crawler/edge/mysql"
	"github.com/drakejin/crawler/internal/crawler"
	storagedb "github.com/drakejin/crawler/internal/storage/db"
	"net/http"
)

var client = &http.Client{}

func main() {

	sql, err := edgemysql.New("admin:passwd@tcp(localhost:23306)/indexer?parseTime=True")
	if err != nil {
		panic(err)
	}
	if err = sql.Ping(); err != nil {
		panic(err)
	}
	storageDB := storagedb.New(sql, true)
	storageDB.Client()

	c := crawler.New(storageDB.Client(), 10, "20221121_2208")

	c.Crawler(context.Background(), nil, "https://9gag.com/trending")
	// c.Crawler(ctx, "", "https://www.hostinger.com/tutorials/uri-vs-url")
}
