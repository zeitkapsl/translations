package models

import (
	"fmt"
	"sort"
	"strings"
)

// TranslationRow represents a single translation entry
type TranslationRow struct {
	Key     string
	Comment string
	Type    string // "singular" or "plural"
	Values  map[string]TranslationValues
}

// TranslationValues holds translation values for different quantities
type TranslationValues struct {
	One   string // For singular or plural "one"
	Other string // For plural "other"
}

// TranslationManager manages all translations and language metadata
type TranslationManager struct {
	Translations []TranslationRow
	Languages    []string
	BasePath     string
}

// NewTranslationManager creates a new translation manager
func NewTranslationManager(basePath string) *TranslationManager {
	return &TranslationManager{
		Translations: make([]TranslationRow, 0),
		Languages:    make([]string, 0),
		BasePath:     basePath,
	}
}

// EnsureLanguage adds a language if it doesn't exist, with normalization
func (tm *TranslationManager) EnsureLanguage(lang string) {
	normalizedLang := tm.normalizeLanguageCode(lang)

	for _, l := range tm.Languages {
		if l == normalizedLang {
			return
		}
	}
	tm.Languages = append(tm.Languages, normalizedLang)
	sort.Strings(tm.Languages)
}

// normalizeLanguageCode normalizes common language code variations
func (tm *TranslationManager) normalizeLanguageCode(lang string) string {
	switch lang {
	case "en-US":
		return "en" // Treat en-US as en
	case "de-DE":
		return "de" // Treat de-DE as de
	default:
		return lang
	}
}

// SetTranslation adds or updates a singular translation
func (tm *TranslationManager) SetTranslation(key, lang, value, comment, transType string) {
	normalizedLang := tm.normalizeLanguageCode(lang)

	// Find existing row
	for i, row := range tm.Translations {
		if row.Key == key {
			if row.Values == nil {
				row.Values = make(map[string]TranslationValues)
			}
			values := row.Values[normalizedLang]
			values.One = value
			row.Values[normalizedLang] = values
			if comment != "" && row.Comment == "" {
				row.Comment = comment
			}
			tm.Translations[i] = row
			return
		}
	}

	// Create new row
	row := TranslationRow{
		Key:     key,
		Comment: comment,
		Type:    transType,
		Values:  map[string]TranslationValues{normalizedLang: {One: value}},
	}
	tm.Translations = append(tm.Translations, row)
}

// SetPluralTranslation adds or updates a plural translation
func (tm *TranslationManager) SetPluralTranslation(key, lang, oneValue, otherValue, comment string) {
	normalizedLang := tm.normalizeLanguageCode(lang)

	// Find existing row
	for i, row := range tm.Translations {
		if row.Key == key {
			if row.Values == nil {
				row.Values = make(map[string]TranslationValues)
			}
			row.Values[normalizedLang] = TranslationValues{One: oneValue, Other: otherValue}
			if comment != "" && row.Comment == "" {
				row.Comment = comment
			}
			row.Type = "plural"
			tm.Translations[i] = row
			return
		}
	}

	// Create new row
	row := TranslationRow{
		Key:     key,
		Comment: comment,
		Type:    "plural",
		Values:  map[string]TranslationValues{normalizedLang: {One: oneValue, Other: otherValue}},
	}
	tm.Translations = append(tm.Translations, row)
}

// SetPluralSingular sets only the singular form of a plural translation
func (tm *TranslationManager) SetPluralSingular(key, lang, value, comment string) {
	normalizedLang := tm.normalizeLanguageCode(lang)

	// Find existing row
	for i, row := range tm.Translations {
		if row.Key == key {
			if row.Values == nil {
				row.Values = make(map[string]TranslationValues)
			}
			values := row.Values[normalizedLang]
			values.One = value
			row.Values[normalizedLang] = values
			row.Type = "plural"
			if comment != "" && row.Comment == "" {
				row.Comment = comment
			}
			tm.Translations[i] = row
			return
		}
	}

	// Create new row
	row := TranslationRow{
		Key:     key,
		Comment: comment,
		Type:    "plural",
		Values:  map[string]TranslationValues{normalizedLang: {One: value}},
	}
	tm.Translations = append(tm.Translations, row)
}

// SetPluralOther sets only the other form of a plural translation
func (tm *TranslationManager) SetPluralOther(key, lang, value, comment string) {
	normalizedLang := tm.normalizeLanguageCode(lang)

	// Find existing row
	for i, row := range tm.Translations {
		if row.Key == key {
			if row.Values == nil {
				row.Values = make(map[string]TranslationValues)
			}
			values := row.Values[normalizedLang]
			values.Other = value
			row.Values[normalizedLang] = values
			row.Type = "plural"
			tm.Translations[i] = row
			return
		}
	}

	// Create new row
	row := TranslationRow{
		Key:     key,
		Comment: comment,
		Type:    "plural",
		Values:  map[string]TranslationValues{normalizedLang: {Other: value}},
	}
	tm.Translations = append(tm.Translations, row)
}

