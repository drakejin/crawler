package crawler

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/motemen/go-loghttp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/drakejin/crawler/internal/_const"
	"github.com/drakejin/crawler/internal/model"
	"github.com/drakejin/crawler/internal/storage/db/ent"
)

// need 전체 고루틴 개수 제어자
// recursive ㅇㅇㅇ

type client struct {
	currentRoutinesCount uint32
	maximumConcurrency   int
	client               *http.Client
	storageDB            *ent.Client
}

func New(storageDB *ent.Client, maximumConcurrency int) *client {
	return &client{
		storageDB:          storageDB,
		maximumConcurrency: maximumConcurrency,
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

var (
	ErrNotSupportContentType = errors.New("crawler: only allow 'content-type' about [html]")

	ErrOverMaximumContentLength = errors.New("crawler: page is too big size to indexing, maximum size is 1mb")

	ErrResponseStatusNotOk       = errors.New("crawler: server status is not 200")
	MaximumContentLength   int64 = 1024 * 1024 // 1mb
)

func (c *client) Crawler(ctx context.Context, referredUrl, targetUrl string) {
	// 해당 페이지가 탐색할만한 페이지인지 확인하기 로직
	//     - url이 올바른가?
	// 해당 페이지에서 얻어낸 page link를 다시 rescursive 하게 요청할 수 있어야한다.
	//     - content-type 이 text/html인가? 아닌가?
	//     - content-length 가 1mb이하인가?
	//     - 이 페이지 전에 읽었던 적이 있는가?
	//     - 없었다면 방문한다.
	//     - 실행 버전이 다르면 방문한다. 다시 끌어올 수 있도록
	//     - 있었다고 한다면, 레퍼카운트를 1 올리고 끝
	//     - 있었다고 한다면
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}

	if resp.Status != _const.StatusStrOK {
		if resp.Status == _const.StatusStrTemporaryRedirect ||
			resp.Status == _const.StatusStrMovedPermanently ||
			resp.Status == _const.StatusStrFound ||
			resp.Status == _const.StatusStrSeeOther ||
			resp.Status == _const.StatusStrNotModified ||
			resp.Status == _const.StatusStrUseProxy ||
			resp.Status == _const.StatusStrPermanentRedirect {
			u, err = resp.Location()
			if err != nil {
				err = errors.Wrap(err, "crawler: response header 'location' value is not valid")
				return
			}
			c.Crawler(ctx, referredUrl, u.String())
			return
		}
		log.Warn().Err(ErrResponseStatusNotOk).Send()
		return
	}

	switch checkContentType(resp.Header) {
	case model.ContentTypeHTML:
		if p, ps, err := ParseHTML(resp.Body); err != nil {
		} else {
			if err = c.Save(ctx, p, ps); err != nil {
				log.Warn().Err(err).Send()
			}
			log.Info().Interface("page", p).Str("url", ps.URL).Send()
		}
	default:
		log.Warn().Err(ErrNotSupportContentType).Send()
		return
	}
}

func (c *client) Save(ctx context.Context, p *model.Page, ps *model.PageSource) error {
	// TODO: 20221117: indexed_url 이라는 체제 대신에 domain / path / querystring 셋을 조합키로 둬서 쓰는 방식을 검토해볼예정

	// 이거 하다 말았
	//tx, err := c.storageDB.BeginTx(ctx, &sql.TxOptions{
	//	ReadOnly:  false,
	//	Isolation: nativesql.LevelDefault,
	//})
	//if err != nil {
	//	return err
	//}
	//entgoClient := tx.Client()
	//referredPage, err := entgoClient.Page.Query().Where(entpage.IndexedURLEQ()).Select(entpage.FieldID).Only(ctx)
	//if err != nil {
	//	return err
	//}
	//entgoClient.Page.Create().
	//	OnConflict(
	//		// Update the row with the new values
	//		// the was proposed for insertion.
	//		sql.ResolveWithNewValues(),
	//	).
	//	// Override some of the fields with custom
	//	// update values.
	//	Update(func(u *ent.PageUpsert) {
	//		u.SetReferredID(referredPage.ID)
	//	})

	return nil
}
