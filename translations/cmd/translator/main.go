package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/zeitkapsl/translations/internal/android"
	"github.com/zeitkapsl/translations/internal/csv"
	"github.com/zeitkapsl/translations/internal/ios"
	"github.com/zeitkapsl/translations/internal/models"
	"github.com/zeitkapsl/translations/internal/translator"
	"github.com/zeitkapsl/translations/internal/web"
)

func main() {
	// Define commands
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	addLangCmd := flag.NewFlagSet("add-language", flag.ExitOnError)
	addRegionCmd := flag.NewFlagSet("add-region", flag.ExitOnError)
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	translateCmd := flag.NewFlagSet("auto-translate", flag.ExitOnError)

	// Define flags for commands
	importDirFlag := importCmd.String("dir", ".", "Base directory containing translation folders")
	langFlag := addLangCmd.String("lang", "", "Language code to add (e.g., es)")
	regionFlag := addRegionCmd.String("region", "", "Region code to add (e.g., de-AT)")
	platformFlag := exportCmd.String("platform", "", "Platform to export to (ios|android|json|csv|all)")
	exportDirFlag := exportCmd.String("dir", ".", "Base directory for export")
	translateServiceFlag := translateCmd.String("service", "", "Translation service to use (azure|deepl)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse the command
	switch os.Args[1] {
	case "import":
		importCmd.Parse(os.Args[2:])
		importTranslations(*importDirFlag)
	case "add-language":
		addLangCmd.Parse(os.Args[2:])
		if *langFlag == "" {
			fmt.Println("Please provide a language code with --lang")
			os.Exit(1)
		}
		addLanguage(*langFlag)
	case "add-region":
		addRegionCmd.Parse(os.Args[2:])
		if *regionFlag == "" {
			fmt.Println("Please provide a region code with --region")
			os.Exit(1)
		}
		addRegion(*regionFlag)
	case "export":
		exportCmd.Parse(os.Args[2:])
		if *platformFlag == "" {
			fmt.Println("Please provide a platform with --platform")
			os.Exit(1)
		}
		exportTranslations(*platformFlag, *exportDirFlag)
	case "stats":
		statsCmd.Parse(os.Args[2:])
		showStats()
	case "auto-translate":
		translateCmd.Parse(os.Args[2:])
		autoTranslate(*translateServiceFlag)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: translator <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  import [--dir=<directory>]       Import translations from all platforms")
	fmt.Println("  add-language --lang=<code>       Add a new language (e.g., es, fr)")
	fmt.Println("  add-region --region=<code>       Add a regional variant (e.g., de-AT)")
	fmt.Println("  export --platform=<platform> [--dir=<directory>]")
	fmt.Println("                                   Export to platform format (ios|android|json|csv|all)")
	fmt.Println("  stats                            Show translation statistics")
	fmt.Println("  auto-translate                   Auto-translate missing strings")
	fmt.Println("  help                             Show this help")
	fmt.Println("")
	fmt.Println("Directory Structure Expected:")
	fmt.Println("  android/app/src/main/res/values*/strings.xml")
	fmt.Println("  ios/Zeitkapsl/Supporting Files/Localizable.xcstrings")
	fmt.Println("  web/src/lib/i18n/*.json")
}

// importTranslations imports translations from all platform-specific formats
func importTranslations(directory string) {
	// Create a new translation set or load existing one if available
	var ts *models.TranslationSet

	// Check if translations.csv already exists
	if _, err := os.Stat(csv.DefaultCSVFile); err == nil {
		fmt.Println("Found existing translations.csv, attempting to load...")

		// Try to load with new format first
		ts, err = csv.LoadFromCSV("")
		if err != nil {
			fmt.Println("Attempting to load CSV ...")
		}
	} else {
		ts = models.NewTranslationSet()
	}

	fmt.Println("Importing translations...")
	importSuccessCount := 0

	// Import from iOS .xcstrings
	if err := ios.ImportFromXCStrings(ts, directory); err != nil {
		fmt.Printf("Warning: iOS import failed: %v\n", err)
	} else {
		importSuccessCount++
	}

	// Import from Android strings.xml
	if err := android.ImportFromAndroid(ts, directory); err != nil {
		fmt.Printf("Warning: Android import failed: %v\n", err)
	} else {
		importSuccessCount++
	}

	// Import from Web/Backend JSON files
	if err := web.ImportFromJSON(ts, directory); err != nil {
		fmt.Printf("Warning: Web import failed: %v\n", err)
	} else {
		importSuccessCount++
	}

	if importSuccessCount == 0 {
		fmt.Println("No translations were imported from any platform")
		return
	}

	// Save to CSV
	if err := csv.SaveToCSV(ts, ""); err != nil {
		fmt.Printf("Error saving translations: %v\n", err)
		return
	}

	fmt.Printf("Import completed successfully from %d platform(s)\n", importSuccessCount)
}

// addLanguage adds a new language to the translation set
func addLanguage(lang string) {
	if !isValidLanguageCode(lang) {
		fmt.Printf("Invalid language code: %s. Expected format: xx (e.g., en, de, fr)\n", lang)
		return
	}

	ts, err := csv.LoadFromCSV("")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No translations.csv file found. Run import first.")
		} else {
			fmt.Printf("Error loading translations: %v\n", err)
		}
		return
	}

	added := ts.AddLanguage(lang)
	if !added {
		fmt.Printf("Language %s already exists in translations\n", lang)
		return
	}

	// Save updated translations
	if err := csv.SaveToCSV(ts, ""); err != nil {
		fmt.Printf("Error saving translations: %v\n", err)
		return
	}

	fmt.Printf("Language %s added successfully\n", lang)
}

