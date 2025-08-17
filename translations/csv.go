package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

// DefaultCSVFile is the default filename for CSV exports
const DefaultCSVFile = "translations.csv"

// SaveToCSV saves a translation set to a CSV file
func SaveToCSV(tm *Translations, filename string) error {
	if filename == "" {
		filename = DefaultCSVFile
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Sort languages
	sort.Strings(tm.Languages)

	// Write header
	header := []string{"app", "key", "comment"}
	for _, lang := range tm.Languages {
		header = append(header, lang)
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write translations
	for _, trans := range tm.Translations {
		record := []string{trans.App, trans.Key, trans.Comment}
		// Add translation values for each language
		for _, lang := range tm.Languages {
			record = append(record, trans.Values[lang])
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	fmt.Printf("Saved translations to %s\n", filename)
	return nil
}

// LoadFromCSV loads a translation set from a CSV file
func LoadFromCSV(tm *Translations, filename string) error {
	if filename == "" {
		filename = DefaultCSVFile
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Parse header
	header := records[0]
	if len(header) < 3 || header[0] != "app" || header[1] != "key" || header[2] != "comment" {
		return fmt.Errorf("invalid CSV header: expected at least 3 columns with key, comment, type")
	}

	// Extract languages from header
	var languages []string
	langMap := make(map[string]int) // language -> starting column index

	for i := 3; i < len(header); i++ {
		lang := header[i]
		languages = append(languages, lang)
		langMap[lang] = i
	}

	tm.Languages = languages
	sort.Strings(tm.Languages)

	// Parse rows
	tm.Translations = make([]TranslationRow, 0, len(records)-1)
	for rowIdx, record := range records[1:] {
		if len(record) < 3 {
			fmt.Printf("Warning: skipping row %d with insufficient columns\n", rowIdx+2)
			continue
		}

		app := record[0]
		key := record[1]
		comment := record[2]

		// Create translation
		trans := TranslationRow{
			App:     app,
			Key:     key,
			Comment: comment,
			Values:  make(map[string]string),
		}

		// Parse language translations
		for _, lang := range languages {
			colIdx, ok := langMap[lang]
			if !ok {
				continue
			}
			trans.Values[lang] = record[colIdx]
		}

		tm.Translations = append(tm.Translations, trans)
	}

	fmt.Printf("Loaded %d translations from %s\n", len(tm.Translations), filename)
	return nil
}
