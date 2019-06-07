package main

import (
	"fmt"
	"io/ioutil"
	"github.com/gocolly/colly"
	"github.com/h2non/filetype"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("olymp.innopolis.ru"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		kind, _ := filetype.Match(r.Body)
		if kind == filetype.Unknown {
			fmt.Println("Unknown file type")
			return
		}

		fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
		ioutil.WriteFile("some.pdf", r.Body, 0644)

	})
	// Start scraping on https://hackerspaces.org
	c.Visit("https://olymp.innopolis.ru/stem/")
}