// addRegion adds a regional variant to the translation set
func addRegion(region string) {
	parts := strings.Split(region, "-")
	if len(parts) != 2 {
		fmt.Println("Invalid region format. Expected format: xx-YY (e.g., de-AT)")
		return
	}

	baseLang := parts[0]

	ts, err := csv.LoadFromCSV("")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No translations.csv file found. Run import first.")
		} else {
			fmt.Printf("Error loading translations: %v\n", err)
		}
		return
	}

	// Check if base language exists
	baseExists := false
	for _, lang := range ts.Languages {
		if lang == baseLang {
			baseExists = true
			break
		}
	}

	if !baseExists {
		fmt.Printf("Base language %s does not exist. Add it first with add-language\n", baseLang)
		return
	}

	added := ts.AddLanguage(region)
	if !added {
		fmt.Printf("Region %s already exists in translations\n", region)
		return
	}

	// Save updated translations
	if err := csv.SaveToCSV(ts, ""); err != nil {
		fmt.Printf("Error saving translations: %v\n", err)
		return
	}

	fmt.Printf("Region %s added successfully\n", region)
}

// exportTranslations exports translations to platform-specific formats
func exportTranslations(platform, directory string) {
	ts, err := csv.LoadFromCSV("")
	if err != nil {
		fmt.Printf("Error loading translations: %v\n", err)
		return
	}

	switch platform {
	case "ios":
		if err := ios.ExportToXCStrings(ts, directory); err != nil {
			fmt.Printf("Error exporting to iOS format: %v\n", err)
		}
	case "android":
		if err := android.ExportToAndroid(ts, directory); err != nil {
			fmt.Printf("Error exporting to Android format: %v\n", err)
		}
	case "json", "web":
		if err := web.ExportToJSON(ts, directory); err != nil {
			fmt.Printf("Error exporting to JSON format: %v\n", err)
		}
	case "csv":
		// CSV is already saved, just notify
		fmt.Println("CSV format is the source format, no export needed")
	case "all":
		// Export to all platforms
		exportCount := 0

		if err := ios.ExportToXCStrings(ts, directory); err != nil {
			fmt.Printf("Warning: iOS export failed: %v\n", err)
		} else {
			exportCount++
		}

		if err := android.ExportToAndroid(ts, directory); err != nil {
			fmt.Printf("Warning: Android export failed: %v\n", err)
		} else {
			exportCount++
		}

		if err := web.ExportToJSON(ts, directory); err != nil {
			fmt.Printf("Warning: Web export failed: %v\n", err)
		} else {
			exportCount++
		}

		if exportCount > 0 {
			fmt.Printf("Export completed successfully to %d platform(s)\n", exportCount)
		} else {
			fmt.Println("No exports were successful")
		}
	default:
		fmt.Printf("Unknown platform: %s. Expected ios, android, json, csv, or all\n", platform)
	}
}

