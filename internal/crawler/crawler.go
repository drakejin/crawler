package crawler

import (
	"context"
	nativesql "database/sql"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
	"github.com/motemen/go-loghttp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"time"

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
	browser              browser.Browsable
}

func New(storageDB *ent.Client, maximumConcurrency int, crawlingVersion string) *client {
	bow := surf.NewBrowser()
	bow.SetTransport(&loghttp.Transport{
		LogRequest: func(req *http.Request) {
			log.Debug().Interface("header", req.Header).Msgf("[request] %s %s", req.Method, req.URL)
		},
		LogResponse: func(resp *http.Response) {
			log.Debug().Interface("header", resp.Header).Msgf("[response] %d %s", resp.StatusCode, resp.Request.URL)
		},
	})
	bow.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	return &client{
		browser:            bow,
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
	err = c.browser.Open(u.String())
	if err != nil {
		log.Warn().Err(err).Send()
		return
	}
	statusCode := c.browser.StatusCode()
	if statusCode != _const.StatusOK {
		if statusCode == _const.StatusTemporaryRedirect ||
			statusCode == _const.StatusMovedPermanently ||
			statusCode == _const.StatusFound ||
			statusCode == _const.StatusSeeOther ||
			statusCode == _const.StatusNotModified ||
			statusCode == _const.StatusUseProxy ||
			statusCode == _const.StatusPermanentRedirect {
			u, err = url.Parse(c.browser.ResponseHeaders().Get(_const.HeaderLocation))
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

	switch checkContentType(c.browser.ResponseHeaders()) {
	case model.ContentTypeHTML:
		if p, ps, err := ParseHTML(c.crawlingVersion, u, c.browser.Dom()); err != nil {
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
		SetSourceID(referredPage.ID).
		SetTargetID(page.ID).
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
