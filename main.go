package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

// Make request to tureng.com and return document
func getDocument(text string) (*goquery.Document, error) {

	userAgent := "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:51.0) Gecko/20100101 Firefox/51.0"
	url := fmt.Sprintf("http://www.tureng.com/en/turkish-english/%s", text)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(res)
}

// Translate given text and print
func translate(text string) {
	doc, err := getDocument(text)

	if err != nil {
		log.Fatal(err)
	}

	trElems := doc.Find("table.searchResultsTable").Eq(0).Find("tbody tr")

	// Check if translating from English or Turkish
	// Find translation language
	isFromEng := true
	if trElems.Eq(0).Find(".c2").Text() == "Turkish" {
		isFromEng = false
	}

	trElems = trElems.Not(".mobile-category-row").Not("[style]").Slice(0, 10)
	trElems.Each(func(i int, s *goquery.Selection) {
		en := s.Find("td[lang=en]").Find("a").Text()
		tr := s.Find("td[lang=tr]").Find("a").Text()

		if en == "" {
			return
		}

		wordType := strings.TrimSpace(s.Find("td[lang=en]").Find("i").Text())
		if wordType == "" {
			if isFromEng {
				fmt.Printf("%s - %s\n", en, tr)
			} else {
				fmt.Printf("%s - %s\n", tr, en)
			}
		} else {
			if isFromEng {
				fmt.Printf("%s(%s) - %s\n", en, wordType, tr)
			} else {
				fmt.Printf("%s - %s(%s)\n", tr, en, wordType)
			}
		}
	})
}

func main() {
	text := os.Args[1]

	if text == "" {
		os.Exit(1)
	}

	translate(text)
}