// showStats shows detailed translation statistics with source information
func showStats() {
	ts, err := csv.LoadFromCSV("")
	if err != nil {
		fmt.Printf("Error loading translations: %v\n", err)
		return
	}

	fmt.Println("Translation Statistics:")
	fmt.Println("======================")

	// Count total strings and types by source
	totalStrings := 0
	singularCount := 0
	pluralCount := 0

	sourceStats := map[models.TranslationSource]int{
		models.SourceIOS:     0,
		models.SourceAndroid: 0,
		models.SourceWeb:     0,
		models.SourceManual:  0,
	}

	for _, trans := range ts.Translations {
		if _, ok := trans.Translations["en"]; ok {
			totalStrings++
			sourceStats[trans.Source]++
			if trans.Type == models.TypeSingular {
				singularCount++
			} else {
				pluralCount++
			}
		}
	}

	fmt.Printf("Total strings: %d (%d singular, %d plural)\n", totalStrings, singularCount, pluralCount)
	fmt.Println("\nStrings by source:")
	fmt.Printf("  iOS: %d\n", sourceStats[models.SourceIOS])
	fmt.Printf("  Android: %d\n", sourceStats[models.SourceAndroid])
	fmt.Printf("  Web: %d\n", sourceStats[models.SourceWeb])
	fmt.Printf("  Manual: %d\n", sourceStats[models.SourceManual])
	fmt.Println()

	// Show language stats
	fmt.Println("Language coverage:")
	fmt.Println("-----------------")

	// Organize languages by base language
	baseLanguages := map[string][]string{}
	for _, lang := range ts.Languages {
		if strings.Contains(lang, "-") {
			// Regional variant
			baseLang := strings.Split(lang, "-")[0]
			baseLanguages[baseLang] = append(baseLanguages[baseLang], lang)
		} else {
			// Base language
			if _, ok := baseLanguages[lang]; !ok {
				baseLanguages[lang] = []string{}
			}
		}
	}

	// Sort base languages for consistent output
	sortedBaseLanguages := make([]string, 0, len(baseLanguages))
	for lang := range baseLanguages {
		sortedBaseLanguages = append(sortedBaseLanguages, lang)
	}
	sort.Strings(sortedBaseLanguages)

	// Calculate stats for each language
	for _, baseLang := range sortedBaseLanguages {
		if baseLang == "en" {
			continue // Skip English as it's the source language
		}

		// Base language stats
		singularTranslated := 0.0
		pluralTranslated := 0.0

		for _, trans := range ts.Translations {
			if _, ok := trans.Translations["en"]; !ok {
				continue // Skip if no English source
			}

			if langTrans, ok := trans.Translations[baseLang]; ok {
				if trans.Type == models.TypeSingular {
					if _, ok := langTrans[models.QuantityOne]; ok {
						singularTranslated++
					}
				} else {
					// For plurals, check if both forms are translated
					oneOk := false
					otherOk := false

					if _, ok := langTrans[models.QuantityOne]; ok {
						oneOk = true
					}
					if _, ok := langTrans[models.QuantityOther]; ok {
						otherOk = true
					}

					if oneOk && otherOk {
						pluralTranslated++
					} else if oneOk || otherOk {
						// Partial translation (only one form)
						pluralTranslated += 0.5
					}
				}
			}
		}

		totalTranslated := singularTranslated + pluralTranslated
		percentComplete := 0.0
		if totalStrings > 0 {
			percentComplete = float64(totalTranslated) / float64(totalStrings) * 100
		}

		missing := totalStrings - int(totalTranslated)

		fmt.Printf("%s: %.1f/%d (%.1f%%) - %d missing\n",
			baseLang, totalTranslated, totalStrings, percentComplete, missing)

		fmt.Printf("  Singular: %.0f/%d (%.1f%%)\n",
			singularTranslated, singularCount,
			float64(singularTranslated)/float64(singularCount)*100)

		fmt.Printf("  Plural: %.1f/%d (%.1f%%)\n",
			pluralTranslated, pluralCount,
			float64(pluralTranslated)/float64(pluralCount)*100)

		// Regional variants stats
		regions := baseLanguages[baseLang]
		sort.Strings(regions)

		for _, region := range regions {
			singularTranslated := 0.0
			pluralTranslated := 0.0

			for _, trans := range ts.Translations {
				if _, ok := trans.Translations["en"]; !ok {
					continue
				}

				if langTrans, ok := trans.Translations[region]; ok {
					if trans.Type == models.TypeSingular {
						if _, ok := langTrans[models.QuantityOne]; ok {
							singularTranslated++
						}
					} else {
						// For plurals, check if both forms are translated
						oneOk := false
						otherOk := false

						if _, ok := langTrans[models.QuantityOne]; ok {
							oneOk = true
						}
						if _, ok := langTrans[models.QuantityOther]; ok {
							otherOk = true
						}

						if oneOk && otherOk {
							pluralTranslated++
						} else if oneOk || otherOk {
							// Partial translation
							pluralTranslated += 0.5
						}
					}
				}
			}

			totalTranslated := singularTranslated + pluralTranslated
			percentComplete := 0.0
			if totalStrings > 0 {
				percentComplete = float64(totalTranslated) / float64(totalStrings) * 100
			}

			missing := totalStrings - int(totalTranslated)

			fmt.Printf("  %s: %.1f/%d (%.1f%%) - %d missing\n",
				region, totalTranslated, totalStrings, percentComplete, missing)
		}

		fmt.Println()
	}

	// Show keys with missing translations
	fmt.Println("\nKeys with missing translations:")
	fmt.Println("-----------------------------")

	missingByLang := make(map[string][]string)

	for _, trans := range ts.Translations {
		// Skip if no English source
		if _, ok := trans.Translations["en"]; !ok {
			continue
		}

		for _, lang := range ts.Languages {
			if lang == "en" {
				continue
			}

			if trans.Type == models.TypeSingular {
				if langTrans, ok := trans.Translations[lang]; !ok || langTrans[models.QuantityOne] == "" {
					missingByLang[lang] = append(missingByLang[lang], fmt.Sprintf("%s (%s)", trans.Key, trans.Source))
				}
			} else {
				// For plurals, check both forms
				if langTrans, ok := trans.Translations[lang]; !ok ||
					langTrans[models.QuantityOne] == "" || langTrans[models.QuantityOther] == "" {
					missingByLang[lang] = append(missingByLang[lang], fmt.Sprintf("%s (plural, %s)", trans.Key, trans.Source))
				}
			}
		}
	}

	// Show missing translations for each language (limit to first 10)
	for _, baseLang := range sortedBaseLanguages {
		if baseLang == "en" {
			continue
		}

		if missing, ok := missingByLang[baseLang]; ok && len(missing) > 0 {
			fmt.Printf("%s: %d missing\n", baseLang, len(missing))

			// Show first 10 missing keys
			showCount := 10
			if len(missing) < showCount {
				showCount = len(missing)
			}

			for i := 0; i < showCount; i++ {
				fmt.Printf("  - %s\n", missing[i])
			}

			if len(missing) > showCount {
				fmt.Printf("  ... and %d more\n", len(missing)-showCount)
			}

			fmt.Println()
		}
	}
}

