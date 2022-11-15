package crawler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/motemen/go-loghttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// need 전체 고루틴 개수 제어자
// recursive ㅇㅇㅇ

type client struct {
	currentRoutinesCount uint32
	maximumConcurrency   int
	startpoint           string
	client               *http.Client
}

func New(startpoint string, currentRoutinesCount uint32, maximumConcurrency int) *client {
	return &client{
		startpoint:           startpoint,
		currentRoutinesCount: currentRoutinesCount,
		maximumConcurrency:   maximumConcurrency,
		client: &http.Client{
			Transport: &loghttp.Transport{
				LogRequest: func(req *http.Request) {
					log.Debug().Msgf("[%p] %s %s", req, req.Method, req.URL)
				},
				LogResponse: func(resp *http.Response) {
					log.Debug().Msgf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)
				},
			},
		},
	}
}

func (c *client) crawler(ctx context.Context, url string) (map[string]string, error) {
	// 해당 페이지가 탐색할만한 페이지인지 확인하기 로직
	//     - content-type 이 text/html인가? 아닌가?
	// 해당 페이지에서 얻어낸 page link를 다시 rescursive 하게 요청할 수 있어야한다.
	//     - 이 페이지 전에 읽었던 적이 있는가?
	//     - 없었다면 방문한다.
	//     - 실행 버전이 다르면 방문한다. 다시 끌어올 수 있도록
	//     - 있었다고 한다면, 레퍼카운트를 1 올리고 끝
	//     - 있었다고 한다면
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
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
	return nil, nil
}
