package main

import (
	"context"
	"net/http"
	"time"

	edgemysql "github.com/drakejin/crawler/edge/mysql"
	"github.com/drakejin/crawler/internal/crawler"
	storagedb "github.com/drakejin/crawler/internal/storage/db"
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

	c := crawler.New(storageDB.Client(), 10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.Crawler(ctx, "", "https://gall.dcinside.com/board/view/?id=dcbest&no=92769")
}
