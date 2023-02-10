package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var googleDomains = map[string]string{}

type SearchResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm, languageCode, countryCode string, pages, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBase, searchTerm, count, languageCode, start)
			toScrape = append(toScrape, scrapeURL)
		}
	} else {
		err := fmt.Errorf("Country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape, nil
}

func getScrapeClient(proxyString interface{}) *http.Client {
	switch v := proxyString.(type) {
	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	default:
		return &http.Client{}
	}

}

func scrapeClientRequest(searchURL string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("Scraper received a non - 200 status code suggesting a ban")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func googleResultParsing(response *http.Response, rank int) ([]SearchResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []SearchResult{}
	sel := doc.Find("div.g")
	rank++
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" && link != "#" && !strings.HasPrefix(link, "/") {
			result := SearchResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank++
		}

	}
	return results, err
}

func GoogleScrape(searchTerm, languageCode, countryCode string, proxyString interface{}, pages, count, backOff int) ([]SearchResult, error) {
	results := []SearchResult{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTerm, languageCode, countryCode, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := scrapeClientRequest(page, proxyString)
		if err != nil {
			return nil, err
		}
		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)

		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(time.Duration(backOff) * time.Second)
	}
	return results, nil
}

func main() {
	res, err := GoogleScrape("leksyking Felix Ogundipe", "en", "com", nil, 1, 30, 10)
	if err == nil {
		for _, res := range res {
			fmt.Println(res)
		}
	}
}
