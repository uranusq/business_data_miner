package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// CrawlSite ... Crawl choosen URL and saves found files
func CrawlSite(urlSite string, allowedDomains []string, saveto string, maxMB float32, maxLoad uint, workMinutes time.Duration, wg *sync.WaitGroup) {
	f := logToFile(saveto)
	defer f.Close()
	defer wg.Done()

	var loadedSize uint
	maxLoadSize := maxLoad * 1024
	waitTime := time.Minute * workMinutes

	c := colly.NewCollector()
	c.AllowedDomains = allowedDomains
	// Reduce maximum response body size to 1M
	size := int(1024 * 1024 * maxMB)
	c.MaxBodySize = size

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("[Visiting]", r.URL.String())
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Printf("Something went wrong: %v\n", err)
	})

	start := time.Now()
	c.OnResponse(func(r *colly.Response) {
		elapsed := time.Since(start)
		ext := ExtensionByContent(r.Body)
		if elapsed > waitTime {
			fmt.Println(urlSite, " time end")
			wg.Done()
			return
		} else if ext == ".none" {
			return
		} else if loadedSize > maxLoadSize && ext != ".pdf" && ext != ".doc" {
			fmt.Println(urlSite, " loaded: ", loadedSize)
			return
		}
		filename := EscapeURL(r.Request.URL.EscapedPath())
		ioutil.WriteFile(saveto+"/"+filename+ext, r.Body, 0644)
		loadedSize += uint(len(r.Body) / 1024)
	})
	// From where scrapping starts

	c.Visit(urlSite)
	fmt.Println(urlSite, " finished")
}

func logToFile(location string) *os.File {
	f, err := os.OpenFile(location+"/parselog.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("createLogFile error opening file: ", err)
	}
	log.SetOutput(f)
	log.Println("Log started\n-------------------------------\n")
	return f
}
