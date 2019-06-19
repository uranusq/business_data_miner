package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// CrawlSite ... Crawl choosen URL and saves found files
func CrawlSite(urlSite string, allowedDomains []string, saveto string, maxMB float32) {
	f := logToFile(saveto)
	defer f.Close()

	c := colly.NewCollector()
	c.AllowedDomains = allowedDomains
	// Reduce maximum response body size to 1M
	size := int(1024 * 1024 * maxMB)
	c.MaxBodySize = size

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("[Visiting]", r.URL.String())
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: %v\n", err)
	})
	c.OnResponse(func(r *colly.Response) {
		ext := ExtensionByContent(r.Body)
		filename := EscapeURL(r.Request.URL.EscapedPath())
		ioutil.WriteFile(saveto+"/"+filename+ext, r.Body, 0644)

	})
	// From where scrapping starts
	c.Visit(urlSite)
}

func logToFile(location string) *os.File {
	f, err := os.OpenFile(location+"/parselog.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("createLogFile error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println("Log started\n-------------------------------\n")
	return f
}
