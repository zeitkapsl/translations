package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// DefaultCSVFile is the default filename for CSV exports
const DefaultCSVFile = "translations.csv"

// SaveToCSV saves a translation set to a CSV file
func SaveToCSV(tm *models.TranslationManager, filename string) error {
	if filename == "" {
		filename = DefaultCSVFile
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Sort languages
	sort.Strings(tm.Languages)

	// Write header
	header := []string{"key", "comment", "type"}
	for _, lang := range tm.Languages {
		header = append(header, lang+"_one", lang+"_other")
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write translations
	for _, trans := range tm.Translations {
		record := []string{trans.Key, trans.Comment, trans.Type}

		// Add translation values for each language
		for _, lang := range tm.Languages {
			oneVal := ""
			otherVal := ""

			if langTrans, ok := trans.Values[lang]; ok {
				oneVal = langTrans.One
				otherVal = langTrans.Other
			}

			record = append(record, oneVal, otherVal)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	fmt.Printf("Saved translations to %s\n", filename)
	return nil
}

// LoadFromCSV loads a translation set from a CSV file
func LoadFromCSV(tm *models.TranslationManager, filename string) error {
	if filename == "" {
		filename = DefaultCSVFile
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Parse header
	header := records[0]
	if len(header) < 3 || header[0] != "key" || header[1] != "comment" || header[2] != "type" {
		return fmt.Errorf("invalid CSV header: expected at least 3 columns with key, comment, type")
	}

	// Extract languages from header
	var languages []string
	langMap := make(map[string]int) // language -> starting column index

	for i := 3; i < len(header); i += 2 {
		if i+1 >= len(header) {
			break
		}

		// Extract language from column names like "en_one", "en_other"
		langOne := header[i]
		langOther := header[i+1]

		if !strings.HasSuffix(langOne, "_one") || !strings.HasSuffix(langOther, "_other") {
			continue
		}

		lang := strings.TrimSuffix(langOne, "_one")
		expectedOther := lang + "_other"

		if langOther != expectedOther {
			continue
		}

		languages = append(languages, lang)
		langMap[lang] = i
	}

	tm.Languages = languages
	sort.Strings(tm.Languages)

	// Parse rows
	tm.Translations = make([]models.TranslationRow, 0, len(records)-1)
	for rowIdx, record := range records[1:] {
		if len(record) < 3 {
			fmt.Printf("Warning: skipping row %d with insufficient columns\n", rowIdx+2)
			continue
		}

		key := record[0]
		comment := record[1]
		typeStr := record[2]

		// Create translation
		trans := models.TranslationRow{
			Key:     key,
			Comment: comment,
			Type:    typeStr,
			Values:  make(map[string]models.TranslationValues),
		}

		// Parse language translations
		for _, lang := range languages {
			colIdx, ok := langMap[lang]
			if !ok || colIdx+1 >= len(record) {
				continue
			}

			oneVal := record[colIdx]
			otherVal := record[colIdx+1]

			if oneVal != "" || otherVal != "" {
				trans.Values[lang] = models.TranslationValues{
					One:   oneVal,
					Other: otherVal,
				}
			}
		}

		tm.Translations = append(tm.Translations, trans)
	}

	fmt.Printf("Loaded %d translations from %s\n", len(tm.Translations), filename)
	return nil
}
