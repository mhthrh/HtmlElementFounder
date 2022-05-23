package Model

import (
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"strings"
)

type Link struct {
	Address    string `json:"address"`
	Internal   bool   `json:"internal"`
	StatusCode int    `json:"statusCode"`
	Err        error  `json:"err"`
}
type Element struct {
	Start string
	End   string
	Value string
	l     *logrus.Entry
}

func New(s, e, v string, l *logrus.Entry) *Element {
	return &Element{
		Start: s,
		End:   e,
		Value: v,
		l:     l,
	}
}

func (e *Element) NewParser(content string, chn *chan string) {
	for _, split := range strings.Split(content, e.Start) {
		first := strings.Index(strings.ReplaceAll(split, "'", "\""), e.Value) + len(e.Value)
		last := strings.Index(strings.ReplaceAll(split[first:], "'", "\""), e.End)
		if first != -1 && last != -1 {
			e.l.Printf("%s value added.\n", split[first:first+last])
			*chn <- split[first : first+last]
		}

	}
}

func (e *Element) NewParser1(content string) []Link {
	var result []Link

	for _, split := range strings.Split(content, e.Start) {
		first := strings.Index(strings.ReplaceAll(split, "'", "\""), e.Value) + len(e.Value)
		last := strings.Index(strings.ReplaceAll(split[first:], "'", "\""), e.End)
		if first != -1 && last != -1 {
			result = append(result, Link{
				Address:    split[first : first+last],
				Internal:   Helper.LinkValid(split[first : first+last]),
				StatusCode: 0,
				Err:        nil,
			})
		}
	}
	return result
}
func (e *Element) OldParser(content string) []Link {
	var links []Link
	z := html.NewTokenizer(strings.NewReader(content))
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, Link{
							Address: attr.Val,
							//	Internal: Helper.LinkValid(attr.Val),
						})
					}

				}
			}

		}
	}
}
