package crawler

import (
	"context"
	nativesql "database/sql"
	"net/http"
	"net/url"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/motemen/go-loghttp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/drakejin/crawler/internal/_const"
	"github.com/drakejin/crawler/internal/model"
	"github.com/drakejin/crawler/internal/storage/db/ent"
	entpage "github.com/drakejin/crawler/internal/storage/db/ent/page"
	entpagereferred "github.com/drakejin/crawler/internal/storage/db/ent/pagereferred"
)

// need 전체 고루틴 개수 제어자
// recursive ㅇㅇㅇ

type client struct {
	currentRoutinesCount uint32
	maximumConcurrency   int
	crawlingVersion      string
	storageDB            *ent.Client
	client               *http.Client
}

func New(storageDB *ent.Client, maximumConcurrency int, crawlingVersion string) *client {
	tp := &loghttp.Transport{
		LogRequest: func(req *http.Request) {
			log.Debug().Interface("header", req.Header).Msgf("[request] %s %s", req.Method, req.URL)
		},
		LogResponse: func(resp *http.Response) {
			log.Debug().Interface("header", resp.Header).Msgf("[response] %d %s", resp.StatusCode, resp.Request.URL)
		},
	}
	//tp := &http.Transport{}
	//err := http2.ConfigureTransport(tp)
	//if err != nil {
	//	panic(err)
	//}
	return &client{
		client: &http.Client{
			Transport: tp,
		},
		storageDB:          storageDB,
		maximumConcurrency: maximumConcurrency,
		crawlingVersion:    crawlingVersion,
	}
}

var (
	ErrNotSupportContentType = errors.New("crawler: only allow 'content-type' about [html]")

	ErrOverMaximumContentLength = errors.New("crawler: page is too big size to indexing, maximum size is 1mb")

	ErrResponseStatusNotOk       = errors.New("crawler: server status is not 200")
	ErrURLSizeOverMaximum        = errors.New("crawler: url length is over maximum size")
	MaximumContentLength   int64 = 1024 * 1024 // 1mb
)

var userAgentMap = map[string]string{
	"9gag.com": "APIs-Google (+https://developers.google.com/webmasters/APIs-Google.html)",

	"www.reddit.com": "PostmanRuntime/7.29.0",
}

