package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ImportFromJSON(tm *Translations, app, baseDirectory string) error {
	// Check if the primary path exists
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		return fmt.Errorf("web directory not found at %s", baseDirectory)
	}

	return filepath.Walk(baseDirectory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			parts := strings.Split(filepath.Base(path), ".")
			if len(parts) < 2 {
				return nil

			}
			fmt.Println("processing: " + path)
			lang := parts[0]
			return importFromJavaScriptFile(tm, app, lang, path)
		}
		return nil
	})
}

func importFromJavaScriptFile(tm *Translations, app, lang, filePath string) error {
	fmt.Printf("Importing from JavaScript file: %s...\n", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	var parseData map[string]string

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&parseData)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}
	for k, v := range parseData {
		tm.SetTranslation(app, k, lang, v, "")
	}
	return nil
}

func ExportToJson(tm *Translations, app, baseDirectory string) error {
	tm.Sort()

	if err := os.MkdirAll(baseDirectory, 0755); err != nil {
		return err
	}

	translations := tm.GetTranslationsForApp(app)

	// Export each language as a separate JSON file
	for _, lang := range tm.Languages {
		object := make(map[string]string)
		containsValue := false
		for _, t := range translations {
			if v, ok := t.Values[lang]; ok && v != "" {
				object[t.Key] = v
				containsValue = true
			}
		}

		if containsValue {
			targetPath := path.Join(baseDirectory, lang+".json")
			f, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer f.Close()
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			err = enc.Encode(object)
			if err != nil {
				return err
			}
			fmt.Println("Create translation: " + targetPath)
		}
	}
	fmt.Printf("Exported web translations to separate JSON files in %s\n", baseDirectory)
	return nil
}
