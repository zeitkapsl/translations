package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"path"
	"path/filepath"
)

type Module struct {
	App        string
	Path       string
	ImportFunc func(translations *Translations) error
	ExportFunc func(translations *Translations) error
}

func getModules(tm *Translations, basePath string) []Module {
	result := make([]Module, 0)

	// Android Module
	result = append(result, Module{
		App:  "android",
		Path: basePath,
		ImportFunc: func(tr *Translations) error {
			return ImportFromAndroid(tm, basePath)
		},
		ExportFunc: func(translations *Translations) error {
			return ExportToAndroid(tm, basePath)
		},
	})

	// iOS Module - ADDED THIS
	result = append(result, Module{
		App:  "ios",
		Path: filepath.Join(basePath, "ios", "Zeitkapsl", "Supporting Files"),
		ImportFunc: func(tr *Translations) error {
			return ImportFromXCStrings(tm, basePath)
		},
		ExportFunc: func(translations *Translations) error {
			return ExportToXCStrings(tm, basePath)
		},
	})

	// Server Emails Module
	serverMailsPath := path.Join(basePath, "server", "pkg", "mail", "templates")
	result = append(result, Module{
		App:  "server_emails",
		Path: serverMailsPath,
		ImportFunc: func(tr *Translations) error {
			return ImportFromJSON(tm, "server_emails", serverMailsPath)
		},
		ExportFunc: func(translations *Translations) error {
			return ExportToJson(tm, "server_emails", serverMailsPath)
		},
	})

	// Core Module
	coreTranslations := path.Join(basePath, "core", "pkg", "i18n")
	result = append(result, Module{
		App:  "core",
		Path: coreTranslations,
		ImportFunc: func(tr *Translations) error {
			return ImportFromJSON(tm, "core", coreTranslations)
		},
		ExportFunc: func(translations *Translations) error {
			return ExportToJson(tm, "core", coreTranslations)
		},
	})

	// Web Module
	webTranslations := path.Join(basePath, "web", "static", "translations")
	result = append(result, Module{
		App:  "web",
		Path: webTranslations,
		ImportFunc: func(tr *Translations) error {
			return ImportFromJSON(tm, "web", webTranslations)
		},
		ExportFunc: func(translations *Translations) error {
			return ExportToJson(tm, "web", webTranslations)
		},
	})

	return result
}

func main() {
	var basePath string
	var csvFile string

	rootCmd := &cobra.Command{
		Use:   "zeitkapsl-translations",
		Short: "Translation management tool for zeitkapsl",
		Long:  "A CLI tool to manage translations across iOS, Android, and Web platforms",
	}

	rootCmd.PersistentFlags().StringVar(&basePath, "base-path", "../", "Base path to the translation files")
	rootCmd.PersistentFlags().StringVar(&csvFile, "csv", "translations.csv", "CSV file path")

	tm := NewTranslations(basePath)
	modules := getModules(tm, basePath)

	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import translations from all platforms",
		Run: func(cmd *cobra.Command, args []string) {
			for _, module := range modules {
				fmt.Printf("Importing %s from %s\n", module.App, module.Path)
				err := module.ImportFunc(tm)
				if err != nil {
					fmt.Printf("Warning: Failed importing %s: %s\n", module.App, err.Error())
					// Continue with other modules instead of fatal error
					continue
				}
			}
			fmt.Printf("Saving to CSV: %s\n", csvFile)
			tm.Sort()
			if err := SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}
			fmt.Println("Import completed successfully!")
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

			tm := NewTranslations(basePath)
			if err := LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			tm.AddLanguage(lang)
			tm.Sort()
			if err := SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}
			fmt.Printf("Added language %s successfully!\n", lang)
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

			if err := LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}
			tm.AddRegion(region)
			if err := SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}
			fmt.Printf("Added region %s successfully!\n", region)
		},
	}
	addRegionCmd.Flags().String("region", "", "Region code to add (e.g., de-AT, en-US)")

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export translations to platform-specific formats",
		Run: func(cmd *cobra.Command, args []string) {
			if err := LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			platform, _ := cmd.Flags().GetString("platform")
			
			for _, m := range modules {
				if platform != "all" && m.App != platform {
					continue
				}
				
				fmt.Printf("Exporting %s to %s\n", m.App, m.Path)
				err := m.ExportFunc(tm)
				if err != nil {
					fmt.Printf("Warning: Failed to export %s: %s\n", m.App, err.Error())
					continue
				}
			}
			fmt.Println("Export completed successfully!")
		},
	}
	exportCmd.Flags().String("platform", "all", "Platform to export to (ios|android|web|core|server_emails|all)")

	// Auto-translate command
	autoTranslateCmd := &cobra.Command{
		Use:   "auto-translate",
		Short: "Auto-translate missing strings using AI (English as source)",
		Run: func(cmd *cobra.Command, args []string) {
			serviceType, _ := cmd.Flags().GetString("service")

			var service TranslationService
			switch serviceType {
			case "auto", "":
				service = GetTranslationService()
				if service == nil {
					fmt.Println("No translation service configured. Please set DEEPL_API_KEY or AZURE_TRANSLATOR_KEY/AZURE_TRANSLATOR_REGION environment variables.")
					return
				}
			default:
				log.Fatalf("Unknown service: %s. Use 'deepl', 'azure', or 'auto'", serviceType)
			}

			if err := LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			count, err := AutoTranslateFromEnglish(tm, service)
			if err != nil {
				log.Fatalf("Auto-translate failed: %v", err)
			}

			if err := SaveToCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to save CSV: %v", err)
			}

			fmt.Printf("Auto-translation completed. Translated %d strings from English.\n", count)
		},
	}
	autoTranslateCmd.Flags().String("service", "auto", "Translation service to use (auto|deepl|azure)")

	// Status command - NEW
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show translation status and statistics",
		Run: func(cmd *cobra.Command, args []string) {
			if err := LoadFromCSV(tm, csvFile); err != nil {
				log.Fatalf("Failed to load CSV: %v", err)
			}

			fmt.Printf("Translation Status:\n")
			fmt.Printf("==================\n")
			fmt.Printf("Total languages: %d\n", len(tm.Languages))
			fmt.Printf("Total translation keys: %d\n", len(tm.Translations))
			
			// Count by app
			appCounts := make(map[string]int)
			for _, trans := range tm.Translations {
				appCounts[trans.App]++
			}
			
			fmt.Printf("\nKeys by platform:\n")
			for app, count := range appCounts {
				fmt.Printf("  %s: %d keys\n", app, count)
			}
			
			fmt.Printf("\nLanguages: %v\n", tm.Languages)
		},
	}

	rootCmd.AddCommand(importCmd, addLangCmd, addRegionCmd, exportCmd, autoTranslateCmd, statusCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}