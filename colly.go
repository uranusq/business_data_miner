package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	cly "github.com/gocolly/colly"
)

// CollyResultChan ... result of work of `CrawlSite` function
type CollyResultChan struct {
	URL    string
	Error  error
	Loaded uint
	Done   bool
}

// CrawlSite ... Crawl choosen URL and saves found files
func CrawlSite(urlSite string, saveto string, maxMB int, maxLoad uint, workMinutes int, resultChan chan CollyResultChan) {
	url := "https://" + urlSite

	var loadedSize uint
	maxLoadSize := maxLoad * 1024
	waitTime := time.Minute * time.Duration(workMinutes)
	exit := false
	c := cly.NewCollector()
	c.AllowedDomains = []string{"www." + urlSite, "sso." + urlSite, urlSite, "s3.amazonaws.com"}
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   180 * time.Second,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 60 * time.Second,
	})
	// Reduce maximum response body size to 1M
	size := int(1024 * 1024 * maxMB)
	c.MaxBodySize = size

	c.OnHTML("a[href]", func(e *cly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *cly.Request) {
		r.Headers.Set("User-Agent", randomOption(userAgents))
		//fmt.Println("[Visiting]", r.URL.String())
	})
	c.OnError(func(_ *cly.Response, err error) {
		resultChan <- CollyResultChan{URL: url, Error: err}
	})

	start := time.Now()
	c.OnResponse(func(r *cly.Response) {
		elapsed := time.Since(start)
		ext := ExtensionByContent(r.Body)
		if elapsed > waitTime && !exit {
			//fmt.Println(url, " time end")
			resultChan <- CollyResultChan{URL: urlSite, Done: true, Loaded: loadedSize}
			exit = true
			return
		} else if ext == ".none" {
			return
		} else if loadedSize > maxLoadSize && ext != ".pdf" && ext != ".doc" {
			return
		}
		filename := EscapeURL(r.Request.URL.EscapedPath())
		r.Save(saveto + "/" + filename + randString(6) + ext)
		loadedSize += uint(len(r.Body) / 1024)
	})

	c.Visit(url)
	// If it crawled less than 25 Kb - try again, but with `www.` domain
	if loadedSize < 1024*25 {
		url = "https://www." + urlSite
		c.Visit(url)
	}
	resultChan <- CollyResultChan{URL: urlSite, Done: true, Loaded: loadedSize}
}
