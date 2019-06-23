package main

import (
	"fmt"
	"sync"
	"time"

	cc "mygo/gocommoncrawl"

	d "./db"
)

// Miner ...
type Miner struct {
	db d.Database
}

// CommonCrawl ...
func (m Miner) CommonCrawl(saveTo string, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(saveTo + "/log.txt")
	resChan := make(chan cc.Result)
	companies := m.db.GetCommon()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[CommonCrawl] Error occured: %v\n", r.Error)
				//fmt.Printf("Error occured: %v\n", r.Error)
			} else if r.Progress > 0 {
				//fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
			} else if r.Done {
				m.db.CommonFinished(r.URL)
				logger.Printf("Common done: %v\n", r.URL)
				fmt.Printf("Commo done: %v\n", r.URL)
				done++
				innerWg.Done()
			}
			if done == len(companies) {
				break
			}
		}
		innerWg.Done()
	}()

	for _, c := range companies {
		saveFolder := saveTo + cc.EscapeURL(c.URL)

		waitTime := time.Second * 1
		start := time.Now()

		go cc.FetchURLData(c.URL, saveFolder, resChan, 30, "", 53)

		elapsed := time.Since(start)
		if elapsed < waitTime {
			time.Sleep(waitTime - elapsed)
		}
	}
	innerWg.Wait()
}

// GoogleCrawl ...
func (m Miner) GoogleCrawl(saveTo string, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(saveTo + "/log.txt")
	resChan := make(chan GoogleResultChan)
	companies := m.db.GetGoogle()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[GoogleCrawl] Error occured: %v\n", r.Error)
			} else if r.Progress > 0 {
				//fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
			} else if r.Done {
				m.db.GoogleFinished(r.URL)
				logger.Printf("Google done: %v\n", r.URL)
				fmt.Printf("Google done: %v\n", r.URL)
				done++
				innerWg.Done()
			}
			if done == len(companies) {
				break
			}
		}
		innerWg.Done()
	}()

	for _, c := range companies {
		saveFolder := saveTo + cc.EscapeURL(c.URL)
		err := createDir(saveFolder)
		if err != nil {
			fmt.Println("[GoogleCrawl] error: ", err)
		}

		waitTime := time.Second * 30
		start := time.Now()

		go FetchURLFiles(c.URL, "pdf", saveFolder, 20, resChan)

		elapsed := time.Since(start)
		if elapsed < waitTime {
			time.Sleep(waitTime - elapsed)
		}
	}
	innerWg.Wait()
}

// CollyCrawl ...
func (m Miner) CollyCrawl(saveTo string, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(saveTo + "/log.txt")
	resChan := make(chan CollyResultChan)
	companies := m.db.GetColly()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[CollyCrawl] Error occured: %v\n", r.Error)
				//fmt.Printf("[CollyCrawl] Error occured: %v\n", r.Error)
			} else if r.Done {
				m.db.CollyFinished(r.URL)
				logger.Printf("Colly done: %v\n", r.URL)
				fmt.Printf("Colly done: %v\n", r.URL)
				done++
				innerWg.Done()
			}
			if done == len(companies) {
				break
			}
		}
		innerWg.Done()
	}()

	for _, c := range companies {

		saveFolder := saveTo + cc.EscapeURL(c.URL)
		err := createDir(saveFolder)
		if err != nil {
			panic(err)
		}
		go CrawlSite(c.URL, saveFolder, 20, 50, 30, resChan)
	}

	innerWg.Wait()
}

func main() {
	saveTo, dbFile := "./data", "prod.db"
	// commonProcs, googleProcs, collyProcs := 20, 2, 20
	var wg sync.WaitGroup
	wg.Add(3) // 3 When all miners used
	miner := Miner{}
	miner.db = d.Database{}

	miner.db.OpenInitialize(dbFile)
	miner.db.PrintInfo()
	defer miner.db.Close()

	// 1. Use CommonCrawl to retrive indexed HTML pages of given site
	go miner.CommonCrawl(saveTo+"/common/", &wg)

	// 2. Use Google search with to find cached files
	go miner.GoogleCrawl(saveTo+"/google/", &wg)

	// 3. Crawl site with gocolly to find unindexed documents
	go miner.CollyCrawl(saveTo+"/colly/", &wg)
	wg.Wait()
}
