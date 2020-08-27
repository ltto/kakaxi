package kakaxi

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReplaceHTML(body []byte, prefixs ...string) (bbody []byte, err error) {
	buffer := bytes.NewBuffer(body)
	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return body, err
	}
	split(doc, "src", prefixs)
	split(doc, "href", prefixs)
	split(doc, "actions", prefixs)
	html, err := doc.Html()
	if err != nil {
		return body, err
	}
	bbody = []byte(html)
	return
}

func split(doc *goquery.Document, attr string, prefixs []string) {
	if attr == "" {
		return
	}
	doc.Find("[" + attr + "^=http]").Each(func(i int, s *goquery.Selection) {
		gets, _ := s.Attr(attr)
		for _, prefix := range prefixs {
			if strings.Index(gets, prefix) == 0 {
				val := gets[len(prefix):]
				if len(val) != 0 && val[0] == '/' {
					val = val[1:]
				}
				s.SetAttr(attr, val)
			}
		}
	})
}
