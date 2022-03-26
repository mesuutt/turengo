package main

import (
	"testing"
)

// TODO: mock website

func TestDisplayCount(t *testing.T) {
	displayCount := 2

	flags := &Flags{
		DisplayCount: displayCount,
		TypeFilters:  []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN},
	}

	result, _ := Translate("brave", flags)

	if len(result.Translations) != displayCount {
		t.Errorf("Result count should not be greater than %d. %d", displayCount, len(result.Translations))
	}
}

func TestTranslate(t *testing.T) {
	displayCount := 10

	flags := &Flags{
		DisplayCount: displayCount,
		TypeFilters:  []WordType{VERB, NOUN, ADJECTIVE, ADVERB, UNKNOWN},
	}

	result, _ := Translate("brave", flags)

	if len(result.Translations) == 0 {
		t.Errorf("Translation count should not be greater than %d. %d", displayCount, len(result.Translations))
	}

	if len(result.Translations) == 0 {
		t.Error("Translation count should be greater than 0")
	}
}

func TestWordTypeFiltering(t *testing.T) {
	t.Run("get only verbs", func(t *testing.T) {
		conf := &Flags{
			DisplayCount: 100,
			TypeFilters:  []WordType{VERB}, // Get only verbs
		}

		result, _ := Translate("brave", conf)

		for _, item := range result.Translations {
			if item.Type != VERB {
				t.Fatal("Translation count should equal to 6")
			}
		}
	})

	t.Run("get only adverbs", func(t *testing.T) {
		flags := &Flags{
			DisplayCount: 100,
			TypeFilters:  []WordType{ADVERB}, // Get only adverbs
		}

		result, _ := Translate("brave", flags)

		if len(result.Translations) != 0 {
			t.Errorf("Translation count should equal to 0")
		}
	})
}

func TestGettingSuggestions(t *testing.T) {
	flags := &Flags{
		DisplayCount: 100,
		TypeFilters:  []WordType{NOUN},
	}

	result, _ := Translate("happyoooo", flags)

	if len(result.Translations) > 0 {
		t.Errorf("Translation count for 'happyoooo' should be 0.")
	}

	if len(result.Suggestions) == 0 {
		t.Errorf("Should be at least one suggestion for 'happyoooo'")
	}
}
