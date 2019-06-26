package main

import (
	"fmt"
	"sync"
	"time"

	cc "mygo/gocommoncrawl"

	d "./db"
	"github.com/BurntSushi/toml"
)

// Miner ... Holds reference of database and do grouping of methods
type Miner struct {
	db d.Database
}

// CommonCrawl ...
func (m Miner) CommonCrawl(config commonConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(config.Path + "/log.txt")
	resChan := make(chan cc.Result)
	companies := m.db.GetCommon()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[CommonCrawl] Error occured: %v\n", r.Error)
				if config.Debug {
					fmt.Printf("Error occured: %v\n", r.Error)
				}
			} else if r.Progress > 0 && config.Debug {
				fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
			} else if r.Done {
				m.db.CommonFinished(r.URL)
				logger.Printf("Common done: %v\n", r.URL)
				if config.Debug {
					fmt.Printf("Commo done: %v\n", r.URL)
				}
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
		saveFolder := config.Path + "/" + cc.EscapeURL(c.URL)

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
func (m Miner) GoogleCrawl(config googleConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(config.Path + "/log.txt")
	resChan := make(chan GoogleResultChan)
	companies := m.db.GetGoogle()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[GoogleCrawl] Error occured: %v\n", r.Error)
				if config.Debug {
					fmt.Printf("[GoogleCrawl] Error occured: %v\n", r.Error)
				}
			} else if r.Progress > 0 && config.Debug {
				fmt.Printf("Progress %v: %v/%v\n", r.URL, r.Progress, r.Total)
			} else if r.Done {
				m.db.GoogleFinished(r.URL)
				logger.Printf("Google done: %v\n", r.URL)
				if config.Debug {
					fmt.Printf("Google done: %v\n", r.URL)
				}
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
		saveFolder := config.Path + "/" + cc.EscapeURL(c.URL)
		err := createDir(saveFolder)
		if err != nil && config.Debug {
			fmt.Println("[GoogleCrawl] error: ", err)
		}

		waitTime := time.Second * 30
		start := time.Now()

		go FetchURLFiles(c.URL, "pdf", saveFolder, 35, resChan)

		elapsed := time.Since(start)
		if elapsed < waitTime {
			time.Sleep(waitTime - elapsed)
		}
	}
	innerWg.Wait()
}

// CollyCrawl ...
func (m Miner) CollyCrawl(config collyConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := logToFile(config.Path + "/log.txt")
	resChan := make(chan CollyResultChan)
	companies := m.db.GetColly()
	var innerWg sync.WaitGroup
	innerWg.Add(len(companies) + 1)

	go func() {
		done := 0
		for r := range resChan {
			if r.Error != nil {
				logger.Printf("[CollyCrawl] Error occured: %v\n", r.Error)
				if config.Debug {
					fmt.Printf("[CollyCrawl] Error occured: %v\n", r.Error)
				}
			} else if r.Done && r.Loaded > 0 {
				m.db.CollyFinished(r.URL)
				logger.Printf("Colly done: %v\n", r.URL)
				if config.Debug {
					fmt.Printf("Colly done: %v\n", r.URL)
				}
				done++
				innerWg.Done()
			} else if r.Done && r.Loaded == 0 {
				logger.Printf("Colly failed: %v\n", r.URL)
				if config.Debug {
					fmt.Printf("Colly failed: %v\n", r.URL)
				}
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
		saveFolder := config.Path + "/" + cc.EscapeURL(c.URL)
		err := createDir(saveFolder)
		if err != nil {
			panic(err)
		}
		go CrawlSite(c.URL, saveFolder, 35, 50, 30, resChan)
	}

	innerWg.Wait()
}

func main() {
	// Try to load config file, if error -> meaningless to proceed
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println("Config load error:")
		panic(err)
	}

	// Initialize miner and database
	miner := Miner{}
	miner.db = d.Database{}
	miner.db.OpenInitialize(config.General.Database)
	miner.db.PrintInfo()
	defer miner.db.Close()

	var wg sync.WaitGroup
	// 1. Use CommonCrawl to retrive indexed HTML pages of given site
	if config.Common.Use {
		wg.Add(1)
		go miner.CommonCrawl(config.Common, &wg)
	}

	// 2. Use Google search with to find cached files
	if config.Google.Use {
		wg.Add(1)
		go miner.GoogleCrawl(config.Google, &wg)
	}

	// 3. Crawl site with gocolly to find unindexed documents
	if config.Colly.Use {
		wg.Add(1)
		go miner.CollyCrawl(config.Colly, &wg)
	}
	wg.Wait()
}