func (c *client) Crawler(ctx context.Context, referredPage *model.Page, targetUrl string) {
	// 해당 페이지가 탐색할만한 페이지인지 확인하기 로직
	//     - url이 올바른가?
	//     - url 길이가 1000자가 넘는가? failed
	// 해당 페이지에서 얻어낸 page link를 다시 rescursive 하게 요청할 수 있어야한다.
	//     - content-type 이 text/html인가? 아닌가?
	//     - content-length 가 1mb이하인가?
	//     - 이 페이지 전에 읽었던 적이 있는가?
	//     - 없었다면 방문한다.
	//     - 실행 버전이 다르면 방문한다. 다시 끌어올 수 있도록
	//     - 있었다고 한다면, 레퍼카운트를 1 올리고 끝
	//     - 있었다고 한다면

	log.Debug().Str("visit", targetUrl).Send()
	if len(targetUrl) > 750 {
		log.Error().Str("url", targetUrl).Err(ErrURLSizeOverMaximum).Send()
		return
	}
	tx, err := c.storageDB.BeginTx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: nativesql.LevelDefault,
	})
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	txClient := tx.Client()
	defer tx.Rollback()
	ok, err := WasVisit(ctx, txClient, c.crawlingVersion, targetUrl)
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}
	if ok {
		log.Warn().Str("visit", targetUrl).Msg("was visit")
		return
	}

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
	// https://developers.google.com/search/docs/crawling-indexing/overview-google-crawlers 의 userAgent를 활용하자.
	if ua, ok := userAgentMap[u.Host]; ok {
		req.Header.Set(_const.HeaderUserAgent, ua)
		req.Header.Set(_const.HeaderAccept, "*/*")
		req.Header.Set(_const.HeaderConnection, "keep-alive")
		req.Header.Set(_const.HeaderAcceptEncoding, "gzip, deflate, br")
	} else {
		req.Header.Set(_const.HeaderUserAgent, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	}
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}

	if resp.StatusCode != _const.StatusOK {
		if resp.StatusCode == _const.StatusTemporaryRedirect ||
			resp.StatusCode == _const.StatusMovedPermanently ||
			resp.StatusCode == _const.StatusFound ||
			resp.StatusCode == _const.StatusSeeOther ||
			resp.StatusCode == _const.StatusNotModified ||
			resp.StatusCode == _const.StatusUseProxy ||
			resp.StatusCode == _const.StatusPermanentRedirect {
			u, err = resp.Location()
			if err != nil {
				err = errors.Wrap(err, "crawler: response header 'location' value is not valid")
				return
			}
			c.Crawler(ctx, referredPage, u.String())
			return
		}
		log.Warn().Err(ErrResponseStatusNotOk).Send()
		return
	}

	switch checkContentType(resp.Header) {
	case model.ContentTypeHTML:
		if p, ps, err := ParseHTML(c.crawlingVersion, resp.Request.URL, resp.Body); err != nil {
			log.Error().Err(err).Send()
		} else {
			page, err := Save(ctx, txClient, p, ps)
			if err != nil {
				log.Warn().Err(err).Send()
				return
			}
			if err = AddReferredIfNotExist(ctx, txClient, referredPage, page); err != nil {
				log.Warn().Err(err).Send()
				return
			}
			if err = tx.Commit(); err != nil {
				log.Warn().Err(err).Send()
				return
			}
			for _, l := range p.Links {
				c.Crawler(ctx, page, l)
			}
		}
	default:
		log.Warn().Err(ErrNotSupportContentType).Send()
		return
	}
}

func WasVisit(ctx context.Context, tx *ent.Client, crawlingVersion, url string) (bool, error) {
	return tx.Page.Query().Where(
		entpage.CrawlingVersionEQ(crawlingVersion),
		entpage.URLEQ(url),
	).Exist(ctx)
}

