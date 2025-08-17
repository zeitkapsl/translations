package main

import (
	"sort"
	"strings"
)

// TranslationRow represents a single translation entry
type TranslationRow struct {
	App     string
	Key     string
	Comment string
	Values  map[string]string
}

func (tr TranslationRow) IsPlural() bool {
	if strings.HasSuffix(tr.Key, ".singular") {
		return true
	}
	if strings.HasSuffix(tr.Key, ".plural") {
		return true
	}
	return false
}

func (tr TranslationRow) GetSingularKey() string {
	if strings.HasSuffix(tr.Key, ".singular") {
		return tr.Key[:len(tr.Key)-len(".singular")]
	} else if strings.HasSuffix(tr.Key, ".plural") {
		return tr.Key[:len(tr.Key)-len(".plural")]
	}
	return tr.Key
}

// TranslationValues holds translation values for different quantities
type TranslationValues struct {
	One   map[string]string // For singular or plural "one"
	Other map[string]string // For plural "other"
}

// Translations manages all translations and language metadata
type Translations struct {
	Translations []TranslationRow
	Languages    []string
	BasePath     string
}

// Translations creates a new translation manager
func NewTranslations(basePath string) *Translations {
	return &Translations{
		Translations: make([]TranslationRow, 0),
		Languages:    make([]string, 0),
		BasePath:     basePath,
	}
}

// EnsureLanguage adds a language if it doesn't exist, with normalization
func (tm *Translations) EnsureLanguage(lang string) {

	for _, l := range tm.Languages {
		if l == lang {
			return
		}
	}
	tm.Languages = append(tm.Languages, lang)
	sort.Strings(tm.Languages)
}

// SetTranslation adds or updates a singular translation
func (tm *Translations) SetTranslation(app, key, lang, value, comment string) {

	// Find existing row
	for i, row := range tm.Translations {
		if row.Key == key && row.App == app {
			if row.Values == nil {
				row.Values = make(map[string]string)
			}
			row.Values[lang] = value
			if comment != "" && row.Comment == "" {
				row.Comment = comment
			}
			tm.Translations[i] = row
			return
		}
	}

	// Create new row
	row := TranslationRow{
		App:     app,
		Key:     key,
		Comment: comment,
		Values:  map[string]string{lang: value},
	}
	tm.Translations = append(tm.Translations, row)
}

// SetPluralTranslation adds or updates a plural translation
func (tm *Translations) SetPluralTranslation(app, key, lang, oneValue, otherValue, comment string) {
	tm.SetTranslation(app, key+".singular", lang, oneValue, comment)
	tm.SetTranslation(app, key+".plural", lang, otherValue, comment)
}

// SetPluralSingular sets only the singular form of a plural translation
func (tm *Translations) SetPluralSingular(app, key, lang, value, comment string) {
	tm.SetTranslation(app, key+".singular", lang, value, comment)
}

// SetPluralOther sets only the other form of a plural translation
func (tm *Translations) SetPluralOther(app, key, lang, value, comment string) {
	tm.SetTranslation(app, key+".plural", lang, value, comment)
}

// AddLanguage adds a new language
func (tm *Translations) AddLanguage(lang string) {
	tm.EnsureLanguage(lang)
}

// AddRegion adds a new region
func (tm *Translations) AddRegion(region string) {
	tm.EnsureLanguage(region)
}

func (tm *Translations) GetRow(app, key string) *TranslationRow {

	for _, v := range tm.Translations {
		if v.App == app && v.Key == key {
			return &v
		}
	}
	return nil
}

func (tm *Translations) GetPlural(app, key string) *TranslationValues {
	singular := tm.GetRow(app, key+".singular")
	plural := tm.GetRow(app, key+".plural")
	if singular != nil && plural != nil {
		return &TranslationValues{
			One:   singular.Values,
			Other: plural.Values,
		}
	}
	return nil
}

func (tm *Translations) GetTranslationsForApp(app string) []TranslationRow {
	result := make([]TranslationRow, 0)
	for _, v := range tm.Translations {
		if v.App == app {
			result = append(result, v)
		}
	}
	return result
}

func (tm *Translations) Sort() {
	// Sort translations by App + Key
	sort.Slice(tm.Translations, func(i, j int) bool {
		if tm.Translations[i].App == tm.Translations[j].App {
			return tm.Translations[i].Key < tm.Translations[j].Key
		}
		return tm.Translations[i].App < tm.Translations[j].App
	})
}



