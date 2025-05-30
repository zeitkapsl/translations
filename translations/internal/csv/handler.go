package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/zeitkapsl/translations/internal/models"
)

const (
	DefaultCSVFile   = "translations.csv"
	KeyColumn        = 0
	CommentColumn    = 1
	TypeColumn       = 2
	QuantityColumn   = 3
	FirstLangColumn  = 4
)

// SaveToCSV saves the TranslationSet to a CSV file
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

	// Sort languages for consistent ordering
	sortedLanguages := make([]string, len(ts.Languages))
	copy(sortedLanguages, ts.Languages)
	
	sort.Slice(sortedLanguages, func(i, j int) bool {
		// Base languages come before regions
		iHasRegion := strings.Contains(sortedLanguages[i], "-")
		jHasRegion := strings.Contains(sortedLanguages[j], "-")
		
		if iHasRegion != jHasRegion {
			return !iHasRegion // Base languages first
		}
		
		// Group regions with their base language
		baseLangI := strings.Split(sortedLanguages[i], "-")[0]
		baseLangJ := strings.Split(sortedLanguages[j], "-")[0]
		
		if baseLangI != baseLangJ {
			return baseLangI < baseLangJ
		}
		
		return sortedLanguages[i] < sortedLanguages[j]
	})

	// Write header row
	header := []string{"key", "comment", "type", "quantity"}
	header = append(header, sortedLanguages...)
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write translations
	for _, trans := range ts.Translations {
		// For singular translations, write a single row
		if trans.Type == models.TypeSingular {
			row := []string{trans.Key, trans.Comment, string(trans.Type), string(models.QuantityOne)}
			for _, lang := range sortedLanguages {
				value := ""
				if translations, ok := trans.Translations[lang]; ok {
					if val, ok := translations[models.QuantityOne]; ok {
						value = val
					}
				}
				row = append(row, value)
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("error writing CSV row: %v", err)
			}
		} else {
			// For plural translations, write a row for each quantity
			quantities := []models.TranslationQuantity{models.QuantityOne, models.QuantityOther}
			for _, quantity := range quantities {
				row := []string{trans.Key, trans.Comment, string(trans.Type), string(quantity)}
				for _, lang := range sortedLanguages {
					value := ""
					if translations, ok := trans.Translations[lang]; ok {
						if val, ok := translations[quantity]; ok {
							value = val
						}
					}
					row = append(row, value)
				}
				if err := writer.Write(row); err != nil {
					return fmt.Errorf("error writing CSV row: %v", err)
				}
			}
		}
	}

	return nil
}

// LoadFromCSV loads translations from a CSV file
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
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Parse header
	header := rows[0]
	if len(header) < 5 || header[KeyColumn] != "key" || header[CommentColumn] != "comment" || 
	   header[TypeColumn] != "type" || header[QuantityColumn] != "quantity" {
		return nil, fmt.Errorf("invalid CSV format, expected columns: key, comment, type, quantity, languages")
	}

	// Extract languages from header
	languages := header[FirstLangColumn:]

	// Create TranslationSet
	ts := &models.TranslationSet{
		Translations: []models.Translation{},
		Languages:    languages,
	}

	// Track processed keys to consolidate plural forms
	processedKeys := make(map[string]int) // key -> index in ts.Translations

	// Parse rows
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < len(header) {
			fmt.Printf("Warning: Row %d has fewer columns than header, skipping\n", i+1)
			continue
		}

		key := row[KeyColumn]
		comment := row[CommentColumn]
		typeName := models.TranslationType(row[TypeColumn])
		quantity := models.TranslationQuantity(row[QuantityColumn])

		// Check if we already have this key
		if idx, ok := processedKeys[key]; ok {
			// Update existing translation
			for j, lang := range languages {
				if j+FirstLangColumn < len(row) && row[j+FirstLangColumn] != "" {
					if ts.Translations[idx].Translations[lang] == nil {
						ts.Translations[idx].Translations[lang] = make(map[models.TranslationQuantity]string)
					}
					ts.Translations[idx].Translations[lang][quantity] = row[j+FirstLangColumn]
				}
			}
			// Update type if needed
			if typeName == models.TypePlural {
				ts.Translations[idx].Type = models.TypePlural
			}
		} else {
			// Create new translation
			trans := models.Translation{
				Key:     key,
				Comment: comment,
				Type:    typeName,
				Translations: map[string]map[models.TranslationQuantity]string{},
			}

			// Parse translations
			for j, lang := range languages {
				if j+FirstLangColumn < len(row) && row[j+FirstLangColumn] != "" {
					if trans.Translations[lang] == nil {
						trans.Translations[lang] = make(map[models.TranslationQuantity]string)
					}
					trans.Translations[lang][quantity] = row[j+FirstLangColumn]
				}
			}

			ts.Translations = append(ts.Translations, trans)
			processedKeys[key] = len(ts.Translations) - 1
		}
	}

	return ts, nil
}