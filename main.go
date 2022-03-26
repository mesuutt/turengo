package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ryanuber/columnize"
)

const UserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:51.0) Gecko/20100101 Firefox/51.0"

type WordType uint8

const (
	NOUN WordType = iota + 1
	VERB
	ADJECTIVE
	ADVERB
	UNKNOWN
)

type Flags struct {
	DisplayCount int // Max display count
	TypeFilters  []WordType
}

type Translation struct {
	Type     WordType // Word type(VERB, NOUN, ADJECTIVE, ADVERB)
	Text     string   // Translation text
	Meaning  string   // Meaning of text
	Category string
}

func (t *Translation) WordTypeDisplay() string {
	switch t.Type {
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

type PageContent struct {
	FromLang     string
	Text         string // Searched text
	ResultCount  int    // Keeps found total translation count in grabbed document
	Translations []Translation
	Suggestions  []string
}

func main() {
	flags := getFlags()
	pageContent, err := Translate(strings.Join(flag.Args(), " "), flags)
	if err != nil {
		log.Fatal(err)
	}

	printTranslations(pageContent)
}

// Translate translates given text
func Translate(text string, flags *Flags) (*PageContent, error) {
	doc, err := getDocument(text)
	if err != nil {
		return nil, err
	}

	result := &PageContent{Text: text}

	// Find translation language
	result.FromLang = "en"

	if doc.Find("table.searchResultsTable tbody tr").Eq(0).Find(".c2").Text() == "Türkçe" {
		result.FromLang = "tr"
	}

	tables := doc.Find("table.searchResultsTable")

	// There is a header row in each translation group table. So subtract it from ResultCount
	result.ResultCount = tables.Find("tbody tr").Not(".mobile-category-row").Not("[style]").Length() - tables.Length()

	if result.ResultCount <= 0 {
		fmt.Printf("There is no translation found for '%s' \n", text)
		result.Suggestions = extractSuggestions(doc)
		return result, nil
	}

	tables.Each(func(_ int, tableSel *goquery.Selection) {
		trElems := tableSel.Find("tbody tr").Not(".mobile-category-row").Not("[style]")
		trElems.EachWithBreak(func(i int, trSel *goquery.Selection) bool {
			// Ignore first row, because it is header row
			if i == 0 {
				return true // continue
			}

			if len(result.Translations) == flags.DisplayCount {
				return false
			}

			trans := Translation{}
			trans.Category = trSel.Find("td").Eq(1).Text()
			en := trSel.Find("td[lang=en]").Find("a").Text()
			tr := trSel.Find("td[lang=tr]").Find("a").Text()

			wordTypeStr := strings.TrimSpace(trSel.Find("td[lang=en]").Find("i").Text())
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
				trans.Meaning = tr
				trans.Text = en
			} else {
				trans.Meaning = en
				trans.Text = tr
			}

			for _, wordType := range flags.TypeFilters {
				if trans.Type == wordType {
					result.Translations = append(result.Translations, trans)
					break
				}
			}

			return true
		})

	})

	return result, nil
}

func getDocument(text string) (*goquery.Document, error) {
	url := fmt.Sprintf("http://www.tureng.com/en/turkish-english/%s", text)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return goquery.NewDocumentFromReader(res.Body)
}

func extractSuggestions(doc *goquery.Document) []string {
	var suggestions []string
	doc.Find(".suggestion-list a").Each(func(i int, s *goquery.Selection) {
		suggestions = append(suggestions, s.Text())
	})

	return suggestions
}

func printUsage() {
	fmt.Printf("Usage: %s TEXT \n", os.Args[0])
	flag.PrintDefaults()
}

func printTranslations(pageContent *PageContent) {
	if pageContent.ResultCount > 0 {
		var lines []string
		for _, item := range pageContent.Translations {
			if item.Type == UNKNOWN {
				lines = append(lines, fmt.Sprintf("%s | %s | %s\n", item.Category, item.Text, item.Meaning))
			} else {
				lines = append(lines, fmt.Sprintf("%s | %s | %s (%s)\n", item.Category, item.Text, item.Meaning, item.WordTypeDisplay()))
			}
		}

		fmt.Println(columnize.SimpleFormat(lines))
		fmt.Printf("===== [ Total: %d ] =====\n", pageContent.ResultCount)
		return
	}

	fmt.Printf("There is no translation found for '%s' \n", pageContent.Text)
	if len(pageContent.Suggestions) > 0 {
		fmt.Printf("\n==== Suggestions ====\n")
		for _, item := range pageContent.Suggestions {
			fmt.Printf("%v\n", item)
		}
	}

}

func getFlags() *Flags {
	defaultDisplayCount := 10

	displayCount := flag.Int("c", defaultDisplayCount, "Max display count")
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

	flags := &Flags{DisplayCount: *displayCount}

	// Read displayCount from ENV if flag not specified and env var exists
	if *displayCount == defaultDisplayCount {
		if dc := os.Getenv("TURENGO_DEFAULT_DISPLAY_COUNT"); dc != "" {
			i, _ := strconv.Atoi(dc)
			flags.DisplayCount = i
		}
	}

	if *includeVerbsPtr {
		flags.TypeFilters = append(flags.TypeFilters, VERB)
	}
	if *includeNounsPtr {
		flags.TypeFilters = append(flags.TypeFilters, NOUN)
	}
	if *includeAdjectivesPtr {
		flags.TypeFilters = append(flags.TypeFilters, ADJECTIVE)
	}
	if *includeAdverbsPtr {
		flags.TypeFilters = append(flags.TypeFilters, ADVERB)
	}
	if len(flags.TypeFilters) == 0 {
		flags.TypeFilters = []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN}
	}

	return flags
}
