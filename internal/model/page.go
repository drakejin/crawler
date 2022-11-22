package model

import (
	"github.com/drakejin/crawler/internal/storage/db/ent"
	"time"

	"github.com/google/uuid"

	"github.com/drakejin/crawler/internal/storage/db/ent/page"
)

type ContentType int

var (
	ContentTypeNotSupport ContentType
	ContentTypeHTML       ContentType
	ContentTypeImage      ContentType
	ContentTypeVideo      ContentType
	ContentTypePDF        ContentType
)

type Page struct {
	// ID of the ent.
	ID    uuid.UUID `json:"id,omitempty"`
	Links []string  `json:"links"`

	// CrawlingVersion holds the value of the "crawling_version" field.
	CrawlingVersion string `json:"crawling_version,omitempty"`
	// domain www.example.com
	Domain string `json:"domain,omitempty"`
	// port number
	Port string `json:"port,omitempty"`
	// is used tls/ssl layer flag
	IsHTTPS bool `json:"is_https,omitempty"`
	// url.path
	Path string `json:"path,omitempty"`
	// url.querystring
	Querystring string `json:"querystring,omitempty"`
	// this mean url
	URL string `json:"url,omitempty"`
	// how many times referred
	CountReferred int64 `json:"count_referred,omitempty"`
	// 해당 row는 쓸 수 있는지? 없는지?
	Status page.Status `json:"status,omitempty"`
	// first indexed time
	CreatedAt time.Time `json:"created_at,omitempty"`
	// first indexed time by which system
	CreatedBy string `json:"created_by,omitempty"`
	// modified time
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// modified by which system
	UpdatedBy string `json:"updated_by,omitempty"`
	// html title tag
	Title string `json:"title,omitempty"`
	// basic meta tags 'description'
	Description string `json:"description,omitempty"`
	// basic meta tags 'keywords'
	Keywords string `json:"keywords,omitempty"`
	// basic meta tags 'content-language'
	ContentLanguage string `json:"content_language,omitempty"`
	// twitter meta tags 'card'
	TwitterCard string `json:"twitter_card,omitempty"`
	// twitter meta tags 'url'
	TwitterURL string `json:"twitter_url,omitempty"`
	// twitter meta tags 'title'
	TwitterTitle string `json:"twitter_title,omitempty"`
	// twitter meta tags 'description'
	TwitterDescription string `json:"twitter_description,omitempty"`
	// twitter meta tags 'image'
	TwitterImage string `json:"twitter_image,omitempty"`
	// og meta tags 'site_name'
	OgSiteName string `json:"og_site_name,omitempty"`
	// og meta tags 'locale'
	OgLocale string `json:"og_locale,omitempty"`
	// og meta tags 'title'
	OgTitle string `json:"og_title,omitempty"`
	// og meta tags 'description'
	OgDescription string `json:"og_description,omitempty"`
	// og meta tags 'type'
	OgType string `json:"og_type,omitempty"`
	// og meta tags 'url'
	OgURL string `json:"og_url,omitempty"`
	// og meta tags 'image'
	OgImage string `json:"og_image,omitempty"`
	// og meta tags 'image:type'
	OgImageType string `json:"og_image_type,omitempty"`
	// og meta tags 'image:url'
	OgImageURL string `json:"og_image_url,omitempty"`
	// og meta tags 'image:secure_url'
	OgImageSecureURL string `json:"og_image_secure_url,omitempty"`
	// og meta tags 'image:width'
	OgImageWidth string `json:"og_image_width,omitempty"`
	// og meta tags 'image:height'
	OgImageHeight string `json:"og_image_height,omitempty"`
	// og meta tags 'video'
	OgVideo string `json:"og_video,omitempty"`
	// og meta tags 'video:type'
	OgVideoType string `json:"og_video_type,omitempty"`
	// og meta tags 'video:url'
	OgVideoURL string `json:"og_video_url,omitempty"`
	// og meta tags 'video:secure_url'
	OgVideoSecureURL string `json:"og_video_secure_url,omitempty"`
	// og meta tags 'video:width'
	OgVideoWidth string `json:"og_video_width,omitempty"`
	// og meta tags 'video:height'
	OgVideoHeight string `json:"og_video_height,omitempty"`
}

type PageSource struct {
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// html view source code
	Source string `json:"source,omitempty"`
}

func ParsePageFromEnt(v *ent.Page) *Page {
	return &Page{
		ID:              v.ID,
		CrawlingVersion: v.CrawlingVersion,
		Domain:          v.Domain,
		Port:            v.Port,
		IsHTTPS:         v.IsHTTPS,
		Path:            v.Path,
		Querystring:     v.Querystring,
		URL:             v.URL,

		CountReferred: v.CountReferred,
		Status:        v.Status,
		CreatedAt:     v.CreatedAt,
		CreatedBy:     v.CreatedBy,
		UpdatedAt:     v.UpdatedAt,
		UpdatedBy:     v.UpdatedBy,

		Title:           v.Title,
		Description:     v.Description,
		Keywords:        v.Keywords,
		ContentLanguage: v.ContentLanguage,

		TwitterCard:        v.TwitterCard,
		TwitterURL:         v.TwitterURL,
		TwitterTitle:       v.TwitterTitle,
		TwitterDescription: v.TwitterDescription,
		TwitterImage:       v.TwitterImage,

		OgSiteName:       v.OgSiteName,
		OgLocale:         v.OgLocale,
		OgTitle:          v.OgTitle,
		OgDescription:    v.OgDescription,
		OgType:           v.OgType,
		OgURL:            v.OgURL,
		OgImage:          v.OgImage,
		OgImageType:      v.OgImageType,
		OgImageURL:       v.OgImageURL,
		OgImageSecureURL: v.OgImageSecureURL,
		OgImageWidth:     v.OgImageWidth,
		OgImageHeight:    v.OgImageHeight,
		OgVideo:          v.OgVideo,
		OgVideoType:      v.OgVideoType,
		OgVideoURL:       v.OgVideoURL,
		OgVideoSecureURL: v.OgVideoSecureURL,
		OgVideoWidth:     v.OgVideoWidth,
		OgVideoHeight:    v.OgVideoHeight,
	}
}
