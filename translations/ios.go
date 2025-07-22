package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
func ImportFromXCStrings(tm *Translations, baseDirectory string) error {
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
						tm.SetPluralSingular("ios", baseKey, lang, localization.StringUnit.Value, comment)
					} else if strings.HasSuffix(key, ".plural") {
						tm.SetPluralOther("ios", baseKey, lang, localization.StringUnit.Value, comment)
					} else {
						tm.SetTranslation("ios", key, lang, localization.StringUnit.Value, comment)
					}
				} else {
					tm.SetTranslation("ios", key, lang, localization.StringUnit.Value, comment)
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
func ExportToXCStrings(tm *Translations, baseDirectory string) error {
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

	// Get all iOS translations
	iosTranslations := tm.GetTranslationsForApp("ios")

	// Track processed plural keys to avoid duplicates
	processedPlurals := make(map[string]bool)

	// Process all translations
	for _, trans := range iosTranslations {
		// Handle plural translations
		if trans.IsPlural() {
			baseKey := trans.GetSingularKey()
			
			// Skip if we've already processed this plural set
			if processedPlurals[baseKey] {
				continue
			}

			// Get both singular and plural forms
			pluralValues := tm.GetPlural("ios", baseKey)
			if pluralValues == nil {
				continue
			}

			// Create singular entry
			if hasNonEmptyValues(pluralValues.One) {
				singularKey := baseKey + ".singular"
				singularEntry := XCStringsEntry{
					Comment:       trans.Comment + " (singular)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for singular
				for lang, value := range pluralValues.One {
					if value != "" {
						singularEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: value,
							},
						}
					}
				}

				if len(singularEntry.Localizations) > 0 {
					xcstrings.Strings[singularKey] = singularEntry
				}
			}

			// Create plural entry
			if hasNonEmptyValues(pluralValues.Other) {
				pluralKey := baseKey + ".plural"
				pluralEntry := XCStringsEntry{
					Comment:       trans.Comment + " (plural)",
					Localizations: make(map[string]XCStringsLocalization),
				}

				// Add localizations for plural
				for lang, value := range pluralValues.Other {
					if value != "" {
						pluralEntry.Localizations[lang] = XCStringsLocalization{
							StringUnit: XCStringsUnit{
								State: "translated",
								Value: value,
							},
						}
					}
				}

				if len(pluralEntry.Localizations) > 0 {
					xcstrings.Strings[pluralKey] = pluralEntry
				}
			}

			processedPlurals[baseKey] = true
			continue
		}

		// Handle regular singular strings
		key := trans.Key

		if len(trans.Values) == 0 {
			continue
		}

		// Create entry
		entry := XCStringsEntry{
			Comment:       trans.Comment,
			Localizations: make(map[string]XCStringsLocalization),
		}

		// Add localizations
		for lang, value := range trans.Values {
			if value != "" {
				entry.Localizations[lang] = XCStringsLocalization{
					StringUnit: XCStringsUnit{
						State: "translated",
						Value: value,
					},
				}
			}
		}

		if len(entry.Localizations) > 0 {
			xcstrings.Strings[key] = entry
		}
	}

	// Write to file
	data, err := json.MarshalIndent(xcstrings, "", "  ")
	if err != nil {
		return fmt.Errorf("error generating iOS format: %v", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("error writing iOS format: %v", err)
	}

	fmt.Printf("Successfully exported %d entries to iOS format: %s\n", len(xcstrings.Strings), outputFile)
	return nil
}

func isPluralKey(key string) bool {
	return strings.HasSuffix(key, ".plural") || strings.HasSuffix(key, ".singular")
}

func normalizeKey(key string) string {
	if strings.HasSuffix(key, ".singular") {
		return key[:len(key)-9] // Remove .singular
	}
	if strings.HasSuffix(key, ".plural") {
		return key[:len(key)-7] // Remove .plural
	}
	return key
}

func hasNonEmptyValues(values map[string]string) bool {
	if values == nil {
		return false
	}
	for _, value := range values {
		if value != "" {
			return true
		}
	}
	return false
}