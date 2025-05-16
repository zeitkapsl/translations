package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zeitkapsl/translations/internal/csv"
	"github.com/zeitkapsl/translations/internal/models"
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
	langFlag := addLangCmd.String("lang", "", "Language code to add (e.g., es)")
	regionFlag := addRegionCmd.String("region", "", "Region code to add (e.g., de-AT)")
	platformFlag := exportCmd.String("platform", "", "Platform to export to (ios|android|json|csv)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse the command
	switch os.Args[1] {
	case "import":
		importCmd.Parse(os.Args[2:])
		importTranslations()
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
		exportTranslations(*platformFlag)
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
	fmt.Println("  import                        Import translations from all platforms")
	fmt.Println("  add-language --lang=<code>    Add a new language (e.g., es, fr)")
	fmt.Println("  add-region --region=<code>    Add a regional variant (e.g., de-AT)")
	fmt.Println("  export --platform=<platform>  Export to platform format (ios|android|json|csv)")
	fmt.Println("  stats                         Show translation statistics")
	fmt.Println("  auto-translate                Auto-translate missing strings")
	fmt.Println("  help                          Show this help")
}

// importTranslations imports translations from all platform-specific formats
func importTranslations() {
	// Create a new translation set
	ts := models.NewTranslationSet()
	
	fmt.Println("Importing translations...")
	
	// TODO: Implement platform-specific imports
	// For now, just created a sample translation example
	ts.AddOrUpdateTranslation("welcome_message", "Greeting on app start screen", "en", "Welcome!")
	ts.AddOrUpdateTranslation("photo_count.singular", "Shown when there's 1 photo", "en", "1 photo")
	ts.AddOrUpdateTranslation("photo_count.plural", "Shown when there are multiple photos", "en", "%d photos")
	
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
func exportTranslations(platform string) {
	fmt.Printf("Export to %s platform not yet implemented\n", platform)
	// TODO: Implement
}

// showStats shows translation statistics
func showStats() {
	fmt.Println("Statistics not yet implemented")
	// TODO: Implement
}

// autoTranslate uses AI to translate missing strings
func autoTranslate() {
	fmt.Println("Auto-translation not yet implemented")
	// TODO: Implement
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