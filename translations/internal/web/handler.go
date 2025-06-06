package web

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// ImportFromJSON imports translations from Web/Backend JSON files or JavaScript translation files
func ImportFromJSON(ts *models.TranslationSet, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Look for Web project structure: web/src/lib/i18n/
	webI18nPath := filepath.Join(baseDirectory, "web", "src", "lib", "i18n")
	
	// Also try alternative paths if the standard one doesn't exist
	alternativePaths := []string{
		filepath.Join(baseDirectory, "web", "i18n"),
		filepath.Join(baseDirectory, "web", "translations"),
		filepath.Join(baseDirectory, "web", "locales"),
		filepath.Join(baseDirectory, "web"),
	}
	
	// Check if the primary path exists, if not try alternatives
	if _, err := os.Stat(webI18nPath); os.IsNotExist(err) {
		found := false
		for _, altPath := range alternativePaths {
			if _, err := os.Stat(altPath); err == nil {
				webI18nPath = altPath
				found = true
				fmt.Printf("Using alternative web path: %s\n", webI18nPath)
				break
			}
		}
		if !found {
			return fmt.Errorf("Web i18n directory not found at %s or any alternative paths", webI18nPath)
		}
	}

	// First try to find and parse a translations.js file
	jsFile := filepath.Join(webI18nPath, "translations.js")
	if _, err := os.Stat(jsFile); err == nil {
		return importFromJavaScriptFile(ts, jsFile)
	}

	// Fallback to JSON files like en.json, de.json, etc.
	files, err := filepath.Glob(filepath.Join(webI18nPath, "*.json"))
	if err != nil {
		return fmt.Errorf("error finding JSON files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no JSON translation files or translations.js found in %s", webI18nPath)
	}

	importCount := 0

	for _, file := range files {
		// Extract language code from filename
		filename := filepath.Base(file)
		lang := strings.TrimSuffix(filename, ".json")

		// Skip files that don't look like language files
		if len(lang) > 5 || !isValidLanguageCode(lang) {
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
				ts.AddOrUpdatePluralTranslationWithSource(baseKey, "", lang, models.QuantityOne, value, models.SourceWeb)
				importCount++
			} else if strings.HasSuffix(key, ".plural") {
				baseKey := strings.TrimSuffix(key, ".plural")
				ts.AddOrUpdatePluralTranslationWithSource(baseKey, "", lang, models.QuantityOther, value, models.SourceWeb)
				importCount++
			} else {
				// Regular string
				ts.AddOrUpdateTranslationWithSource(key, "", lang, value, models.SourceWeb)
				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from JSON files\n", importCount)
	return nil
}

// importFromJavaScriptFile parses a translations.js file with the specific format used
func importFromJavaScriptFile(ts *models.TranslationSet, filePath string) error {
	fmt.Printf("Importing from JavaScript file: %s...\n", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	content := string(data)
	
	// Parse the JavaScript file to extract translations
	translations, err := parseJavaScriptTranslations(content)
	if err != nil {
		return fmt.Errorf("error parsing JavaScript translations: %v", err)
	}

	importCount := 0

	// Process each language
	for lang, langTranslations := range translations {
		// Skip invalid language codes
		if !isValidLanguageCode(lang) {
			fmt.Printf("Skipping invalid language code: %s\n", lang)
			continue
		}

		// Add language if not already present
		ts.AddLanguage(lang)

		// Process translations for this language
		for key, value := range langTranslations {
			// Check if it's a plural form
			if strings.HasSuffix(key, ".singular") {
				baseKey := strings.TrimSuffix(key, ".singular")
				ts.AddOrUpdatePluralTranslationWithSource(baseKey, "", lang, models.QuantityOne, value, models.SourceWeb)
				importCount++
			} else if strings.HasSuffix(key, ".plural") {
				baseKey := strings.TrimSuffix(key, ".plural")
				ts.AddOrUpdatePluralTranslationWithSource(baseKey, "", lang, models.QuantityOther, value, models.SourceWeb)
				importCount++
			} else {
				// Regular string
				ts.AddOrUpdateTranslationWithSource(key, "", lang, value, models.SourceWeb)
				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from JavaScript file\n", importCount)
	return nil
}

// parseJavaScriptTranslations parses the specific JavaScript translation format
func parseJavaScriptTranslations(content string) (map[string]map[string]string, error) {
	// Remove import statements and export default
	lines := strings.Split(content, "\n")
	var objectStart int = -1
	var braceCount int = 0
	var inObject bool = false

	// Find where the object starts
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export default") {
			objectStart = i
			break
		}
	}

	if objectStart == -1 {
		return nil, fmt.Errorf("could not find export default statement")
	}

	// Extract the object content
	var objectLines []string
	for i := objectStart; i < len(lines); i++ {
		line := lines[i]
		
		// Count braces to know when the object ends
		for _, char := range line {
			if char == '{' {
				braceCount++
				inObject = true
			} else if char == '}' {
				braceCount--
			}
		}

		if inObject {
			objectLines = append(objectLines, line)
		}

		// If we've closed all braces, we're done
		if inObject && braceCount == 0 {
			break
		}
	}

	// Join all object lines
	objectContent := strings.Join(objectLines, "\n")
	
	// Remove export default and clean up
	objectContent = strings.Replace(objectContent, "export default", "", 1)
	objectContent = strings.TrimSpace(objectContent)
	
	// Remove trailing semicolon if present
	if strings.HasSuffix(objectContent, ";") {
		objectContent = objectContent[:len(objectContent)-1]
	}

	// Parse the object manually
	return parseObjectContent(objectContent)
}

// parseObjectContent manually parses the JavaScript object
func parseObjectContent(content string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	
	// Remove outer braces
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") {
		content = content[1:]
	}
	if strings.HasSuffix(content, "}") {
		content = content[:len(content)-1]
	}

	// Split by languages (look for top-level keys)
	langSections := extractLanguageSections(content)
	
	for lang, section := range langSections {
		translations, err := parseTranslationSection(section)
		if err != nil {
			fmt.Printf("Warning: Error parsing section for %s: %v\n", lang, err)
			continue
		}
		result[lang] = translations
	}

	return result, nil
}

// extractLanguageSections splits the content into language sections
func extractLanguageSections(content string) map[string]string {
	sections := make(map[string]string)
	
	// Find language keys like "en-US": { ... }
	lines := strings.Split(content, "\n")
	var currentLang string
	var currentSection []string
	var braceCount int
	var inSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line starts a new language section
		if strings.Contains(line, ":") && (strings.Contains(line, "\"en-") || strings.Contains(line, "\"de-") || strings.Contains(line, "'en-") || strings.Contains(line, "'de-")) {
			// Save previous section if exists
			if currentLang != "" && len(currentSection) > 0 {
				sections[currentLang] = strings.Join(currentSection, "\n")
			}

			// Extract language code
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				langPart := strings.TrimSpace(parts[0])
				langPart = strings.Trim(langPart, "\"'")
				currentLang = langPart
				currentSection = []string{}
				braceCount = 0
				inSection = true

				// Add the rest of the line (after the colon)
				rest := strings.Join(parts[1:], ":")
				currentSection = append(currentSection, rest)
			}
		} else if inSection {
			currentSection = append(currentSection, line)
		}

		// Count braces
		for _, char := range line {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
			}
		}

		// If braces are balanced and we're at the end of a section
		if inSection && braceCount == 0 && strings.Contains(line, "}") {
			if currentLang != "" {
				sections[currentLang] = strings.Join(currentSection, "\n")
				currentLang = ""
				currentSection = []string{}
				inSection = false
			}
		}
	}

	// Don't forget the last section
	if currentLang != "" && len(currentSection) > 0 {
		sections[currentLang] = strings.Join(currentSection, "\n")
	}

	return sections
}

