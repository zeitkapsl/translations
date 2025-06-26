package web

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// ImportFromWeb imports translations from Web/Backend JavaScript files
func ImportFromWeb(tm *models.TranslationManager, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	// Look for Web project structure: web/src/lib/i18n/
	webI18nPath := filepath.Join(baseDirectory, "web", "src", "lib", "i18n")

	// Check if the primary path exists
	if _, err := os.Stat(webI18nPath); os.IsNotExist(err) {
		return fmt.Errorf("web i18n directory not found at %s", webI18nPath)
	}

	// Try to find and parse a translations.js file
	jsFile := filepath.Join(webI18nPath, "translations.js")
	if _, err := os.Stat(jsFile); err == nil {
		return importFromJavaScriptFile(tm, jsFile)
	}

	return fmt.Errorf("no translations.js file found in %s", webI18nPath)
}

// importFromJavaScriptFile parses a translations.js file with the specific format used
func importFromJavaScriptFile(tm *models.TranslationManager, filePath string) error {
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
		tm.EnsureLanguage(lang)

		// Process translations for this language
		for key, value := range langTranslations {
			// Check if it's a plural form
			if strings.HasSuffix(key, ".singular") {
				baseKey := strings.TrimSuffix(key, ".singular")
				tm.SetPluralSingular(baseKey, lang, value, "")
				importCount++
			} else if strings.HasSuffix(key, ".plural") {
				baseKey := strings.TrimSuffix(key, ".plural")
				tm.SetPluralOther(baseKey, lang, value, "")
				importCount++
			} else {
				// Regular string
				tm.SetTranslation(key, lang, value, "", "singular")
				importCount++
			}
		}
	}

	fmt.Printf("Imported %d translations from JavaScript file\n", importCount)
	return nil
}

// ExportToWeb exports translations to separate JSON files for each language
func ExportToWeb(tm *models.TranslationManager, baseDirectory string) error {
	if baseDirectory == "" {
		baseDirectory = "."
	}

	webPath := filepath.Join(baseDirectory, "web", "src", "lib", "i18n")

	if err := os.MkdirAll(webPath, 0755); err != nil {
		return err
	}

	// Export each language as a separate JSON file
	for _, lang := range tm.Languages {
		if err := writeLanguageJSONFile(tm, webPath, lang); err != nil {
			return fmt.Errorf("error writing JSON file for %s: %v", lang, err)
		}
	}

	fmt.Printf("Exported web translations to separate JSON files in %s\n", webPath)
	return nil
}

// writeLanguageJSONFile writes a single language to a JSON file
func writeLanguageJSONFile(tm *models.TranslationManager, webPath, lang string) error {
	translations := make(map[string]interface{})

	// Add translations for this language
	for _, row := range tm.Translations {
		values := tm.GetTranslationWithFallback(row, lang)

		if row.Type == "plural" {
			// For plurals, create separate .singular and .plural keys
			if values.One != "" {
				translations[row.Key+".singular"] = values.One
			}
			if values.Other != "" {
				translations[row.Key+".plural"] = values.Other
			}
		} else {
			// Regular singular translation
			if values.One != "" {
				translations[row.Key] = values.One
			}
		}
	}

	// Skip if no translations for this language
	if len(translations) == 0 {
		fmt.Printf("Skipping %s.json - no translations found\n", lang)
		return nil
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON for %s: %v", lang, err)
	}

	// Write to file
	filePath := filepath.Join(webPath, lang+".json")
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing %s: %v", filePath, err)
	}

	fmt.Printf("Exported %s.json with %d translations\n", lang, len(translations))
	return nil
}

// getTranslationWithFallback gets translation with regional fallback support
func getTranslationWithFallback(tm *models.TranslationManager, row models.TranslationRow, lang string) models.TranslationValues {
	return tm.GetTranslationWithFallback(row, lang)
}

// parseJavaScriptTranslations parses the complex JavaScript translation format
func parseJavaScriptTranslations(content string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// Find the export default object
	exportStart := strings.Index(content, "export default {")
	if exportStart == -1 {
		return nil, fmt.Errorf("could not find export default")
	}

	// Extract the object content
	objectContent := content[exportStart+len("export default {"):]

	// Find matching closing brace
	braceCount := 1
	objectEnd := 0
	for i, char := range objectContent {
		if char == '{' {
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 {
				objectEnd = i
				break
			}
		}
	}

	if objectEnd == 0 {
		return nil, fmt.Errorf("could not find end of object")
	}

	objectContent = objectContent[:objectEnd]

	// Parse language blocks with improved handling
	languages := extractLanguageSections(objectContent)

	for lang, section := range languages {
		translations := parseLanguageBlock(section)
		if len(translations) > 0 {
			result[lang] = translations
		}
	}

	return result, nil
}

