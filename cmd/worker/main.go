package main

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	edgemysql "github.com/drakejin/crawler/edge/mysql"
	storagedb "github.com/drakejin/crawler/internal/storage/db"
	"net/http"
	"time"
)

var client = &http.Client{}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://gall.dcinside.com/board/view/?id=dcbest&no=92769",
		nil,
	)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	// Find the review items
	meta := make(map[string]string)
	doc.Find("head meta").Each(func(i int, s *goquery.Selection) {
		var name, content string
		if v, ok := s.Attr("name"); ok {
			name = v
		}

		if v, ok := s.Attr("content"); ok {
			content = v
		}
		if name != "" && content != "" {
			meta[name] = content
		}
	})

	sql, err := edgemysql.New("admin:passwd@tcp(localhost:23306)/indexer?parseTime=True")
	if err != nil {
		panic(err)
	}
	if err = sql.Ping(); err != nil {
		panic(err)
	}
	storageDB := storagedb.New(sql, true)
	storageDB.Client()

}
