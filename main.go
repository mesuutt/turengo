package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DefaultDisplayCount = 10
const DefaultDisplayCountEnv = "TURENGO_DEFAULT_DISPLAY_COUNT"

type Flags struct {
	DisplayCount int // Max display count
	TypeFilters  []WordType
}

func main() {
	flags := getFlags()
	pageContent, err := Translate(strings.Join(flag.Args(), " "), flags)
	if err != nil {
		log.Fatal(err)
	}

	pageContent.PrintAsTable()
}

func getFlags() *Flags {
	displayCount := flag.Int("c", DefaultDisplayCount, "Max display count")
	verbFlag := flag.Bool("v", false, "Filter verbs")
	nounFlag := flag.Bool("n", false, "Filter nouns")
	adverbFlag := flag.Bool("adv", false, "Filter adverbs")
	adjectiveFlag := flag.Bool("adj", false, "Filter adjectives")

	flag.Usage = func() {
		fmt.Printf("Usage: %s TEXT \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	flags := &Flags{DisplayCount: *displayCount}

	// Read displayCount from ENV if flag not specified and env var exists
	if *displayCount == DefaultDisplayCount {
		if dc := os.Getenv(DefaultDisplayCountEnv); dc != "" {
			i, _ := strconv.Atoi(dc)
			flags.DisplayCount = i
		}
	}

	if *verbFlag {
		flags.TypeFilters = append(flags.TypeFilters, VERB)
	}
	if *nounFlag {
		flags.TypeFilters = append(flags.TypeFilters, NOUN)
	}
	if *adjectiveFlag {
		flags.TypeFilters = append(flags.TypeFilters, ADJECTIVE)
	}
	if *adverbFlag {
		flags.TypeFilters = append(flags.TypeFilters, ADVERB)
	}
	if len(flags.TypeFilters) == 0 {
		flags.TypeFilters = []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN}
	}

	return flags
}
