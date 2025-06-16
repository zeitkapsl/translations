package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/zeitkapsl/translations/internal/android"
	"github.com/zeitkapsl/translations/internal/csv"
	"github.com/zeitkapsl/translations/internal/ios"
	"github.com/zeitkapsl/translations/internal/models"
	"github.com/zeitkapsl/translations/internal/translator"
	"github.com/zeitkapsl/translations/internal/web"
)

func main() {
	var basePath string
	var csvFile string

	rootCmd := &cobra.Command{
		Use:   "zeitkapsl-translations",
		Short: "Translation management tool for zeitkapsl",
		Long:  "A CLI tool to manage translations across iOS, Android, and Web platforms",
	}

	rootCmd.PersistentFlags().StringVar(&basePath, "base-path", ".", "Base path to the translation files")
	rootCmd.PersistentFlags().StringVar(&csvFile, "csv", "translations.csv", "CSV file path")

	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import translations from all platforms",
		Run: func(cmd *cobra.Command, args []string) {
			tm := models.NewTranslationManager(basePath)

			fmt.Println("Importing Android translations...")
			if err := android.ImportFromAndroid(tm, basePath); err != nil {
				log.Printf("Android import failed: %v", err)
			}

			fmt.Println("Importing iOS translations...")
			if err := ios.ImportFromXCStrings(tm, basePath); err != nil {
				log.Printf("iOS import failed: %v", err)
			}

			fmt.Println("Importing Web translations...")
			if err := web.ImportFromWeb(tm, basePath); err != nil {
				log.Printf("Web import failed: %v", err)
			}

			fmt.Printf("Saving to CSV: %s\n", csvFile)
			if err := csv.SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}

			fmt.Println("Import completed successfully!")
			tm.GetStats()
		},
	}

	// Add language command
	addLangCmd := &cobra.Command{
		Use:   "add-language",
		Short: "Add a new language",
		Run: func(cmd *cobra.Command, args []string) {
			lang, _ := cmd.Flags().GetString("lang")
			if lang == "" {
				log.Fatal("--lang flag is required")
			}

			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			tm.AddLanguage(lang)
			if err := csv.SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}
		},
	}
	addLangCmd.Flags().String("lang", "", "Language code to add (e.g., es, fr)")

	// Add region command
	addRegionCmd := &cobra.Command{
		Use:   "add-region",
		Short: "Add a new region to a language",
		Run: func(cmd *cobra.Command, args []string) {
			region, _ := cmd.Flags().GetString("region")
			if region == "" {
				log.Fatal("--region flag is required")
			}

			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			tm.AddRegion(region)
			if err := csv.SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}
		},
	}
	addRegionCmd.Flags().String("region", "", "Region code to add (e.g., de-AT, en-US)")

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export translations to platform-specific formats",
		Run: func(cmd *cobra.Command, args []string) {
			platform, _ := cmd.Flags().GetString("platform")

			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			switch platform {
			case "ios":
				fmt.Println("Exporting to iOS...")
				if err := ios.ExportToXCStrings(tm, basePath); err != nil {
					log.Fatalf("iOS export failed: %v", err)
				}
			case "android":
				fmt.Println("Exporting to Android...")
				if err := android.ExportToAndroid(tm, basePath); err != nil {
					log.Fatalf("Android export failed: %v", err)
				}
			case "json", "web":
				fmt.Println("Exporting to Web...")
				if err := web.ExportToWeb(tm, basePath); err != nil {
					log.Fatalf("Web export failed: %v", err)
				}
			case "csv":
				fmt.Printf("Exporting to CSV: %s\n", csvFile)
				if err := csv.SaveToCSV(tm, csvFile); err != nil {
					log.Fatalf("CSV export failed: %v", err)
				}
			case "all", "":
				fmt.Println("Exporting to all platforms...")
				if err := ios.ExportToXCStrings(tm, basePath); err != nil {
					log.Printf("iOS export failed: %v", err)
				}
				if err := android.ExportToAndroid(tm, basePath); err != nil {
					log.Printf("Android export failed: %v", err)
				}
				if err := web.ExportToWeb(tm, basePath); err != nil {
					log.Printf("Web export failed: %v", err)
				}
			default:
				log.Fatalf("Unknown platform: %s", platform)
			}

			fmt.Println("Export completed successfully!")
		},
	}
	exportCmd.Flags().String("platform", "all", "Platform to export to (ios|android|json|csv|all)")

	// Stats command
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show translation statistics",
		Run: func(cmd *cobra.Command, args []string) {
			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}
			tm.GetStats()
		},
	}

	// Detailed stats command
	detailedStatsCmd := &cobra.Command{
		Use:   "detailed-stats",
		Short: "Show detailed translation statistics with source analysis",
		Run: func(cmd *cobra.Command, args []string) {
			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}
			tm.GetDetailedStats()
		},
	}

	// Auto-translate command
	autoTranslateCmd := &cobra.Command{
		Use:   "auto-translate",
		Short: "Auto-translate missing strings using AI (English as source)",
		Run: func(cmd *cobra.Command, args []string) {
			serviceType, _ := cmd.Flags().GetString("service")

			var service translator.TranslationService
			switch serviceType {
			// case "deepl":
			// 	deeplKey := os.Getenv("DEEPL_API_KEY")
			// 	if deeplKey == "" {
			// 		log.Fatal("DeepL API key not set. Please set DEEPL_API_KEY environment variable.")
			// 	}
			// 	service = translator.NewDeepLTranslator(deeplKey)
			// case "azure":
			// 	azureKey := os.Getenv("AZURE_TRANSLATOR_KEY")
			// 	azureRegion := os.Getenv("AZURE_TRANSLATOR_REGION")
			// 	azureEndpoint := os.Getenv("AZURE_TRANSLATOR_ENDPOINT")
			// 	if azureKey == "" || azureRegion == "" {
			// 		log.Fatal("Azure Translator credentials not set. Please set AZURE_TRANSLATOR_KEY and AZURE_TRANSLATOR_REGION environment variables.")
			// 	}
			// 	service = translator.NewAzureTranslator(azureKey, azureRegion, azureEndpoint)
			case "auto", "":
				service = translator.GetTranslationService()
				if service == nil {
					fmt.Println("No translation service configured. Please set DEEPL_API_KEY or AZURE_TRANSLATOR_KEY/AZURE_TRANSLATOR_REGION environment variables.")
					return
				}

			default:
				log.Fatalf("Unknown service: %s. Use 'deepl', 'azure', or 'auto'", serviceType)
			}

			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			count, err := translator.AutoTranslateFromEnglish(tm, service)
			if err != nil {
				log.Fatalf("Auto-translate failed: %v", err)
			}

			if err := csv.SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}

			fmt.Printf("Auto-translation completed. Translated %d strings from English.\n", count)
		},
	}
	autoTranslateCmd.Flags().String("service", "auto", "Translation service to use (deepl|azure|auto)")

	// Cleanup command
	cleanupCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up redundant language variants",
		Run: func(cmd *cobra.Command, args []string) {
			tm := models.NewTranslationManager(basePath)
			if err := csv.LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			cleaned := tm.CleanupLanguages()

			if err := csv.SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}

			fmt.Printf("Cleanup completed. Merged %d language variants.\n", cleaned)
			tm.GetStats()
		},
	}

	rootCmd.AddCommand(importCmd, addLangCmd, addRegionCmd, exportCmd, statsCmd, detailedStatsCmd, autoTranslateCmd, cleanupCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
