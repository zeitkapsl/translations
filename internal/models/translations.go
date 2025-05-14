package models

import "strings"

// TranslationType represents the type of translation (singular, plural with specific quantities)
type TranslationType string

const (
	TypeSingular TranslationType = "singular"
	TypePlural   TranslationType = "plural"
)

// TranslationQuantity represents the quantity for plural forms
type TranslationQuantity string

const (
	QuantityOne   TranslationQuantity = "one"
	QuantityOther TranslationQuantity = "other"
)

// Translation represents a single translation entry
type Translation struct {
	Key          string
	Comment      string
	Type         TranslationType                           // singular or plural
	Quantities   map[TranslationQuantity]string            // For plural forms (only used if Type is TypePlural)
	Translations map[string]map[TranslationQuantity]string // language -> quantity -> translation
	// For singular translations, we only use QuantityOne
}

// TranslationSet contains all translations and language metadata
type TranslationSet struct {
	Translations []Translation
	Languages    []string // All languages including regions
}

// NewTranslationSet creates a new empty translation set with default language
func NewTranslationSet() *TranslationSet {
	return &TranslationSet{
		Translations: []Translation{},
		Languages:    []string{"en"}, // Default language is English
	}
}

// AddLanguage adds a new language if it doesn't already exist
func (ts *TranslationSet) AddLanguage(lang string) bool {
	// Check if language already exists
	for _, l := range ts.Languages {
		if l == lang {
			return false // Language already exists
		}
	}

	// Add new language
	ts.Languages = append(ts.Languages, lang)
	return true
}

// HasTranslation checks if a key has a translation in the specified language
func (ts *TranslationSet) HasTranslation(key, lang string) bool {
	for _, t := range ts.Translations {
		if t.Key == key {
			if translations, ok := t.Translations[lang]; ok && len(translations) > 0 {
				return true
			}
			return false
		}
	}
	return false
}

// GetTranslation gets a translation for a key in the specified language
// Returns empty string if not found
func (ts *TranslationSet) GetTranslation(key, lang string) string {
	for _, t := range ts.Translations {
		if t.Key == key {
			if translations, ok := t.Translations[lang]; ok {
				// For singular, return the "one" quantity (or any available)
				if t.Type == TypeSingular {
					if val, ok := translations[QuantityOne]; ok {
						return val
					}
				}
				// For plural, return the "other" quantity if requested
				if val, ok := translations[QuantityOther]; ok {
					return val
				}
			}
			return ""
		}
	}
	return ""
}

// GetPluralTranslation gets plural translations for a key
func (ts *TranslationSet) GetPluralTranslation(key, lang string, quantity TranslationQuantity) string {
	for _, t := range ts.Translations {
		if t.Key == key {
			if translations, ok := t.Translations[lang]; ok {
				if val, ok := translations[quantity]; ok {
					return val
				}
			}
			return ""
		}
	}
	return ""
}

// AddOrUpdateTranslation adds a new translation or updates an existing one
func (ts *TranslationSet) AddOrUpdateTranslation(key, comment, lang string, value string) {
	ts.AddOrUpdatePluralTranslation(key, comment, lang, QuantityOne, value)
}

// AddOrUpdatePluralTranslation adds a new translation or updates an existing one with plural support
func (ts *TranslationSet) AddOrUpdatePluralTranslation(key, comment, lang string, quantity TranslationQuantity, value string) {
	// Find existing translation
	for i, t := range ts.Translations {
		if t.Key == key {
			// Update existing
			if t.Translations == nil {
				ts.Translations[i].Translations = make(map[string]map[TranslationQuantity]string)
			}

			if ts.Translations[i].Translations[lang] == nil {
				ts.Translations[i].Translations[lang] = make(map[TranslationQuantity]string)
			}

			ts.Translations[i].Translations[lang][quantity] = value

			// Update comment if provided and current is empty
			if comment != "" && ts.Translations[i].Comment == "" {
				ts.Translations[i].Comment = comment
			}

			// Update type if adding plural forms
			if quantity != QuantityOne {
				ts.Translations[i].Type = TypePlural
			}

			return
		}
	}

	// Create new
	newTrans := Translation{
		Key:     key,
		Comment: comment,
		Type:    TypeSingular, // Default to singular
		Translations: map[string]map[TranslationQuantity]string{
			lang: {
				quantity: value,
			},
		},
	}

	// If adding a plural form, update type
	if quantity != QuantityOne {
		newTrans.Type = TypePlural
	}

	ts.Translations = append(ts.Translations, newTrans)
}

// IsPluralKey determines if a key is for a plural translation
// Often keys have patterns like "key.singular" and "key.plural"
func IsPluralKey(key string) bool {
	return strings.HasSuffix(key, ".plural")
}

// NormalizeKey extracts the base key without plural suffix
func NormalizeKey(key string) string {
	if strings.HasSuffix(key, ".singular") {
		return key[:len(key)-9] // Remove .singular
	}
	if strings.HasSuffix(key, ".plural") {
		return key[:len(key)-7] // Remove .plural
	}
	return key
}
