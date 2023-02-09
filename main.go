package main

import (
	"fmt"
	"math/rand"
	"strings"
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

func GoogleScrape(searchTerm, languageCode, countryCode string, pages, count int) ([]SearchResult, error) {
	results := []SearchResult{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTerm, languageCode, countryCode, pages, count)
}

func main() {
	res, err := GoogleScrape("leksyking Felix Ogundipe", "en", "com", 1, 30)
	if err == nil {
		for _, res := range res {
			fmt.Println(res)
		}
	}
}
