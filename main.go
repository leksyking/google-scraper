package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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
