package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"

	"github.com/drakejin/crawler/internal/_const"
	"github.com/drakejin/crawler/internal/model"
)

func checkContentType(header http.Header) model.ContentType {
	if v := header.Get(_const.HeaderContentType); v == "" {
		return model.ContentTypeNotSupport
	} else {
		if strings.ContainsAny(v, _const.MIMETextHTML) {
			return model.ContentTypeHTML
		} else if strings.ContainsAny(v, "image") {
			return model.ContentTypeImage
		} else if strings.ContainsAny(v, "pdf") {
			return model.ContentTypePDF
		} else if strings.ContainsAny(v, "video") {
			return model.ContentTypeVideo
		} else if strings.ContainsAny(v, _const.MIMEOctetStream) {
			return model.ContentTypeNotSupport // something byte download
		}
	}
	return model.ContentTypeNotSupport
}

func ParseHTML(crawlingVersion string, originUrl *url.URL, body io.ReadCloser) (*model.Page, *model.PageSource, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, nil, err
	}
	if err = body.Close(); err != nil {
		return nil, nil, err
	}

	htmlSource, err := doc.Html()
	if err != nil {
		return nil, nil, err
	}
	if MaximumContentLength < int64(len(htmlSource)) {
		return nil, nil, ErrOverMaximumContentLength
	}

	var isHTTP bool
	if originUrl.Scheme == _const.SchemaHTTPS {
		isHTTP = true
	}

	p := ParseHTMLHead(doc.Find("head"))
	p.ID = uuid.New()
	p.CrawlingVersion = crawlingVersion
	p.Domain = originUrl.Host
	p.IsHTTPS = isHTTP
	p.Path = originUrl.Path
	p.Port = originUrl.Port()
	p.Querystring = originUrl.Query().Encode()
	p.URL = originUrl.String()

	copyUrl, _ := url.Parse(originUrl.String())
	copyUrl.RawQuery = ""

	p.Links = ParseHTMLLinks(copyUrl, doc.Find("body"))

	ps := &model.PageSource{
		ID:     uuid.New(),
		Source: htmlSource,
	}
	return p, ps, nil
}

func ParseHTMLLinks(originUrl *url.URL, bodyDom *goquery.Selection) []string {
	var r []string
	bodyDom.Find("a[href]").Each(func(_ int, aTag *goquery.Selection) {
		if v, ok := aTag.Attr("href"); ok {
			strLength := len([]rune(v))
			if strings.Contains(v, "http") && strLength < 1000 {
				r = append(r, v)
				return
			}

			//// FIXME: 대충 짜본거.. 테스트 코드들이 필요함. url들이 이뻐야함.
			//// filepath package를 이용해서.. 상대경로 코드들은 .. 다.. 치환해버리기
			//if -1 < strings.Index(v, "/") && strings.Index(v, "/") < 3 {
			//	relativePath := []rune(v)
			//	startThreeRunesString := string(relativePath[:3])
			//	startDoubleRunesString := string(relativePath[:2])
			//	startRuneString := string(relativePath[:1])
			//	if startThreeRunesString == "../" {
			//		r = append(r, fmt.Sprintf("%s/%s", originUrl.String(), v))
			//		return
			//	}
			//	if startDoubleRunesString == ".." {
			//		r = append(r, fmt.Sprintf("%s/%s", originUrl.String(), v))
			//		return
			//	} else if startDoubleRunesString == "//" {
			//		r = append(r, fmt.Sprintf("%s:%s", originUrl.Scheme, v))
			//		return
			//	} else if startDoubleRunesString == "./" {
			//		r = append(r, fmt.Sprintf("%s/%s", originUrl.String(), string(relativePath[2:])))
			//		return
			//	}
			//
			//	if startRuneString == "/" {
			//		r = append(r, fmt.Sprintf("%s/%s", originUrl.String(), v))
			//		return
			//	}
			//}
		}
	})
	return r
}

func ParseHTMLHead(headDom *goquery.Selection) *model.Page {
	p := &model.Page{}

	p.Title = headDom.Find("title").Text()

	headDom.Find("meta").Each(func(i int, s *goquery.Selection) {
		var key, value string
		if v, ok := s.Attr("name"); ok {
			key = v
		} else {
			if subV, subOK := s.Attr("property"); subOK {
				key = subV
			}
		}

		if v, ok := s.Attr("content"); ok {
			value = v
		}
		switch key {
		case _const.TagMetaKeywords:
			p.Keywords = value
		case _const.TagMetaContentLanguage:
			p.ContentLanguage = value
		case _const.TagMetaDescription:
			p.Description = value
		case _const.TagMetaTwitterCard:
			p.TwitterCard = value
		case _const.TagMetaTwitterURL:
			p.TwitterURL = value
		case _const.TagMetaTwitterTitle:
			p.TwitterTitle = value
		case _const.TagMetaTwitterDescription:
			p.TwitterDescription = value
		case _const.TagMetaTwitterImage:
			p.TwitterImage = value

		case _const.TagMetaOpenGraphSiteName:
			p.OgSiteName = value
		case _const.TagMetaOpenGraphLocale:
			p.OgLocale = value
		case _const.TagMetaOpenGraphTitle:
			p.OgTitle = value
		case _const.TagMetaOpenGraphDescription:
			p.OgDescription = value
		case _const.TagMetaOpenGraphType:
			p.OgType = value
		case _const.TagMetaOpenGraphURL:
			p.OgURL = value
		case _const.TagMetaOpenGraphImage:
			p.OgImage = value
		case _const.TagMetaOpenGraphImageType:
			p.OgImageType = value
		case _const.TagMetaOpenGraphImageURL:
			p.OgVideoURL = value
		case _const.TagMetaOpenGraphImageSecureURL:
			p.OgVideoSecureURL = value
		case _const.TagMetaOpenGraphImageWidth:
			p.OgVideoWidth = value
		case _const.TagMetaOpenGraphImageHeight:
			p.OgVideoHeight = value
		case _const.TagMetaOpenGraphVideo:
			p.OgVideo = value
		case _const.TagMetaOpenGraphVideoType:
			p.OgVideoType = value
		case _const.TagMetaOpenGraphVideoURL:
			p.OgVideoURL = value
		case _const.TagMetaOpenGraphVideoSecureURL:
			p.OgVideoSecureURL = value
		case _const.TagMetaOpenGraphVideoWidth:
			p.OgVideoWidth = value
		case _const.TagMetaOpenGraphVideoHeight:
			p.OgVideoHeight = value
		}
	})
	return p
}
