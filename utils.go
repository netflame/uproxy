package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	_ "regexp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

// -------------------------------- request and response related func --------------------------------
var (
	defaultTimeout           = 6 * time.Second
	defaultClientWithTimeout = &http.Client{Timeout: defaultTimeout}
	defaultHeaders           = http.Header{"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"}}
)

// get response from a given url
func getResp(url string) (resp *http.Response, err error) {
	resp, err = Get(url)
	if err != nil {
		return nil, err
	}
	if err = checkResp(resp); err != nil {
		return nil, err
	}
	return
}

// Get make a `get` request with default headers
func Get(url string) (resp *http.Response, err error) {
	return do("GET", url, nil)
}

// wrapper Do with default headers
func do(method, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := newRequest(method, url, nil, defaultHeaders)
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return c.Do(req)
}

// get a reqest with custom headers
func newRequest(method, url string, body io.Reader, headers http.Header) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value[0])
	}
	return req, nil
}

// ----------------------- geting sites related func ------------------------
// get urls to scrape from default sites
func urlsToScrape() []string {
	toScrape, err := defaultSites()
	check(err)
	urls := make([]string, 0)
	for _, site := range toScrape.Site {
		urls = append(urls, site.URL)
	}
	return urls
}

// get sites to scrape from `sites.toml`
func defaultSites() (s *sites, err error) {
	file := "sites.toml"
	s, err = getSites(file)
	if err != nil {
		return nil, err
	}
	return
}

// get sites to scrape from a given file
func getSites(file string) (*sites, error) {
	var toScrape sites
	var curDir = curdir()
	siteFile := filepath.Join(curDir, file)

	if _, err := toml.DecodeFile(siteFile, &toScrape); err != nil {
		return nil, err
	}

	return &toScrape, nil
}

// get current dir
func curdir() string {
	dir, err := filepath.Abs(filepath.Dir("./utils.go"))
	check(err)
	return dir
}

// --------------------------- check sth func --------------------------
func check(err error) {
	checkErr(err)
}

// exit when err happens
func checkErr(err error) {
	if err != nil {
		miniLog.fatal(err)
	}
}

// check error gracefully
func checkErr2(err error) {
	if err != nil {
		miniLog.error(err)
	}
}

// check http response
func checkResp(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func checkItem(i item) (err error) {
	switch v := i.(type) {
	case *proxy:
		err = checkProxyURL(v.String())
	default:
		err = fmt.Errorf("Unsupported type: %T", v)
	}
	return
}

func checkProxyURL(pURL string) (err error) {
	proxyurl, _ := url.Parse(pURL)
	c := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyurl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: defaultTimeout,
	}
	resp, err := c.Get("https://httpbin.org/ip")
	if err != nil {
		miniLog.error("(utils.go:checkProxyURL)  ", pURL, "->", err)
		return
	}
	defer resp.Body.Close()

	hbIP := &httpbinIP{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, hbIP)
	if err != nil {
		return
	}

	if hbIP.Origin != pURL {
		err = fmt.Errorf("Expected %s, but Caught %s", pURL, hbIP.Origin)
	}

	return
}

// convert url string to proxy
// url format: schema://ip:port, eg, http://0.0.0.0:8080
func getProxyFromURL(u string) *proxy {
	// r := regexp.MustCompile(`(http[s])://(\d+\.\d+\.\d+\.\d+):(\d+)`)
	// res := r.FindStringSubmatch(url)
	U, _ := url.Parse(u)
	port, _ := strconv.Atoi(U.Port())
	return &proxy{
		Scheme: U.Scheme,
		IP:     U.Hostname(),
		Port:   port,
	}
}

// get next url for a relative href
func getNextURL(cur, next string) string {
	u, _ := url.Parse(cur)
	nextURL := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   next,
	}
	return nextURL.String()
}
