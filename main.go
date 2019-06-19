package main

import (
	"fmt"
	cc "mygo/gocommoncrawl"
	"time"

	d "./db"
)

// Miner ...
type Miner struct {
	db d.Database
}

// CommonCrawl ...
func (m Miner) CommonCrawl(saveTo string) {
	resChann := make(chan cc.Result)
	sites := []string{"medium.com/", "example.com/", "tutorialspoint.com/"}
	for _, url := range sites {
		saveFolder := saveTo + cc.EscapeURL(url)
		go cc.FetchURLData(url, saveFolder, resChann, 30, "")
	}

	for r := range resChann {
		if r.Error != nil {
			fmt.Printf("Error occured: %v\n", r.Error)
		} else if r.Progress > 0 {
			fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
		}
	}
}

// GoogleCrawl ...
func (m Miner) GoogleCrawl(saveTo string) {
	// make channel and add extension to file
	sites := []string{"medium.com/", "innopolis.ru", "example.com/", "tutorialspoint.com/"}
	for _, url := range sites {
		saveFolder := saveTo + cc.EscapeURL(url)
		waitTime := time.Second * 20
		start := time.Now()

		go FetchURLFiles(url, "pdf", saveFolder, 1)

		elapsed := time.Since(start)
		if elapsed < waitTime {
			time.Sleep(waitTime - elapsed)
		}
	}

}

// CollyCrawl ...
func (m Miner) CollyCrawl() {
	CrawlSite("https://olymp.innopolis.ru/stem/", []string{"olymp.innopolis.ru"}, "./data", 1)
}

func main() {

	// saveTo, dbFile := "./data", "test.db"
	// commonProcs, googleProcs, collyProcs := 20, 2, 20

	miner := Miner{}
	// miner.db = d.Database{}

	// miner.db.OpenInitialize(dbFile)
	// miner.db.PrintInfo()
	// defer miner.db.Close()

	// 1. Use CommonCrawl to retrive indexed HTML pages of given site
	//miner.CommonCrawl(saveTo, commonProcs)

	// 2. Use Google search with to find cached files
	miner.GoogleCrawl("./data/google/")

	// 3. Crawl site with gocolly to find unindexed documents
	//collyProcs = 1
	//miner.CollyCrawl()

}
