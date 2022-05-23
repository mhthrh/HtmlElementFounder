package Controler

import (
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper"
	"GitHub.com/mhthrh/HtmlElementsFinder/Model"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var (
	lnkChn chan string
	termin chan bool
	URL    *url.URL
)

type HtmlPage struct {
	Title    string
	Version  string
	HasLogin bool
	Heading  map[string]int
	Links    *[]Model.Link
}

type Controller struct {
	l *logrus.Entry
}

func New(l *logrus.Entry) *Controller {
	return &Controller{l}
}
func (c *Controller) HtmlElementInfo(req Model.Request) (*HtmlPage, error) {
	var Lnks []Model.Link
	lnkChn = make(chan string)
	termin = make(chan bool)

	page, err := Helper.ReadPage(req.Address)
	if err != nil {
		return nil, err
	}
	result := HtmlPage{
		Title:   Model.New("<TITLE>", "</TITLE>", "", c.l).HtmlTag(page),
		Version: Model.New("<!DOCTYPE", ">", "", c.l).HtmlTag(page),
		HasLogin: func(c string) bool {
			var i int
			words := []string{"Login", "username", "password", "submit"}
			for _, s := range words {
				if Model.New(s, "", "", nil).HeadCount(c) > 0 {
					i++
				}
			}
			if i >= len(words)-1 {
				return true
			}
			return false
		}(page),
		Heading: map[string]int{"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0},
		Links:   &Lnks,
	}

	for i := req.ThreadCount; i > 0; i-- {
		go find(req.Address, &lnkChn, &termin, &Lnks)
	}

	Model.New("</a>", "\"", "href=\"", c.l).NewParser(page, &lnkChn)

	for s, _ := range result.Heading {
		result.Heading[s] = Model.New(s, fmt.Sprintf("</%s>", s), "", c.l).HeadCount(page)
	}

	for i := req.ThreadCount; i > 0; i-- {
		termin <- false
	}
	return &result, nil
}

func find(address string, b *chan string, terminate *chan bool, l *[]Model.Link) {
	client := http.Client{
		Timeout: 7 * time.Second,
	}
	for {
		select {
		case s := <-*b:
			u, inter, err := isInternalUrl(address, s)
			res, err := client.Get(u)
			if err != nil {
				*l = append(*l, Model.Link{
					Address:    u,
					Internal:   inter,
					StatusCode: -99,
					Err:        err,
				})
				break
			}
			*l = append(*l, Model.Link{
				Address:    u,
				Internal:   inter,
				StatusCode: res.StatusCode,
				Err:        nil,
			})

		case b := <-*terminate:
			if !b {
				fmt.Println("goroutine died!")
				return
			}
		}
	}
}

func isInternalUrl(mainUrl, subUrl string) (string, bool, error) {
	main, err := url.Parse(mainUrl)
	if err != nil {
		return "", false, err
	}

	sub, err := url.Parse(subUrl)
	if err != nil {
		return "", false, err
	}
	if main.Host == sub.Host || sub.Host == "" {
		return fmt.Sprintf("%s://%s%s", main.Scheme, main.Host, sub.Path), true, nil
	}

	return fmt.Sprintf("%s", sub), false, nil
}