// extractLanguageSections splits the content into language sections with better parsing
func extractLanguageSections(content string) map[string]string {
	sections := make(map[string]string)

	// Split by lines and process
	lines := strings.Split(content, "\n")
	var currentLang string
	var currentLines []string
	var braceCount int
	var inLanguageBlock bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line starts a new language section like "en-US": { or "de-DE": {
		if matches := regexp.MustCompile(`^\s*"([^"]+)"\s*:\s*\{`).FindStringSubmatch(line); len(matches) > 1 {
			// Save previous section
			if currentLang != "" && len(currentLines) > 0 {
				sections[currentLang] = strings.Join(currentLines, "\n")
			}

			// Start new section
			currentLang = matches[1]
			currentLines = []string{}
			braceCount = 1
			inLanguageBlock = true
			continue
		}

		if inLanguageBlock {
			currentLines = append(currentLines, line)

			// Count braces to track nesting
			for _, char := range line {
				if char == '{' {
					braceCount++
				} else if char == '}' {
					braceCount--
				}
			}

			// End of language block
			if braceCount == 0 {
				if currentLang != "" {
					sections[currentLang] = strings.Join(currentLines, "\n")
				}
				currentLang = ""
				currentLines = []string{}
				inLanguageBlock = false
			}
		}
	}

	// Don't forget the last section
	if currentLang != "" && len(currentLines) > 0 {
		sections[currentLang] = strings.Join(currentLines, "\n")
	}

	return sections
}

// parseLanguageBlock parses individual translation key-value pairs with advanced handling
func parseLanguageBlock(block string) map[string]string {
	result := make(map[string]string)

	// Remove the trailing }
	block = strings.TrimSuffix(strings.TrimSpace(block), "},")
	block = strings.TrimSuffix(strings.TrimSpace(block), "}")

	// Split into lines and process each one
	lines := strings.Split(block, "\n")

	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		if line == "" || line == "," {
			i++
			continue
		}

		// Remove trailing comma
		line = strings.TrimSuffix(line, ",")

		// Parse different types of entries
		if entry := parseSimpleKeyValue(line); entry != nil {
			result[entry.key] = entry.value
		} else if entry := parseTemplateString(line); entry != nil {
			result[entry.key] = entry.value
		} else if entry := parseArrayValue(line, lines, &i); entry != nil {
			result[entry.key] = entry.value
		} else if entry := parseFunction(line, lines, &i); entry != nil {
			// Skip functions for now, but you could handle them if needed
			// result[entry.key] = entry.value
		}

		i++
	}

	return result
}

type keyValuePair struct {
	key   string
	value string
}

// parseSimpleKeyValue handles simple key: "value" pairs
func parseSimpleKeyValue(line string) *keyValuePair {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return nil
	}

	key := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])

	// Clean up key
	key = strings.Trim(key, "\"'")

	// Handle simple quoted strings
	if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		value = strings.Trim(value, "\"'")
		// Unescape common escape sequences
		value = strings.ReplaceAll(value, "\\\"", "\"")
		value = strings.ReplaceAll(value, "\\'", "'")
		value = strings.ReplaceAll(value, "\\n", "\n")
		return &keyValuePair{key: key, value: value}
	}

	return nil
}

// parseTemplateString handles template literals like `text ${variable} more text`
func parseTemplateString(line string) *keyValuePair {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return nil
	}

	key := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])

	// Clean up key
	key = strings.Trim(key, "\"'")

	// Handle template literals
	if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
		value = strings.Trim(value, "`")
		// Replace template expressions with placeholders
		value = regexp.MustCompile(`\$\{[^}]+\}`).ReplaceAllString(value, "%s")
		// Handle HTML tags
		value = strings.ReplaceAll(value, "<br\\>", "<br>")
		return &keyValuePair{key: key, value: value}
	}

	return nil
}

// parseArrayValue handles array values like ["jan", "feb", ...]
func parseArrayValue(line string, lines []string, currentIndex *int) *keyValuePair {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return nil
	}

	key := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])

	// Clean up key
	key = strings.Trim(key, "\"'")

	// Check if this starts an array
	if strings.HasPrefix(value, "[") {
		arrayContent := value

		// If array doesn't close on same line, collect more lines
		if !strings.HasSuffix(strings.TrimSpace(value), "]") && !strings.Contains(value, "],") {
			for j := *currentIndex + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				arrayContent += " " + nextLine
				*currentIndex = j
				if strings.HasSuffix(nextLine, "],") || strings.HasSuffix(nextLine, "]") {
					break
				}
			}
		}

		// Extract array elements
		arrayContent = strings.TrimSuffix(strings.TrimSuffix(arrayContent, "],"), "]")
		arrayContent = strings.TrimPrefix(arrayContent, "[")

		// Parse array elements
		elements := []string{}
		for _, element := range strings.Split(arrayContent, ",") {
			element = strings.TrimSpace(element)
			element = strings.Trim(element, "\"'")
			if element != "" {
				elements = append(elements, element)
			}
		}

		// Join array elements with commas for storage
		if len(elements) > 0 {
			return &keyValuePair{key: key, value: strings.Join(elements, ",")}
		}
	}

	return nil
}

// parseFunction handles function values like (param) => `result ${param}`
func parseFunction(line string, lines []string, currentIndex *int) *keyValuePair {
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return nil
	}

	key := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])

	// Clean up key
	key = strings.Trim(key, "\"'")

	// Check if this is a function
	if strings.Contains(value, "=>") {
		// For now, skip functions as they're complex to handle
		// You could extract the template part if needed

		// If function spans multiple lines, skip them
		if !strings.HasSuffix(value, ",") {
			for j := *currentIndex + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				*currentIndex = j
				if strings.HasSuffix(nextLine, ",") || strings.HasSuffix(nextLine, "}") {
					break
				}
			}
		}

		return &keyValuePair{key: key, value: "FUNCTION_PLACEHOLDER"}
	}

	return nil
}

// escapeJS escapes a string for JavaScript
func escapeJS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// isValidLanguageCode checks if a string is a valid language code
func isValidLanguageCode(code string) bool {
	if len(code) == 2 {
		return true
	}
	if len(code) == 5 && code[2] == '-' {
		return true
	}
	return false
}