// GetTranslationWithFallback gets translation with regional fallback support
func (tm *TranslationManager) GetTranslationWithFallback(row TranslationRow, lang string) TranslationValues {
	// Direct match
	if values, exists := row.Values[lang]; exists {
		return values
	}

	// Regional fallback (de-AT -> de)
	if strings.Contains(lang, "-") {
		baseLang := strings.Split(lang, "-")[0]
		if values, exists := row.Values[baseLang]; exists {
			return values
		}
	}

	// English fallback
	if values, exists := row.Values["en"]; exists {
		return values
	}

	return TranslationValues{}
}

// AddLanguage adds a new language
func (tm *TranslationManager) AddLanguage(lang string) {
	tm.EnsureLanguage(lang)
	fmt.Printf("Added language: %s\n", lang)
}

// AddRegion adds a new region
func (tm *TranslationManager) AddRegion(region string) {
	tm.EnsureLanguage(region)
	fmt.Printf("Added region: %s\n", region)
}

// GetStats displays basic translation statistics
func (tm *TranslationManager) GetStats() {
	fmt.Printf("Translation Statistics:\n")
	fmt.Printf("Total languages: %d\n", len(tm.Languages))
	fmt.Printf("Total strings: %d\n", len(tm.Translations))
	fmt.Printf("\nLanguages: %s\n", strings.Join(tm.Languages, ", "))

	// Calculate missing translations
	if len(tm.Translations) == 0 {
		return
	}

	fmt.Printf("\nMissing translations per language:\n")
	for _, lang := range tm.Languages {
		if lang == "en" {
			continue // Skip English as source
		}

		missing := 0
		total := 0

		for _, row := range tm.Translations {
			// Check if there's an English version
			if enValues, ok := row.Values["en"]; ok && (enValues.One != "" || enValues.Other != "") {
				total++

				values := tm.GetTranslationWithFallback(row, lang)
				if row.Type == "plural" {
					if values.One == "" || values.Other == "" {
						missing++
					}
				} else {
					if values.One == "" {
						missing++
					}
				}
			}
		}

		if total > 0 {
			percentage := float64(total-missing) / float64(total) * 100
			fmt.Printf("  %s: %d missing (%.1f%% complete)\n", lang, missing, percentage)
		}
	}
}

