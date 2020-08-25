package kakaxi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReplaceHTML(body []byte, prefix string) (bbody []byte, err error) {
	buffer := bytes.NewBuffer(body)
	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return body, err
	}
	split(doc, "src", prefix)
	split(doc, "href", prefix)
	split(doc, "actions", prefix)
	html, err := doc.Html()
	if err != nil {
		return body, err
	}
	bbody = []byte(html)
	return
}

func split(doc *goquery.Document, attr, prefix string) {
	if attr == "" {
		return
	}
	doc.Find("[" + attr + "^=http]").Each(func(i int, s *goquery.Selection) {
		gets, _ := s.Attr(attr)
		if strings.Index(gets, prefix) == 0 {
			val := gets[len(prefix):]
			fmt.Println(attr, val)
			s.SetAttr(attr, val)
			fmt.Println(attr, val)
		}
	})
}