// autoTranslate uses AI to translate missing strings
func autoTranslate(serviceType string) {
	ts, err := csv.LoadFromCSV("")
	if err != nil {
		fmt.Printf("Error loading translations: %v\n", err)
		return
	}

	// Get translation service based on the specified type
	var service translator.TranslationService
	if serviceType == "azure" {
		// Force Azure service
		azureKey := os.Getenv("AZURE_TRANSLATOR_KEY")
		azureRegion := os.Getenv("AZURE_TRANSLATOR_REGION")
		azureEndpoint := os.Getenv("AZURE_TRANSLATOR_ENDPOINT")

		if azureKey == "" || azureRegion == "" {
			fmt.Println("Azure Translator credentials not set. Please set AZURE_TRANSLATOR_KEY and AZURE_TRANSLATOR_REGION environment variables.")
			return
		}

		if azureEndpoint == "" {
			azureEndpoint = "https://api.cognitive.microsofttranslator.com/translate"
		}

		service = &translator.AzureTranslator{
			Key:      azureKey,
			Region:   azureRegion,
			Endpoint: azureEndpoint,
		}
	} else if serviceType == "deepl" {
		// Force DeepL service
		deeplKey := os.Getenv("DEEPL_API_KEY")
		if deeplKey == "" {
			fmt.Println("DeepL API key not set. Please set DEEPL_API_KEY environment variable.")
			return
		}
		service = &translator.DeepLTranslator{APIKey: deeplKey}
	} else {
		// Use default service selection logic
		service = translator.GetTranslationService()
	}

	fmt.Println("Auto-translating missing strings...")

	// Perform translation
	count, err := translator.AutoTranslate(ts, service)
	if err != nil {
		fmt.Printf("Error during translation: %v\n", err)
		return
	}

	// Save updated translations
	if err := csv.SaveToCSV(ts, ""); err != nil {
		fmt.Printf("Error saving translations: %v\n", err)
		return
	}

	fmt.Printf("Auto-translation completed. Translated %d strings.\n", count)
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
