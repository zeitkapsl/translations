package ios

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// XCStringsFile represents the structure of an .xcstrings file
type XCStringsFile struct {
	SourceLanguage string                    `json:"sourceLanguage"`
	Strings        map[string]XCStringsEntry `json:"strings"`
	Version        string                    `json:"version,omitempty"`
}

// XCStringsEntry represents a single string entry in an .xcstrings file
type XCStringsEntry struct {
	Comment       string                           `json:"comment,omitempty"`
	Localizations map[string]XCStringsLocalization `json:"localizations"`
}

// XCStringsLocalization represents a localization in an .xcstrings file
type XCStringsLocalization struct {
	StringUnit XCStringsUnit `json:"stringUnit"`
}

// XCStringsUnit represents a string unit in an .xcstrings file
type XCStringsUnit struct {
	State string `json:"state"`
	Value string `json:"value"`
}

// ImportFromXCStrings imports translations from iOS .xcstrings files
func ImportFromXCStrings(ts *models.TranslationSet, directory string) error {
	if directory == "" {
		directory = "."
	}

	// Find .xcstrings files
	files, err := filepath.Glob(filepath.Join(directory, "*.xcstrings"))
	if err != nil {
		return fmt.Errorf("error finding .xcstrings files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no .xcstrings files found in %s", directory)
	}

	importCount := 0

	for _, file := range files {
		fmt.Printf("Importing from %s...\n", file)

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		var xcstrings XCStringsFile
		if err := json.Unmarshal(data, &xcstrings); err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		// Add source language if not already present
		ts.AddLanguage(xcstrings.SourceLanguage)

		// Process strings
		for key, entry := range xcstrings.Strings {
			comment := entry.Comment

			// Process localizations
			for lang, localization := range entry.Localizations {
				// Only add translations that are in "translated" state and have a value
				if localization.StringUnit.State == "translated" && localization.StringUnit.Value != "" {
					// Check if it's a plural key
					if models.IsPluralKey(key) {
						baseKey := models.NormalizeKey(key)
						if strings.HasSuffix(key, ".singular") {
							ts.AddOrUpdatePluralTranslation(baseKey, comment, lang, models.QuantityOne, localization.StringUnit.Value)
						} else if strings.HasSuffix(key, ".plural") {
							ts.AddOrUpdatePluralTranslation(baseKey, comment, lang, models.QuantityOther, localization.StringUnit.Value)
						} else {
							ts.AddOrUpdateTranslation(key, comment, lang, localization.StringUnit.Value)
						}
					} else {
						ts.AddOrUpdateTranslation(key, comment, lang, localization.StringUnit.Value)
					}

					// Add language if not already present
					ts.AddLanguage(lang)

					importCount++
				}
			}
		}
	}

	fmt.Printf("Imported %d translations from iOS .xcstrings files\n", importCount)
	return nil
}

// ExportToXCStrings exports translations to iOS .xcstrings format
func ExportToXCStrings(ts *models.TranslationSet, outputFile string) error {
	if outputFile == "" {
		outputFile = "Localizable.xcstrings"
	}

	// Create xcstrings structure
	xcstrings := XCStringsFile{
		Version:        "1.0",
		SourceLanguage: "en", // Assuming English is the source language
		Strings:        make(map[string]XCStringsEntry),
	}

	// Process all translations
	for _, trans := range ts.Translations {
		var key string

		// For plural forms, we need special handling
		if trans.Type == models.TypePlural {
			// iOS handles plurals differently in SwiftUI vs older frameworks
			// We'll create separate entries for singular and plural forms

			// Singular form
			if oneTrans, ok := trans.Translations["en"][models.QuantityOne]; ok && oneTrans != "" {
				singularKey := trans.Key + ".singular"
				singularEntry := XCStringsEntry{
					Comment:       trans.Comment + " (singular)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for singular
				for lang, quantities := range trans.Translations {
					if val, ok := quantities[models.QuantityOne]; ok && val != "" {
						singularEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: val,
							},
						}
					}
				}

				xcstrings.Strings[singularKey] = singularEntry
			}

			// Plural form
			if otherTrans, ok := trans.Translations["en"][models.QuantityOther]; ok && otherTrans != "" {
				pluralKey := trans.Key + ".plural"
				pluralEntry := XCStringsEntry{
					Comment:       trans.Comment + " (plural)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for plural
				for lang, quantities := range trans.Translations {
					if val, ok := quantities[models.QuantityOther]; ok && val != "" {
						pluralEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: val,
							},
						}
					}
				}

				xcstrings.Strings[pluralKey] = pluralEntry
			}

			continue
		}

		// Regular singular strings
		key = trans.Key

		// Skip if no English translation
		enTrans, ok := trans.Translations["en"]
		if !ok || len(enTrans) == 0 {
			continue
		}

		// Create entry
		entry := XCStringsEntry{
			Comment:       trans.Comment,
			Localizations: make(map[string]XCStringsLocalization),
		}

		// Add localizations
		for lang, quantities := range trans.Translations {
			if val, ok := quantities[models.QuantityOne]; ok && val != "" {
				entry.Localizations[lang] = XCStringsLocalization{
					StringUnit: XCStringsUnit{
						State: "translated",
						Value: val,
					},
				}
			}
		}

		xcstrings.Strings[key] = entry
	}

	// Write to file
	data, err := json.MarshalIndent(xcstrings, "", "  ")
	if err != nil {
		return fmt.Errorf("error generating iOS format: %v", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("error writing iOS format: %v", err)
	}

	fmt.Printf("Successfully exported to iOS format: %s\n", outputFile)
	return nil
}
