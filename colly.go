package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

// CollyResultChan ... result of work of `CrawlSite` function
type CollyResultChan struct {
	URL   string
	Error error
	Done  bool
}

// CrawlSite ... Crawl choosen URL and saves found files
func CrawlSite(urlSite string, saveto string, maxMB float32, maxLoad uint, workMinutes time.Duration, resultChan chan CollyResultChan) {
	url := "https://" + urlSite

	var loadedSize uint
	maxLoadSize := maxLoad * 1024
	waitTime := time.Minute * workMinutes
	exit := false
	c := colly.NewCollector()
	c.AllowedDomains = []string{"*." + urlSite, urlSite}
	// Reduce maximum response body size to 1M
	size := int(1024 * 1024 * maxMB)
	c.MaxBodySize = size

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("[Visiting]", r.URL.String())
	})
	c.OnError(func(_ *colly.Response, err error) {
		resultChan <- CollyResultChan{URL: url, Error: err}
	})

	start := time.Now()
	c.OnResponse(func(r *colly.Response) {
		elapsed := time.Since(start)
		ext := ExtensionByContent(r.Body)
		if elapsed > waitTime && !exit {
			fmt.Println(url, " time end")
			resultChan <- CollyResultChan{URL: urlSite, Done: true}
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
	resultChan <- CollyResultChan{URL: urlSite, Done: true}
}
