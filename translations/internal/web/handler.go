package web

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// ImportFromJSON imports translations from Web/Backend JSON files
func ImportFromJSON(ts *models.TranslationSet, directory string) error {
	if directory == "" {
		directory = "."
	}

	// Find JSON files like en.json, de.json, etc.
	files, err := filepath.Glob(filepath.Join(directory, "*.json"))
	if err != nil {
		return fmt.Errorf("error finding JSON files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no JSON translation files found in %s", directory)
	}

	importCount := 0

	for _, file := range files {
		// Extract language code from filename
		filename := filepath.Base(file)
		lang := strings.TrimSuffix(filename, ".json")

		// Skip files that don't look like language files
		if len(lang) > 3 || !isValidLanguageCode(lang) {
			fmt.Printf("Skipping file %s as it doesn't match a language code format\n", filename)
			continue
		}

		fmt.Printf("Importing from %s...\n", file)

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		var jsonTranslations map[string]string
		if err := json.Unmarshal(data, &jsonTranslations); err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		// Add language if not already present
		ts.AddLanguage(lang)

		// Process translations
		for key, value := range jsonTranslations {
			// Check if it's a plural form
			if strings.HasSuffix(key, ".singular") {
				baseKey := strings.TrimSuffix(key, ".singular")
				ts.AddOrUpdatePluralTranslation(baseKey, "", lang, models.QuantityOne, value)
				importCount++
			} else if strings.HasSuffix(key, ".plural") {
				baseKey := strings.TrimSuffix(key, ".plural")
				ts.AddOrUpdatePluralTranslation(baseKey, "", lang, models.QuantityOther, value)
				importCount++
			} else {
				// Regular string
				ts.AddOrUpdateTranslation(key, "", lang, value)
				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from JSON files\n", importCount)
	return nil
}

// ExportToJSON exports translations to Web/Backend JSON files
func ExportToJSON(ts *models.TranslationSet, outputDir string) error {
	if outputDir == "" {
		outputDir = "."
	}

	// For each language, create a .json file
	for _, lang := range ts.Languages {
		// Create JSON structure
		jsonTranslations := make(map[string]string)

		// Add translations
		for _, trans := range ts.Translations {
			if translations, ok := trans.Translations[lang]; ok {
				if trans.Type == models.TypeSingular {
					if val, ok := translations[models.QuantityOne]; ok && val != "" {
						jsonTranslations[trans.Key] = val
					}
				} else {
					// For plurals, we need to use .singular and .plural suffixes
					if val, ok := translations[models.QuantityOne]; ok && val != "" {
						jsonTranslations[trans.Key+".singular"] = val
					}
					if val, ok := translations[models.QuantityOther]; ok && val != "" {
						jsonTranslations[trans.Key+".plural"] = val
					}
				}
			}
		}

		// Skip if no translations for this language
		if len(jsonTranslations) == 0 {
			continue
		}

		// Marshal to JSON
		jsonData, err := json.MarshalIndent(jsonTranslations, "", "  ")
		if err != nil {
			return fmt.Errorf("error generating JSON for %s: %v", lang, err)
		}

		// Write to file
		filePath := filepath.Join(outputDir, lang+".json")
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return fmt.Errorf("error writing %s: %v", filePath, err)
		}

		fmt.Printf("Exported JSON translations to %s\n", filePath)
	}

	return nil
}

// isValidLanguageCode checks if a string is a valid language code
func isValidLanguageCode(code string) bool {
	// Simple validation - language codes are typically 2 characters
	return len(code) == 2
}
