package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

type Translation struct {
	Type string // Word type(verb, noun, adj, adv)
	Text string
}

type TranslationResponse struct {
	FromLang     string
	Text         string // Searched text
	Translations []Translation
	TotalCount   int // Keeps found total translation count in grabbed document
}

type Config struct {
	FromLang string
	MaxCount int // Max grabbing count
}

type Tureng struct {
	Config Config
}

var userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:51.0) Gecko/20100101 Firefox/51.0"

// Make request to tureng.com and return scraped document
func (tureng Tureng) getDocument(text string) (*goquery.Document, error) {

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

//
func (tureng Tureng) translate(text string) (result TranslationResponse, err error) {

	doc, err := tureng.getDocument(text)

	if err != nil {
		log.Fatal(err)
	}

	result = TranslationResponse{Text: text}

	trElems := doc.Find("table.searchResultsTable").Eq(0).Find("tbody tr")

	// Check if translating from English or Turkish
	// Find translation language
	result.FromLang = "en"
	if trElems.Eq(0).Find(".c2").Text() == "Turkish" {
		result.FromLang = "tr"
	}

	trElems = trElems.Not(".mobile-category-row").Not("[style]")
	result.TotalCount = trElems.Length()
	trElems.Each(func(i int, s *goquery.Selection) {
		if len(result.Translations) > tureng.Config.MaxCount {
			return
		}

		trans := Translation{}
		en := s.Find("td[lang=en]").Find("a").Text()
		tr := s.Find("td[lang=tr]").Find("a").Text()

		if en == "" {
			return
		}

		trans.Type = strings.TrimSpace(s.Find("td[lang=en]").Find("i").Text())
		if result.FromLang == "en" {
			trans.Text = tr
		} else {
			trans.Text = en
		}
		result.Translations = append(result.Translations, trans)
	})

	return
}

func main() {
	text := os.Args[1]

	if text == "" {
		os.Exit(1)
	}

	conf := &Config{MaxCount: 3}
	tureng := &Tureng{Config: *conf}

	result, err := tureng.translate(text)
	if err != nil {
		log.Fatal(err)
	}

	for _, trans := range result.Translations {
		if trans.Type != "" {
			fmt.Printf("%s - %s (%s)\n", result.Text, trans.Text, trans.Type)
		} else {
			fmt.Printf("%s - %s\n", trans.Text)
		}
	}

	fmt.Printf("===========\nTotal: %d\n", result.TotalCount)
}