// parseTranslationSection parses individual translation key-value pairs
func parseTranslationSection(section string) (map[string]string, error) {
	translations := make(map[string]string)
	
	// Remove outer braces
	section = strings.TrimSpace(section)
	if strings.HasPrefix(section, "{") {
		section = section[1:]
	}
	if strings.HasSuffix(section, "}") {
		section = section[:len(section)-1]
	}

	// Split by lines and parse each key-value pair
	lines := strings.Split(section, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "," {
			continue
		}

		// Remove trailing comma
		if strings.HasSuffix(line, ",") {
			line = line[:len(line)-1]
		}

		// Parse key: value
		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])

		// Clean up key (remove quotes)
		key = strings.Trim(key, "\"'")
		
		// Clean up value (remove quotes and handle template literals)
		if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
			// Template literal - convert to simple string
			value = strings.Trim(value, "`")
			// Replace ${new Date().getFullYear()} with current year
			value = strings.ReplaceAll(value, "${new Date().getFullYear()}", "2024")
		} else {
			// Regular string - remove quotes
			value = strings.Trim(value, "\"'")
		}

		if key != "" && value != "" {
			translations[key] = value
		}
	}

	return translations, nil
}

// ExportToJSON exports translations to Web/Backend JSON files
func ExportToJSON(ts *models.TranslationSet, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Export to Web project structure: web/src/lib/i18n/
	webI18nPath := filepath.Join(baseDirectory, "web", "src", "lib", "i18n")

	// Get only Web-sourced translations
	webTranslations := ts.GetTranslationsBySource(models.SourceWeb)
	if len(webTranslations) == 0 {
		fmt.Println("No Web translations to export")
		return nil
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(webI18nPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", webI18nPath, err)
	}

	// For each language, create a .json file
	for _, lang := range ts.Languages {
		// Check if there are any Web translations for this language
		hasTranslations := false
		for _, trans := range webTranslations {
			if _, ok := trans.Translations[lang]; ok {
				hasTranslations = true
				break
			}
		}

		if !hasTranslations {
			continue
		}

		// Create JSON structure
		jsonTranslations := make(map[string]string)

		// Add translations (only Web-sourced ones)
		for _, trans := range webTranslations {
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
		filePath := filepath.Join(webI18nPath, lang+".json")
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return fmt.Errorf("error writing %s: %v", filePath, err)
		}

		fmt.Printf("Exported JSON translations to %s\n", filePath)
	}

	return nil
}

// isValidLanguageCode checks if a string is a valid language code
func isValidLanguageCode(code string) bool {
	// Simple validation for language code formats
	if len(code) == 2 {
		// Basic language code (e.g., en, de, fr)
		return true
	}

	if len(code) == 5 && code[2] == '-' {
		// Language with region (e.g., en-US, de-AT)
		return true
	}

	return false
}