func AddReferredIfNotExist(ctx context.Context, tx *ent.Client, referredPage, page *model.Page) error {
	if referredPage == nil {
		return nil
	}
	if page == nil {
		return nil
	}
	ok, err := tx.PageReferred.Query().Where(entpagereferred.SourceID(referredPage.ID), entpagereferred.TargetID(page.ID)).Exist(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	err = tx.PageReferred.Create().
		SetID(uuid.New()).
		SetSourceID(referredPage.ID).
		SetTargetID(page.ID).
		SetCreatedBy("crawler").
		SetCreatedAt(time.Now().In(time.UTC)).
		SetUpdatedBy("crawler").
		SetUpdatedAt(time.Now().In(time.UTC)).
		OnConflict(
			sql.ResolveWith(func(set *sql.UpdateSet) {
				set.Set("id", uuid.New())
			}),
			sql.ConflictWhere(
				sql.And(
					sql.EQ("source_id", referredPage.ID),
					sql.EQ("target_id", page.ID),
				),
			),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Save(ctx context.Context, txClient *ent.Client, p *model.Page, ps *model.PageSource) (*model.Page, error) {
	ids, err := txClient.Page.Query().Where(entpage.URLEQ(p.URL)).IDs(ctx)
	if err != nil {
		return nil, err
	}
	var id uuid.UUID
	if len(ids) == 0 {
		id = p.ID
		pCreator := txClient.Page.Create().
			SetID(id).
			SetCrawlingVersion(p.CrawlingVersion).
			SetDomain(p.Domain).
			SetPort(p.Port).
			SetIsHTTPS(p.IsHTTPS).
			SetPath(p.Path).
			SetQuerystring(p.Querystring).
			SetURL(p.URL).
			SetStatus(entpage.StatusALLOW).
			SetUpdatedAt(time.Now().In(time.UTC)).
			SetUpdatedBy("worker").
			SetCreatedBy("worker").
			SetTitle(p.Title).
			SetKeywords(p.Keywords).
			SetDescription(p.Description).
			SetContentLanguage(p.ContentLanguage).
			SetTwitterCard(p.TwitterCard).
			SetTwitterURL(p.TwitterURL).
			SetTwitterTitle(p.TwitterTitle).
			SetTwitterDescription(p.TwitterDescription).
			SetTwitterImage(p.TwitterImage).
			SetOgSiteName(p.OgSiteName).
			SetOgLocale(p.OgLocale).
			SetOgTitle(p.OgTitle).
			SetOgDescription(p.OgDescription).
			SetOgType(p.OgType).
			SetOgURL(p.OgURL).
			SetOgImage(p.OgImage).
			SetOgImageType(p.OgImageType).
			SetOgImageURL(p.OgImageURL).
			SetOgImageSecureURL(p.OgImageSecureURL).
			SetOgImageWidth(p.OgImageWidth).
			SetOgImageHeight(p.OgImageHeight).
			SetOgVideo(p.OgVideo).
			SetOgVideoType(p.OgVideoType).
			SetOgVideoURL(p.OgVideoURL).
			SetOgVideoSecureURL(p.OgVideoSecureURL).
			SetOgVideoWidth(p.OgVideoWidth).
			SetOgVideoHeight(p.OgVideoHeight)

		page, err := pCreator.Save(ctx)
		if err != nil {
			return nil, err
		}
		err = txClient.PageSource.Create().
			SetID(id).
			SetSource(ps.Source).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		r := model.ParsePageFromEnt(page)
		r.Links = p.Links
		return r, nil
	}
	id = ids[0]
	page, err := txClient.Page.UpdateOneID(id).
		SetCrawlingVersion(p.CrawlingVersion).
		SetDomain(p.Domain).
		SetPort(p.Port).
		SetIsHTTPS(p.IsHTTPS).
		SetPath(p.Path).
		SetQuerystring(p.Querystring).
		SetURL(p.URL).
		SetStatus(entpage.StatusALLOW).
		SetUpdatedAt(time.Now().In(time.UTC)).
		SetUpdatedBy("worker").
		SetCreatedBy("worker").
		SetTitle(p.Title).
		SetKeywords(p.Keywords).
		SetDescription(p.Description).
		SetContentLanguage(p.ContentLanguage).
		SetTwitterCard(p.TwitterCard).
		SetTwitterURL(p.TwitterURL).
		SetTwitterTitle(p.TwitterTitle).
		SetTwitterDescription(p.TwitterDescription).
		SetTwitterImage(p.TwitterImage).
		SetOgSiteName(p.OgSiteName).
		SetOgLocale(p.OgLocale).
		SetOgTitle(p.OgTitle).
		SetOgDescription(p.OgDescription).
		SetOgType(p.OgType).
		SetOgURL(p.OgURL).
		SetOgImage(p.OgImage).
		SetOgImageType(p.OgImageType).
		SetOgImageURL(p.OgImageURL).
		SetOgImageSecureURL(p.OgImageSecureURL).
		SetOgImageWidth(p.OgImageWidth).
		SetOgImageHeight(p.OgImageHeight).
		SetOgVideo(p.OgVideo).
		SetOgVideoType(p.OgVideoType).
		SetOgVideoURL(p.OgVideoURL).
		SetOgVideoSecureURL(p.OgVideoSecureURL).
		SetOgVideoWidth(p.OgVideoWidth).
		SetOgVideoHeight(p.OgVideoHeight).Save(ctx)
	if err != nil {
		return nil, err
	}
	err = txClient.PageSource.UpdateOneID(id).SetSource(ps.Source).Exec(ctx)
	if err != nil {
		return nil, err
	}
	r := model.ParsePageFromEnt(page)
	r.Links = p.Links

	return r, nil
}
