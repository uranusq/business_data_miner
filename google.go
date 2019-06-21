// https://gist.github.com/EdmundMartin/eaea4aaa5d231078cb433b89878dbecf
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GoogleResultChan ... result of work of `FetchURLFiles` function
type GoogleResultChan struct {
	URL      string
	Progress int
	Total    int
	Error    error
}

// GoogleResult ... Result of Google search
type GoogleResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"uk":  "https://www.google.co.uk/search?q=",
	"ru":  "https://www.google.ru/search?q=",
	"fr":  "https://www.google.fr/search?q=",
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func buildGoogleURL(searchTerm string, countryCode string, languageCode string) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleBase, searchTerm, languageCode)
	} else {
		return fmt.Sprintf("%s%s&num=100&hl=%s", googleDomains["com"], searchTerm, languageCode)
	}
}

func googleRequest(searchURL string) (*http.Response, error) {
	baseClient := &http.Client{}
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomOption(userAgents))

	res, err := baseClient.Do(req)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func googleResultParser(response *http.Response) ([]GoogleResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []GoogleResult{}
	sel := doc.Find("div.g")
	rank := 1
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" {
			result := GoogleResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank += 1
		}
	}
	return results, err
}

// GoogleScrape ...
func GoogleScrape(searchTerm string, countryCode string, languageCode string) ([]GoogleResult, error) {
	googleUrl := buildGoogleURL(searchTerm, countryCode, languageCode)
	res, err := googleRequest(googleUrl)
	if err != nil {
		return nil, err
	}
	scrapes, err := googleResultParser(res)
	if err != nil {
		return nil, err
	} else {
		return scrapes, nil
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(saveto string, extension string, url string, maxMegabytes uint64) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := saveto + "/" + FilenameFromURL(url)
	// Create the file
	var out *os.File
	if _, err := os.Stat(filename); err == nil {
		out, err = os.Create(filename + randString(10) + "." + extension)
		if err != nil {
			return err
		}
	} else {
		out, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	defer out.Close()

	// Write the body to file
	//_, err = io.Copy(out, resp.Body)
	megabytes := int64(maxMegabytes * 1024000)
	_, err = io.CopyN(out, resp.Body, megabytes)
	if err != nil {
		return err
	}
	return err
}

// FetchURLFiles ...
func FetchURLFiles(url string, extension string, saveto string, maxMegabytes uint64, resultChan chan GoogleResultChan) {
	// Query google
	query := fmt.Sprintf("site:%v filetype:%v", url, extension)
	res, err := GoogleScrape(query, "ru", "RU")
	if err != nil {
		resultChan <- GoogleResultChan{Error: fmt.Errorf("[FetchURLFiles] error: %v", err), URL: url}
		return
	}

	// Download found files
	for i, r := range res {
		DownloadFile(saveto, extension, r.ResultURL, maxMegabytes)
		if err != nil {
			resultChan <- GoogleResultChan{Error: fmt.Errorf("[FetchURLFiles] error: %v", err), URL: url}
			return
		}
		resultChan <- GoogleResultChan{URL: url, Total: len(res), Progress: i + 1}
	}
}

func randString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func randomOption(options []string) string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(options)
	return options[randNum]
}
