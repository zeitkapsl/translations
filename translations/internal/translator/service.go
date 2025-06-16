package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zeitkapsl/translations/internal/models"
)

// TranslationService interface for different translation providers
type TranslationService interface {
	Translate(text, sourceLang, targetLang string) (string, error)
	Name() string
}

// DeepLTranslator implements TranslationService for DeepL API
type DeepLTranslator struct {
	APIKey string
}

func (d *DeepLTranslator) Name() string {
	return "DeepL"
}

func (d *DeepLTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	if d.APIKey == "" {
		return "", fmt.Errorf("DeepL API key not set. Set DEEPL_API_KEY environment variable")
	}

	url := "https://api-free.deepl.com/v2/translate"

	requestBody, err := json.Marshal(map[string]interface{}{
		"text":        []string{text},
		"source_lang": strings.ToUpper(sourceLang),
		"target_lang": strings.ToUpper(targetLang),
	})
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+d.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(result.Translations) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result.Translations[0].Text, nil
}

// AzureTranslator implements TranslationService for Azure Translator API
type AzureTranslator struct {
	Key      string
	Region   string
	Endpoint string
}

func (a *AzureTranslator) Name() string {
	return "Azure Translator"
}

func (a *AzureTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	if a.Key == "" || a.Region == "" {
		return "", fmt.Errorf("Azure Translator credentials not set. Set AZURE_TRANSLATOR_KEY and AZURE_TRANSLATOR_REGION environment variables")
	}

	url := a.Endpoint
	if url == "" {
		url = "https://api.cognitive.microsofttranslator.com/translate"
	}

	// Add API version and parameters
	url = fmt.Sprintf("%s?api-version=3.0&from=%s&to=%s", url, sourceLang, targetLang)

	// Prepare request body
	requestBody, err := json.Marshal([]map[string]string{
		{"Text": text},
	})
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", a.Key)
	req.Header.Set("Ocp-Apim-Subscription-Region", a.Region)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result []struct {
		Translations []struct {
			Text string `json:"text"`
			To   string `json:"to"`
		} `json:"translations"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if len(result) == 0 || len(result[0].Translations) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	return result[0].Translations[0].Text, nil
}

// GetTranslationService returns an appropriate translation service based on environment variables
func GetTranslationService() TranslationService {
	// Check for Azure API key first
	azureKey := os.Getenv("AZURE_TRANSLATOR_KEY")
	azureRegion := os.Getenv("AZURE_TRANSLATOR_REGION")
	azureEndpoint := os.Getenv("AZURE_TRANSLATOR_ENDPOINT")

	// Default endpoint if not specified
	if azureEndpoint == "" {
		azureEndpoint = "https://api.cognitive.microsofttranslator.com/translate"
	}

	// If Azure credentials are available, use Azure
	if azureKey != "" && azureRegion != "" {
		fmt.Println("Using Azure Translator service")
		return &AzureTranslator{
			Key:      azureKey,
			Region:   azureRegion,
			Endpoint: azureEndpoint,
		}
	}

	// Fall back to DeepL
	deeplKey := os.Getenv("DEEPL_API_KEY")
	if deeplKey != "" {
		fmt.Println("Using DeepL translation service")
		return &DeepLTranslator{APIKey: deeplKey}
	}

	// No translation service available
	return nil
}

// AutoTranslateFromEnglish performs automatic translation from English only
func AutoTranslateFromEnglish(tm *models.TranslationManager, service TranslationService) (int, error) {
	if service == nil {
		return 0, fmt.Errorf("no translation service configured")
	}

	sourceLang := "en"
	translatedCount := 0
	requestCount := 0

	fmt.Printf("Using %s for translation from English to all other languages\n", service.Name())

	// Get all target languages (exclude English and regional variants)
	targetLanguages := []string{}
	for _, lang := range tm.Languages {
		if lang != sourceLang && !strings.Contains(lang, "-") {
			targetLanguages = append(targetLanguages, lang)
		}
	}

	if len(targetLanguages) == 0 {
		return 0, fmt.Errorf("no target languages to translate to")
	}

	fmt.Printf("Target languages: %s\n", strings.Join(targetLanguages, ", "))

	// Debug: Count how many strings have English content
	englishStrings := 0
	for _, row := range tm.Translations {
		if sourceValues, hasEnglish := row.Values[sourceLang]; hasEnglish && (sourceValues.One != "" || sourceValues.Other != "") {
			englishStrings++
		}
	}
	fmt.Printf("Found %d strings with English content\n", englishStrings)

	for i, row := range tm.Translations {
		// Check if we have English source content
		sourceValues, hasEnglish := row.Values[sourceLang]
		if !hasEnglish || (sourceValues.One == "" && sourceValues.Other == "") {
			continue // Skip if no English content
		}

		for _, targetLang := range targetLanguages {
			// Rate limit: pause after every 70 requests
			if requestCount > 0 && requestCount%70 == 0 {
				fmt.Println("Rate limit reached, sleeping for 60 seconds...")
				time.Sleep(60 * time.Second)
			}

			// Check current target translation - get directly from Values, not with fallback
			targetValues, hasTarget := row.Values[targetLang]
			if !hasTarget {
				targetValues = models.TranslationValues{} // Empty if doesn't exist
			}

			if row.Type == "plural" {
				// Translate singular form
				if sourceValues.One != "" && targetValues.One == "" {
					fmt.Printf("Translating [en→%s] singular: %s\n", targetLang, sourceValues.One)
					translatedText, err := service.Translate(sourceValues.One, sourceLang, targetLang)
					requestCount++

					if err != nil {
						fmt.Printf("Error translating singular to %s: %v\n", targetLang, err)
					} else {
						if tm.Translations[i].Values == nil {
							tm.Translations[i].Values = make(map[string]models.TranslationValues)
						}
						values := tm.Translations[i].Values[targetLang]
						values.One = translatedText
						tm.Translations[i].Values[targetLang] = values
						translatedCount++
					}
				}

				// Translate plural form
				if sourceValues.Other != "" && targetValues.Other == "" {
					fmt.Printf("Translating [en→%s] plural: %s\n", targetLang, sourceValues.Other)
					translatedText, err := service.Translate(sourceValues.Other, sourceLang, targetLang)
					requestCount++

					if err != nil {
						fmt.Printf("Error translating plural to %s: %v\n", targetLang, err)
					} else {
						if tm.Translations[i].Values == nil {
							tm.Translations[i].Values = make(map[string]models.TranslationValues)
						}
						values := tm.Translations[i].Values[targetLang]
						values.Other = translatedText
						tm.Translations[i].Values[targetLang] = values
						translatedCount++
					}
				}
			} else {
				// Singular translation
				if sourceValues.One != "" && targetValues.One == "" {
					fmt.Printf("Translating [en→%s]: %s\n", targetLang, sourceValues.One)
					translatedText, err := service.Translate(sourceValues.One, sourceLang, targetLang)
					requestCount++

					if err != nil {
						fmt.Printf("Error translating to %s: %v\n", targetLang, err)
					} else {
						if tm.Translations[i].Values == nil {
							tm.Translations[i].Values = make(map[string]models.TranslationValues)
						}
						tm.Translations[i].Values[targetLang] = models.TranslationValues{One: translatedText}
						translatedCount++
					}
				}
			}
		}
	}

	return translatedCount, nil
}
