package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
type parseFunc map[string]func(*http.Response) error

var pf = parseFunc{
	"kuaidaili": parseKDL,
	"xicidaili": parseXC,
}
*/

// Spider a struct of spider
// Name: spider's name
// ToScrape: site paris with name and url
type Spider struct {
	Name     string
	ToScrape *sites
}

// Scrape scrape sites according related settings
func (s *Spider) Scrape() {
	sitesToScrape := s.ToScrape.Site
	ch := make(chan miniChan, len(sitesToScrape))
	for _, st := range sitesToScrape {
		go Scrape(st, ch) // WARNING：给的是值，协程中也只能传值
	}
	for i := 0; i < len(sitesToScrape); i++ {
		<-ch
	}
}

// Scrape scrape a given site
func Scrape(s site, ch chan miniChan) {
	defer func() {
		ch <- getMCElement()
	}()
	name, url := s.Name, s.URL
	c := make(chan miniChan, 1)
	go scrape(name, url, c)
	<-c
}

// scrape get `url`'s resp and parse it using suitable `parser` through `name`
func scrape(name, url string, c chan miniChan) {
	defer func() {
		c <- getMCElement()
	}()

	resp, err := getResp(url)
	if err != nil {
		miniLog.error("(spiders.go:scrape)", name, err)
		return
	}
	defer resp.Body.Close()

	// err = pf[name](resp)
	switch name {
	case "kuaidaili":
		err = parseKDL(resp)
	case "xicidaili":
		err = parseXC(resp)
	default:
		err = fmt.Errorf("Unsupported site: %s", name)
	}
	if err != nil {
		miniLog.error("(spiders.go:scrape)", err)
		return
	}
}

// parser for kuaidaili
func parseKDL(resp *http.Response) (err error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	r := make([]*proxy, 0)

	doc.Find("#list tbody tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		ip := tr.Find("td").Eq(0).Text()
		// ip := tr.Find("td:nth-child(1)").Text()						// ok
		port, err := strconv.Atoi(tr.Find("td").Eq(1).Text())
		if err != nil {
			return false
		}
		// port := tr.Find("td:nth-child(2)").Text() 					// ok
		scheme := strings.ToLower(tr.Find("td").Eq(3).Text())
		// schema := strings.ToLower(tr.Find("td:nth-child(4)").Text()) // ok
		p := &proxy{
			Scheme: scheme,
			IP:     ip,
			Port:   port,
		}
		r = append(r, p)

		return true
	})

	// plc: channel for pipeline process
	plc := make(chan miniChan, 1)
	go processProxies(r, plc)
	<-plc

	next, exists := doc.Find("#listnav ul .active").Parent().Next().Find("a").Attr("href")
	if exists {
		nextURL := getNextURL(resp.Request.URL.String(), next)
		miniLog.info("kuaidaili next ->  ", nextURL)
		c := make(chan miniChan, 1)
		go scrape("kuaidaili", nextURL, c)
		<-c
	} else {
		miniLog.error("kuaidaili next -> ", next)
	}

	return
}

// parser for xicidaili
func parseXC(resp *http.Response) (err error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	r := make([]*proxy, 0)

	first := doc.Find("#ip_list tr").Eq(1)
	p, err := getProxyForXC(first)
	if err != nil {
		return
	}
	r = append(r, p)

	first.NextAll().EachWithBreak(func(i int, tr *goquery.Selection) bool {
		p, err := getProxyForXC(tr)
		if err != nil {
			return false
		}
		r = append(r, p)

		return true
	})

	plc := make(chan miniChan, 1)
	go processProxies(r, plc)
	<-plc

	next, exists := doc.Find("a.next_page").First().Attr("href")
	if exists {
		nextURL := getNextURL(resp.Request.URL.String(), next)
		miniLog.info("xici next -> ", nextURL)
		c := make(chan miniChan, 1)
		go scrape("xicidaili", nextURL, c)
		<-c
	} else {
		miniLog.error("xici next -> ", next)
	}
	// fmt.Println(strings.Join(r, "\n"))
	return
}

// only works for xicidaili
func getProxyForXC(s *goquery.Selection) (*proxy, error) {
	ip := s.Find("td:nth-child(2)").Text()
	port, err := strconv.Atoi(s.Find("td:nth-child(3)").Text())
	if err != nil {
		return nil, err
	}
	scheme := strings.ToLower(s.Find("td:nth-child(6)").Text())
	p := &proxy{
		Scheme: scheme,
		IP:     ip,
		Port:   port,
	}
	return p, nil
}