// GetDetailedStats displays detailed translation statistics with source analysis
func (tm *TranslationManager) GetDetailedStats() {
	fmt.Printf("Detailed Translation Statistics:\n")
	fmt.Printf("================================\n")
	fmt.Printf("Total languages: %d\n", len(tm.Languages))
	fmt.Printf("Total translation keys: %d\n", len(tm.Translations))
	fmt.Printf("Languages: %s\n\n", strings.Join(tm.Languages, ", "))

	// Analyze source coverage
	sourceCoverage := make(map[string]int)
	totalTranslatableKeys := 0

	for _, row := range tm.Translations {
		hasAnyContent := false
		for lang, values := range row.Values {
			if values.One != "" || values.Other != "" {
				sourceCoverage[lang]++
				hasAnyContent = true
			}
		}
		if hasAnyContent {
			totalTranslatableKeys++
		}
	}

	fmt.Printf("Source Content Analysis:\n")
	fmt.Printf("========================\n")
	fmt.Printf("Keys with translatable content: %d\n", totalTranslatableKeys)
	fmt.Printf("Source language coverage:\n")
	for _, lang := range tm.Languages {
		if count, exists := sourceCoverage[lang]; exists {
			percentage := float64(count) / float64(totalTranslatableKeys) * 100
			fmt.Printf("  %s: %d keys (%.1f%% of translatable content)\n", lang, count, percentage)
		} else {
			fmt.Printf("  %s: 0 keys (0.0%% of translatable content)\n", lang)
		}
	}

	fmt.Printf("\nTranslation Completeness per Language:\n")
	fmt.Printf("======================================\n")

	for _, lang := range tm.Languages {
		if lang == "en" || lang == "de" {
			continue // Skip source languages
		}

		missing := 0
		total := 0
		availableFromEn := 0
		availableFromDe := 0
		availableFromOther := 0

		for _, row := range tm.Translations {
			// Check what sources are available
			enValues, hasEn := row.Values["en"]
			deValues, hasDe := row.Values["de"]

			hasEnContent := hasEn && (enValues.One != "" || enValues.Other != "")
			hasDeContent := hasDe && (deValues.One != "" || deValues.Other != "")

			// Count what's available as source material
			if hasEnContent {
				availableFromEn++
			}
			if hasDeContent {
				availableFromDe++
			}

			// Check if ANY language has content for this key
			hasAnySource := hasEnContent || hasDeContent
			for otherLang, otherValues := range row.Values {
				if otherLang != lang && otherLang != "en" && otherLang != "de" {
					if otherValues.One != "" || otherValues.Other != "" {
						hasAnySource = true
						if !hasEnContent && !hasDeContent {
							availableFromOther++
						}
						break
					}
				}
			}

			if !hasAnySource {
				continue // Skip keys with no source content at all
			}

			total++

			// Check if target language has translation
			targetValues := tm.GetTranslationWithFallback(row, lang)
			if row.Type == "plural" {
				if targetValues.One == "" || targetValues.Other == "" {
					missing++
				}
			} else {
				if targetValues.One == "" {
					missing++
				}
			}
		}

		if total > 0 {
			percentage := float64(total-missing) / float64(total) * 100
			fmt.Printf("%s: %d missing of %d total (%.1f%% complete)\n", lang, missing, total, percentage)
			fmt.Printf("  Available from English: %d keys\n", availableFromEn)
			fmt.Printf("  Available from German: %d keys\n", availableFromDe)
			fmt.Printf("  Available from other languages: %d keys\n", availableFromOther)
			fmt.Printf("  Potential max translations: %d keys\n\n", availableFromEn+availableFromDe+availableFromOther)
		}
	}

	// Show some example missing translations
	fmt.Printf("Example Missing Translations (showing first 5):\n")
	fmt.Printf("==============================================\n")

	for _, lang := range tm.Languages {
		if lang == "en" || lang == "de" {
			continue
		}

		fmt.Printf("\nMissing in %s:\n", lang)
		count := 0
		for _, row := range tm.Translations {
			if count >= 5 {
				break
			}

			// Check if we have source content
			enValues, hasEn := row.Values["en"]
			deValues, hasDe := row.Values["de"]
			hasEnContent := hasEn && (enValues.One != "" || enValues.Other != "")
			hasDeContent := hasDe && (deValues.One != "" || deValues.Other != "")

			if !hasEnContent && !hasDeContent {
				continue
			}

			// Check if target is missing
			targetValues := tm.GetTranslationWithFallback(row, lang)
			isMissing := false
			if row.Type == "plural" {
				isMissing = targetValues.One == "" || targetValues.Other == ""
			} else {
				isMissing = targetValues.One == ""
			}

			if isMissing {
				sourceText := ""
				sourceLang := ""
				if hasEnContent {
					sourceText = enValues.One
					if sourceText == "" {
						sourceText = enValues.Other
					}
					sourceLang = "en"
				} else if hasDeContent {
					sourceText = deValues.One
					if sourceText == "" {
						sourceText = deValues.Other
					}
					sourceLang = "de"
				}

				if sourceText != "" {
					fmt.Printf("  %s (%s): \"%s\"\n", row.Key, sourceLang, sourceText)
					count++
				}
			}
		}

		if count == 0 {
			fmt.Printf("  No missing translations!\n")
		}
	}
}

// CleanupLanguages merges redundant language variants
func (tm *TranslationManager) CleanupLanguages() int {
	mergeCount := 0

	// Define merge rules
	mergeRules := map[string]string{
		"en-US": "en",
		"de-DE": "de",
	}

	// Process each translation row
	for i, row := range tm.Translations {
		newValues := make(map[string]TranslationValues)

		for lang, values := range row.Values {
			targetLang := lang
			if mergeTo, shouldMerge := mergeRules[lang]; shouldMerge {
				targetLang = mergeTo
				mergeCount++
			}

			// Merge values if target already exists
			if existing, exists := newValues[targetLang]; exists {
				// Prefer non-empty values
				merged := existing
				if merged.One == "" && values.One != "" {
					merged.One = values.One
				}
				if merged.Other == "" && values.Other != "" {
					merged.Other = values.Other
				}
				newValues[targetLang] = merged
			} else {
				newValues[targetLang] = values
			}
		}

		tm.Translations[i].Values = newValues
	}

	// Update language list
	newLanguages := []string{}
	for _, lang := range tm.Languages {
		if mergeTo, shouldMerge := mergeRules[lang]; shouldMerge {
			// Check if target doesn't already exist
			found := false
			for _, existing := range newLanguages {
				if existing == mergeTo {
					found = true
					break
				}
			}
			if !found {
				newLanguages = append(newLanguages, mergeTo)
			}
		} else {
			newLanguages = append(newLanguages, lang)
		}
	}

	tm.Languages = newLanguages
	sort.Strings(tm.Languages)

	return mergeCount
}

// Helper function to check if slice contains item
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
