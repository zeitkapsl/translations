package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeitkapsl/translations/internal/android"
	"github.com/zeitkapsl/translations/internal/csv"
	"github.com/zeitkapsl/translations/internal/ios"
	"github.com/zeitkapsl/translations/internal/models"
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
	importDirFlag := importCmd.String("dir", ".", "Directory to import from")
	langFlag := addLangCmd.String("lang", "", "Language code to add (e.g., es)")
	regionFlag := addRegionCmd.String("region", "", "Region code to add (e.g., de-AT)")
	platformFlag := exportCmd.String("platform", "", "Platform to export to (ios|android|json|csv)")
	exportDirFlag := exportCmd.String("dir", ".", "Directory to export to")

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
		autoTranslate()
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
	fmt.Println("                                   Export to platform format (ios|android|json|csv)")
	fmt.Println("  stats                            Show translation statistics")
	fmt.Println("  auto-translate                   Auto-translate missing strings")
	fmt.Println("  help                             Show this help")
}

// importTranslations imports translations from all platform-specific formats
func importTranslations(directory string) {
	// Create a new translation set or load existing one if available
	var ts *models.TranslationSet
	//var err error

	// Check if translations.csv already exists
	if _, err := os.Stat(csv.DefaultCSVFile); err == nil {
		fmt.Println("Found existing translations.csv, loading...")
		ts, err = csv.LoadFromCSV("")
		if err != nil {
			fmt.Printf("Error loading existing translations: %v\n", err)
			ts = models.NewTranslationSet()
		}
	} else {
		ts = models.NewTranslationSet()
	}

	fmt.Println("Importing translations...")

	// Import from iOS .xcstrings
	if err := ios.ImportFromXCStrings(ts, directory); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Import from Android strings.xml
	if err := android.ImportFromAndroid(ts, directory); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Import from Web/Backend JSON files
	if err := web.ImportFromJSON(ts, directory); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Save to CSV
	if err := csv.SaveToCSV(ts, ""); err != nil {
		fmt.Printf("Error saving translations: %v\n", err)
		return
	}

	fmt.Println("Import completed successfully")
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
		if err := ios.ExportToXCStrings(ts, filepath.Join(directory, "Localizable.xcstrings")); err != nil {
			fmt.Printf("Error exporting to iOS format: %v\n", err)
		}
	case "android":
		if err := android.ExportToAndroid(ts, directory); err != nil {
			fmt.Printf("Error exporting to Android format: %v\n", err)
		}
	case "json":
		if err := web.ExportToJSON(ts, directory); err != nil {
			fmt.Printf("Error exporting to JSON format: %v\n", err)
		}
	case "csv":
		// CSV is already saved, just notify
		fmt.Println("CSV format is the source format, no export needed")
	default:
		fmt.Printf("Unknown platform: %s. Expected ios, android, json, or csv\n", platform)
	}
}

// showStats shows translation statistics
func showStats() {
	ts, err := csv.LoadFromCSV("")
	if err != nil {
		fmt.Printf("Error loading translations: %v\n", err)
		return
	}

	fmt.Println("Translation Statistics:")
	fmt.Println("======================")

	// Count total strings (based on English translations)
	totalStrings := 0
	for _, trans := range ts.Translations {
		if enTrans, ok := trans.Translations["en"]; ok && len(enTrans) > 0 {
			totalStrings++
		}
	}

	fmt.Printf("Total number of strings: %d\n\n", totalStrings)

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

	// Calculate stats for each language
	for baseLang, regions := range baseLanguages {
		if baseLang == "en" {
			continue // Skip English as it's the source language
		}

		// Base language stats
		translated := 0
		for _, trans := range ts.Translations {
			if langTrans, ok := trans.Translations[baseLang]; ok && len(langTrans) > 0 {
				translated++
			}
		}

		missing := totalStrings - translated
		percentComplete := 0.0
		if totalStrings > 0 {
			percentComplete = float64(translated) / float64(totalStrings) * 100
		}

		fmt.Printf("%s: %d/%d (%.1f%%) - %d missing\n",
			baseLang, translated, totalStrings, percentComplete, missing)

		// Regional variants stats
		for _, region := range regions {
			translated := 0
			for _, trans := range ts.Translations {
				if langTrans, ok := trans.Translations[region]; ok && len(langTrans) > 0 {
					translated++
				}
			}

			missing := totalStrings - translated
			percentComplete := 0.0
			if totalStrings > 0 {
				percentComplete = float64(translated) / float64(totalStrings) * 100
			}

			fmt.Printf("  %s: %d/%d (%.1f%%) - %d missing\n",
				region, translated, totalStrings, percentComplete, missing)
		}
	}
}

// autoTranslate uses AI to translate missing strings
func autoTranslate() {
	fmt.Println("Auto-translation not yet implemented")
	// TODO: Implement AI translation integration
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
