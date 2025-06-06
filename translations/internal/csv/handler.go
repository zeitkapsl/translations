package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

// DefaultCSVFile is the default filename for CSV exports
const DefaultCSVFile = "translations.csv"

// SaveToCSV saves a translation set to a CSV file
func SaveToCSV(ts *models.TranslationSet, filename string) error {
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

	// Write header
	header := []string{"Key", "Comment", "Type", "Source"}
	for _, lang := range ts.Languages {
		header = append(header, lang+"_one", lang+"_other")
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write translations
	for _, trans := range ts.Translations {
		record := []string{
			trans.Key,
			trans.Comment,
			string(trans.Type),
			string(trans.Source),
		}

		// Add translation values for each language
		for _, lang := range ts.Languages {
			oneVal := ""
			otherVal := ""

			if langTrans, ok := trans.Translations[lang]; ok {
				if val, ok := langTrans[models.QuantityOne]; ok {
					oneVal = val
				}
				if val, ok := langTrans[models.QuantityOther]; ok {
					otherVal = val
				}
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
func LoadFromCSV(filename string) (*models.TranslationSet, error) {
	if filename == "" {
		filename = DefaultCSVFile
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Parse header
	header := records[0]
	if len(header) < 4 {
		return nil, fmt.Errorf("invalid CSV header: expected at least 4 columns")
	}

	// Extract languages from header
	var languages []string
	langMap := make(map[string]int) // language -> starting column index

	for i := 4; i < len(header); i += 2 {
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

	// Create translation set
	ts := &models.TranslationSet{
		Languages:    languages,
		Translations: []models.Translation{},
	}

	// Process each record (skip header)
	for rowIdx, record := range records[1:] {
		if len(record) < 4 {
			fmt.Printf("Warning: skipping row %d with insufficient columns\n", rowIdx+2)
			continue
		}

		key := record[0]
		comment := record[1]
		typeStr := record[2]
		sourceStr := record[3]

		// Parse type
		var transType models.TranslationType
		switch typeStr {
		case string(models.TypeSingular):
			transType = models.TypeSingular
		case string(models.TypePlural):
			transType = models.TypePlural
		default:
			transType = models.TypeSingular // default
		}

		// Parse source
		var source models.TranslationSource
		switch sourceStr {
		case string(models.SourceIOS):
			source = models.SourceIOS
		case string(models.SourceAndroid):
			source = models.SourceAndroid
		case string(models.SourceWeb):
			source = models.SourceWeb
		case string(models.SourceManual):
			source = models.SourceManual
		default:
			source = models.SourceManual // default
		}

		// Create translation
		trans := models.Translation{
			Key:          key,
			Comment:      comment,
			Type:         transType,
			Source:       source,
			Translations: make(map[string]map[models.TranslationQuantity]string),
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
				trans.Translations[lang] = make(map[models.TranslationQuantity]string)
				if oneVal != "" {
					trans.Translations[lang][models.QuantityOne] = oneVal
				}
				if otherVal != "" {
					trans.Translations[lang][models.QuantityOther] = otherVal
				}
			}
		}

		ts.Translations = append(ts.Translations, trans)
	}

	fmt.Printf("Loaded %d translations from %s\n", len(ts.Translations), filename)
	return ts, nil
}
