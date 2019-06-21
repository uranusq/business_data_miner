package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	cc "github.com/karust/gocommoncrawl"

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
	resChann := make(chan GoogleResultChan)
	sites := []string{"medium.com/", "innopolis.ru", "example.com/", "tutorialspoint.com/"}

	go func() {
		for r := range resChann {
			if r.Error != nil {
				fmt.Printf("Error occured: %v\n", r.Error)
			} else if r.Progress > 0 {
				fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
			}
		}
	}()

	for _, url := range sites {
		saveFolder := saveTo + cc.EscapeURL(url)
		err := createDir(saveFolder)
		if err != nil {
			fmt.Println("[GoogleCrawl] error: %v", err)
		}

		waitTime := time.Second * 20
		start := time.Now()

		go FetchURLFiles(url, "pdf", saveFolder, 1, resChann)

		elapsed := time.Since(start)
		if elapsed < waitTime {
			time.Sleep(waitTime - elapsed)
		}
	}

}

// CollyCrawl ...
func (m Miner) CollyCrawl(saveTo string, wg *sync.WaitGroup) {
	defer wg.Done()

	sites := []string{"medium.com", "innopolis.ru", "example.com", "tutorialspoint.com"}
	var innerWg sync.WaitGroup
	innerWg.Add(len(sites))

	for _, url := range sites {
		saveFolder := saveTo + cc.EscapeURL(url)
		err := createDir(saveFolder)
		if err != nil {
			fmt.Println("[CollyCrawl] error: ", err)
		}
		go CrawlSite("https://"+url, []string{"*." + url, url}, saveFolder, 10, 15, 1, &innerWg)
	}
	innerWg.Wait()
}

func createDir(path string) error {
	// Create directory if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModeDir)
		if err != nil {
			return fmt.Errorf("[createDir] error: %v", err)
		}
	}
	return nil
}

func main() {
	// saveTo, dbFile := "./data", "test.db"
	// commonProcs, googleProcs, collyProcs := 20, 2, 20
	var wg sync.WaitGroup
	wg.Add(1) // 3 When all miners used
	miner := Miner{}
	// miner.db = d.Database{}

	// miner.db.OpenInitialize(dbFile)
	// miner.db.PrintInfo()
	// defer miner.db.Close()

	// 1. Use CommonCrawl to retrive indexed HTML pages of given site
	//miner.CommonCrawl(saveTo, commonProcs)

	// 2. Use Google search with to find cached files
	//miner.GoogleCrawl("./data/google/")

	// 3. Crawl site with gocolly to find unindexed documents
	miner.CollyCrawl("./data/colly/", &wg)
	wg.Wait()
}
