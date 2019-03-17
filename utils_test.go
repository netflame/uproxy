package main

import (
	"testing"
)

var siteNum = 2

func TestDefaultSites(t *testing.T) {
	ds, err := defaultSites()
	if err != nil {
		t.Error(err)
	}
	if len(ds.Site) != siteNum {
		t.Error("Failed: Got null data")
	}
	siteNames := []string{"kuaidaili", "xicidaili"}
	for _, site := range ds.Site {
		if name := site.Name; !contains(siteNames, name) {
			t.Error("Failed in `defaultSites`")
		}
	}
}
func TestURLsToScrape(t *testing.T) {
	urls := urlsToScrape()
	if len(urls) != siteNum {
		t.Error("Failed in `urlsToScrape`")
	}
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func TestGetProxyFromURL(t *testing.T) {
	url := "https://0.0.0.0:8080"
	p := getProxyFromURL(url)
	if p.String() != url {
		t.Errorf("Expected %s, Got %s", url, p)
	}
}

func TestGetNEextURL(t *testing.T) {
	cur, next := "http://www.baidu.com/first/whatever", "/second/whatever"
	nextURL := getNextURL(cur, next)
	expected := "http://www.baidu.com/second/whatever"
	if nextURL != expected {
		t.Errorf("Expected %s, Got %s", expected, nextURL)
	}
}
