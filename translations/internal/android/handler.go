package android

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// Resources represents the root element of an Android strings.xml file
type Resources struct {
	XMLName xml.Name         `xml:"resources"`
	Strings []StringResource `xml:"string"`
	Plurals []PluralResource `xml:"plurals"`
}

// StringResource represents a string element in an Android strings.xml file
type StringResource struct {
	XMLName xml.Name `xml:"string"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

// PluralResource represents a plurals element in an Android strings.xml file
type PluralResource struct {
	XMLName xml.Name     `xml:"plurals"`
	Name    string       `xml:"name,attr"`
	Items   []PluralItem `xml:"item"`
}

// PluralItem represents an item element within a plurals element
type PluralItem struct {
	XMLName  xml.Name `xml:"item"`
	Quantity string   `xml:"quantity,attr"`
	Value    string   `xml:",chardata"`
}

// ImportFromAndroid imports translations from Android strings.xml files
func ImportFromAndroid(ts *models.TranslationSet, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Look for Android project structure: android/app/src/main/res/
	androidResPath := filepath.Join(baseDirectory, "android", "app", "src", "main", "res")
	
	// Check if the Android res directory exists
	if _, err := os.Stat(androidResPath); os.IsNotExist(err) {
		return fmt.Errorf("Android res directory not found at %s", androidResPath)
	}

	// Scan for values directories (values, values-en, values-de, etc.)
	entries, err := os.ReadDir(androidResPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", androidResPath, err)
	}

	valuesDir := []string{}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "values") {
			valuesDir = append(valuesDir, entry.Name())
		}
	}

	if len(valuesDir) == 0 {
		return fmt.Errorf("no Android values directories found in %s", androidResPath)
	}

	importCount := 0

	for _, dir := range valuesDir {
		// Extract language code from directory name
		lang := "en" // Default language
		if dir != "values" {
			lang = strings.TrimPrefix(dir, "values-")
			// Handle regional variants like values-de-rAT -> de-AT
			lang = strings.Replace(lang, "-r", "-", 1)
		}

		file := filepath.Join(androidResPath, dir, "strings.xml")
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Warning: Could not read %s: %v\n", file, err)
			continue
		}

		var resources Resources
		if err := xml.Unmarshal(data, &resources); err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		// Add language if not already present
		ts.AddLanguage(lang)

		// Process regular strings
		for _, str := range resources.Strings {
			ts.AddOrUpdateTranslationWithSource(str.Name, "", lang, str.Value, models.SourceAndroid)
			importCount++
		}

		// Process plurals
		for _, plural := range resources.Plurals {
			for _, item := range plural.Items {
				quantity := models.TranslationQuantity(item.Quantity)
				if quantity == "one" {
					ts.AddOrUpdatePluralTranslationWithSource(plural.Name, "", lang, models.QuantityOne, item.Value, models.SourceAndroid)
				} else {
					// Android has many quantity types (zero, one, two, few, many, other)
					// but we map them all to "other" for simplicity
					ts.AddOrUpdatePluralTranslationWithSource(plural.Name, "", lang, models.QuantityOther, item.Value, models.SourceAndroid)
				}
				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from Android strings.xml files\n", importCount)
	return nil
}

// ExportToAndroid exports translations to Android strings.xml files
func ExportToAndroid(ts *models.TranslationSet, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Export to Android project structure: android/app/src/main/res/
	androidResPath := filepath.Join(baseDirectory, "android", "app", "src", "main", "res")

	// Get only Android-sourced translations
	androidTranslations := ts.GetTranslationsBySource(models.SourceAndroid)
	if len(androidTranslations) == 0 {
		fmt.Println("No Android translations to export")
		return nil
	}

	// For each language, create a strings.xml file in the appropriate directory
	for _, lang := range ts.Languages {
		// Check if there are any Android translations for this language
		hasTranslations := false
		for _, trans := range androidTranslations {
			if _, ok := trans.Translations[lang]; ok {
				hasTranslations = true
				break
			}
		}

		if !hasTranslations {
			continue
		}

		// Determine the values directory name
		dirName := "values"
		if lang != "en" {
			// Handle regional variants like de-AT -> values-de-rAT
			if strings.Contains(lang, "-") {
				parts := strings.Split(lang, "-")
				dirName = fmt.Sprintf("values-%s-r%s", parts[0], parts[1])
			} else {
				dirName = "values-" + lang
			}
		}

		dirPath := filepath.Join(androidResPath, dirName)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dirPath, err)
		}

		// Create XML structure
		resources := Resources{
			Strings: []StringResource{},
			Plurals: []PluralResource{},
		}

		// Track processed plural keys to avoid duplicates
		processedPlurals := make(map[string]bool)

		// Add translations (only Android-sourced ones)
		for _, trans := range androidTranslations {
			// Skip keys without translations for this language
			if _, ok := trans.Translations[lang]; !ok {
				continue
			}

			if trans.Type == models.TypeSingular {
				// Add as regular string
				value := ""
				if translations, ok := trans.Translations[lang]; ok {
					if val, ok := translations[models.QuantityOne]; ok {
						value = val
					}
				}

				if value != "" {
					resources.Strings = append(resources.Strings, StringResource{
						Name:  trans.Key,
						Value: value,
					})
				}
			} else if trans.Type == models.TypePlural && !processedPlurals[trans.Key] {
				// Add as plural
				pluralResource := PluralResource{
					Name:  trans.Key,
					Items: []PluralItem{},
				}

				// Add quantity items
				langTranslations := trans.Translations[lang]
				if val, ok := langTranslations[models.QuantityOne]; ok && val != "" {
					pluralResource.Items = append(pluralResource.Items, PluralItem{
						Quantity: string(models.QuantityOne),
						Value:    val,
					})
				}

				if val, ok := langTranslations[models.QuantityOther]; ok && val != "" {
					pluralResource.Items = append(pluralResource.Items, PluralItem{
						Quantity: string(models.QuantityOther),
						Value:    val,
					})
				}

				if len(pluralResource.Items) > 0 {
					resources.Plurals = append(resources.Plurals, pluralResource)
					processedPlurals[trans.Key] = true
				}
			}
		}

		// Skip if no translations for this language
		if len(resources.Strings) == 0 && len(resources.Plurals) == 0 {
			continue
		}

		// Marshal to XML
		xmlData, err := xml.MarshalIndent(resources, "", "    ")
		if err != nil {
			return fmt.Errorf("error generating XML for %s: %v", lang, err)
		}

		// Add XML header
		xmlContent := []byte(xml.Header + string(xmlData))

		// Write to file
		filePath := filepath.Join(dirPath, "strings.xml")
		if err := os.WriteFile(filePath, xmlContent, 0644); err != nil {
			return fmt.Errorf("error writing %s: %v", filePath, err)
		}

		fmt.Printf("Exported Android translations to %s\n", filePath)
	}

	return nil
}