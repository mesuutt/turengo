package main

import (
	"testing"
)

// TODO: @mesut. Test with mocking.

var text = "brave"

func TestDisplayCount(t *testing.T) {
	displayCount := 2

	conf := &Config{
		DisplayCount: displayCount,
	}

	conf.WordTypeFilters = []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN}

	tureng := &Tureng{Config: *conf}
	result, _ := tureng.Translate(text)

	if len(result.TranslationGroups[0].Translations) != displayCount {
		t.Errorf("Result count should not be greater than %d. %d", displayCount, len(result.TranslationGroups[0].Translations))
	}
}

func TestTranslate(t *testing.T) {
	displayCount := 10

	conf := &Config{
		DisplayCount: displayCount,
	}

	conf.WordTypeFilters = []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN}

	tureng := &Tureng{Config: *conf}
	result, _ := tureng.Translate(text)

	if len(result.TranslationGroups[0].Translations) == 0 {
		t.Errorf("Translation count should not be greater than %d. %d", displayCount, len(result.TranslationGroups[0].Translations))
	}

	if len(result.TranslationGroups[0].Translations) == 0 {
		t.Errorf("Translation count for '%s' should be greater than 0", text)
	}
}

func TestGettingMeaningOfWordWithOtherTerms(t *testing.T) {
	conf := &Config{
		DisplayCount:    100,
		WordTypeFilters: []WordType{VERB}, // Get only verbs
	}

	tureng := &Tureng{Config: *conf}
	result, _ := tureng.Translate(text)

	if len(result.TranslationGroups[1].Translations) == 0 {
		t.Errorf("Other meanings of '%s' should be exist. But not found.", text)
	}
}

func TestWordTypeFiltering(t *testing.T) {

	conf := &Config{
		DisplayCount:    100,
		WordTypeFilters: []WordType{VERB}, // Get only verbs
	}

	tureng := &Tureng{Config: *conf}
	result, _ := tureng.Translate(text)

	if len(result.TranslationGroups[0].Translations) != 6 {
		t.Errorf("Translation count should equal to 6")
	}

	conf = &Config{
		DisplayCount:    100,
		WordTypeFilters: []WordType{ADVERB}, // Get only adverbs
	}
	tureng = &Tureng{Config: *conf}
	result, _ = tureng.Translate(text)

	if len(result.TranslationGroups[0].Translations) != 0 {
		t.Errorf("Translation count should equal to 0")
	}
}

func TestGettingSuggestions(t *testing.T) {

	conf := &Config{
		DisplayCount:    100,
		WordTypeFilters: []WordType{NOUN},
	}

	tureng := &Tureng{Config: *conf}
	result, _ := tureng.Translate("happyoooo")

	if len(result.TranslationGroups) > 0 {
		t.Errorf("Translation for 'happyoooo' should not be found.")
	}

	suggs := tureng.GetSuggestions()
	if len(suggs) == 0 {
		t.Errorf("Should be at least one suggestion for 'happyoooo'")
	}
}
