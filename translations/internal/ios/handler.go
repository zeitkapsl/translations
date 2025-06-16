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
func ImportFromXCStrings(tm *models.TranslationManager, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Look for iOS project structure: ios/Zeitkapsl/Supporting Files/
	iOSStringsPath := filepath.Join(baseDirectory, "ios", "Zeitkapsl", "Supporting Files", "Localizable.xcstrings")

	// Check if the iOS strings file exists
	if _, err := os.Stat(iOSStringsPath); os.IsNotExist(err) {
		return fmt.Errorf("iOS Localizable.xcstrings file not found at %s", iOSStringsPath)
	}

	importCount := 0

	fmt.Printf("Importing from %s...\n", iOSStringsPath)

	data, err := os.ReadFile(iOSStringsPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", iOSStringsPath, err)
	}

	var xcstrings XCStringsFile
	if err := json.Unmarshal(data, &xcstrings); err != nil {
		return fmt.Errorf("error parsing %s: %v", iOSStringsPath, err)
	}

	// Add source language if not already present
	tm.EnsureLanguage(xcstrings.SourceLanguage)

	// Process strings
	for key, entry := range xcstrings.Strings {
		comment := entry.Comment

		// Process localizations
		for lang, localization := range entry.Localizations {
			// Only add translations that are in "translated" state and have a value
			if localization.StringUnit.State == "translated" && localization.StringUnit.Value != "" {
				// Check if it's a plural key
				if isPluralKey(key) {
					baseKey := normalizeKey(key)
					if strings.HasSuffix(key, ".singular") {
						tm.SetPluralSingular(baseKey, lang, localization.StringUnit.Value, comment)
					} else if strings.HasSuffix(key, ".plural") {
						tm.SetPluralOther(baseKey, lang, localization.StringUnit.Value, comment)
					} else {
						tm.SetTranslation(key, lang, localization.StringUnit.Value, comment, "singular")
					}
				} else {
					tm.SetTranslation(key, lang, localization.StringUnit.Value, comment, "singular")
				}

				// Add language if not already present
				tm.EnsureLanguage(lang)

				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from iOS .xcstrings file\n", importCount)
	return nil
}

// ExportToXCStrings exports translations to iOS .xcstrings format
func ExportToXCStrings(tm *models.TranslationManager, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Export to iOS project structure: ios/Zeitkapsl/Supporting Files/
	iOSStringsDir := filepath.Join(baseDirectory, "ios", "Zeitkapsl", "Supporting Files")
	outputFile := filepath.Join(iOSStringsDir, "Localizable.xcstrings")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(iOSStringsDir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", iOSStringsDir, err)
	}

	// Create xcstrings structure
	xcstrings := XCStringsFile{
		Version:        "1.0",
		SourceLanguage: "en", // Assuming English is the source language
		Strings:        make(map[string]XCStringsEntry),
	}

	// Process all translations
	for _, trans := range tm.Translations {
		if trans.Type == "plural" {
			// Singular form
			if enTrans, ok := trans.Values["en"]; ok && enTrans.One != "" {
				singularKey := trans.Key + ".singular"
				singularEntry := XCStringsEntry{
					Comment:       trans.Comment + " (singular)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for singular
				for lang, quantities := range trans.Values {
					if quantities.One != "" {
						singularEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: quantities.One,
							},
						}
					}
				}

				xcstrings.Strings[singularKey] = singularEntry
			}

			// Plural form
			if enTrans, ok := trans.Values["en"]; ok && enTrans.Other != "" {
				pluralKey := trans.Key + ".plural"
				pluralEntry := XCStringsEntry{
					Comment:       trans.Comment + " (plural)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for plural
				for lang, quantities := range trans.Values {
					if quantities.Other != "" {
						pluralEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: quantities.Other,
							},
						}
					}
				}

				xcstrings.Strings[pluralKey] = pluralEntry
			}

			continue
		}

		// Regular singular strings
		key := trans.Key

		// Skip if no English translation
		enTrans, ok := trans.Values["en"]
		if !ok || enTrans.One == "" {
			continue
		}

		// Create entry
		entry := XCStringsEntry{
			Comment:       trans.Comment,
			Localizations: make(map[string]XCStringsLocalization),
		}

		// Add localizations
		for lang, quantities := range trans.Values {
			if quantities.One != "" {
				entry.Localizations[lang] = XCStringsLocalization{
					StringUnit: XCStringsUnit{
						State: "translated",
						Value: quantities.One,
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

// isPluralKey determines if a key is for a plural translation
func isPluralKey(key string) bool {
	return strings.HasSuffix(key, ".plural") || strings.HasSuffix(key, ".singular")
}

// normalizeKey extracts the base key without plural suffix
func normalizeKey(key string) string {
	if strings.HasSuffix(key, ".singular") {
		return key[:len(key)-9] // Remove .singular
	}
	if strings.HasSuffix(key, ".plural") {
		return key[:len(key)-7] // Remove .plural
	}
	return key
}
