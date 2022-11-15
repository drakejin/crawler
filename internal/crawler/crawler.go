package crawler

import (
	"context"
	"github.com/motemen/go-loghttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

// need 전체 고루틴 개수 제어자
//

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

func (c *client) crawler(ctx context.Context, url string) map[string]string {
}
