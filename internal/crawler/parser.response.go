package crawler

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

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
		}
	}
	return model.ContentTypeNotSupport
}

func ParseHTML(body io.ReadCloser) (*model.Page, *model.PageSource, error) {
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

	return &model.Page{},
		&model.PageSource{
			Source: htmlSource,
		},
		nil
}
