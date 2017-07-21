package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WordType uint8

// Word type consts
const (
	NOUN WordType = iota + 1
	VERB
	ADJECTIVE
	ADVERB
	UNKNOWN
)

type Translation struct {
	Type     WordType // Word type(VERB, NOUN, ADJACTIVE, ADVERB)
	Text     string
	Category string
}

type Content struct {
	FromLang     string
	Text         string // Searched text
	Translations []Translation
	ResultCount  int // Keeps found total translation count in grabbed document
}

type Config struct {
	FromLang        string
	DisplayCount    int // Max display count
	WordTypeFilters []WordType
}

type Tureng struct {
	Config   Config
	Document *goquery.Document
}

func (trans *Translation) WordTypeShortDisplay() string {
	switch trans.Type {
	case NOUN:
		return "n."
	case VERB:
		return "v."
	case ADJECTIVE:
		return "adj."
	case ADVERB:
		return "adv."
	default:
		return ""
	}
}

var userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:51.0) Gecko/20100101 Firefox/51.0"

// Make request to tureng.com and return grabbed document
func (tureng *Tureng) getDocument(text string) (*goquery.Document, error) {

	url := fmt.Sprintf("http://www.tureng.com/en/turkish-english/%s", text)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return goquery.NewDocumentFromResponse(res)
}

// Translate given text
func (tureng *Tureng) Translate(text string) (result Content, err error) {

	doc, err := tureng.getDocument(text)
	tureng.Document = doc

	if err != nil {
		log.Fatal(err)
	}

	result = Content{Text: text}

	trElems := doc.Find("table.searchResultsTable").Eq(0).Find("tbody tr")

	// Check if translating from English or Turkish
	// Find translation language
	result.FromLang = "en"
	if trElems.Eq(0).Find(".c2").Text() == "Turkish" {
		result.FromLang = "tr"
	}

	trElems = trElems.Not(".mobile-category-row").Not("[style]")

	if trElems.Length() == 0 {
		return result, nil

	}

	// There is a header row in trElems. So subtract it from TotalCount
	result.ResultCount = trElems.Length() - 1

	trElems.Each(func(i int, s *goquery.Selection) {
		if len(result.Translations) >= tureng.Config.DisplayCount {
			return
		}

		trans := Translation{}
		trans.Category = s.Find("td").Eq(1).Text()
		en := s.Find("td[lang=en]").Find("a").Text()
		tr := s.Find("td[lang=tr]").Find("a").Text()
		if en == "" {
			return
		}

		wordTypeStr := strings.TrimSpace(s.Find("td[lang=en]").Find("i").Text())
		switch wordTypeStr {
		case "v.":
			trans.Type = VERB
		case "n.":
			trans.Type = NOUN
		case "adj.":
			trans.Type = ADJECTIVE
		case "adv.":
			trans.Type = ADVERB
		default:
			trans.Type = UNKNOWN
		}

		if result.FromLang == "en" {
			trans.Text = tr
		} else {
			trans.Text = en
		}

		for _, wordType := range tureng.Config.WordTypeFilters {
			if trans.Type == wordType {
				result.Translations = append(result.Translations, trans)
			}
		}

	})

	return
}

func (tureng *Tureng) GetSuggestions() []string {
	suggestions := []string{}
	tureng.Document.Find(".suggestion-list a").Each(func(i int, s *goquery.Selection) {
		suggestions = append(suggestions, s.Text())
	})

	return suggestions
}

func printUsage() {
	fmt.Printf("Usage: %s TEXT \n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	displayCount := flag.Int("c", 10, "Max display count")
	includeVerbsPtr := flag.Bool("v", false, "Filter verbs")
	includeNounsPtr := flag.Bool("n", false, "Filter nouns")
	includeAdverbsPtr := flag.Bool("adv", false, "Filter adverbs")
	includeAdjectivesPtr := flag.Bool("adj", false, "Filter adjectives")

	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	conf := &Config{
		DisplayCount: *displayCount,
	}

	if *includeVerbsPtr {
		conf.WordTypeFilters = append(conf.WordTypeFilters, VERB)
	}
	if *includeNounsPtr {
		conf.WordTypeFilters = append(conf.WordTypeFilters, NOUN)
	}
	if *includeAdjectivesPtr {
		conf.WordTypeFilters = append(conf.WordTypeFilters, ADJECTIVE)
	}
	if *includeAdverbsPtr {
		conf.WordTypeFilters = append(conf.WordTypeFilters, ADVERB)
	}
	if len(conf.WordTypeFilters) == 0 {
		conf.WordTypeFilters = []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN}
	}

	tureng := &Tureng{Config: *conf}
	text := strings.Join(flag.Args(), " ")
	result, err := tureng.Translate(text)

	if err != nil {
		log.Fatal(err)
	}

	if result.ResultCount == 0 {
		fmt.Printf("There is no translation found for '%s' \n", text)
		suggs := tureng.GetSuggestions()

		if len(suggs) > 0 {
			fmt.Printf("\n==== Suggestions ====\n")
			for _, item := range suggs {
				fmt.Printf("%v\n", item)
			}
		}
	} else {
		for _, trans := range result.Translations {
			if trans.Type == UNKNOWN {
				fmt.Printf("%s - %s - %s\n", trans.Category, result.Text, trans.Text)
			} else {
				fmt.Printf("%s - %s - %s (%s)\n", trans.Category, result.Text, trans.Text, trans.WordTypeShortDisplay())
			}
		}

		fmt.Printf("===== [ Total: %d ] =====\n", result.ResultCount)
	}
}
