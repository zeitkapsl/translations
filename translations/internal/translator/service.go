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

type TranslationService interface {
	Translate(text, sourceLang, targetLang string) (string, error)
	Name() string
}

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

type AzureTranslator struct {
	Key        string
	Region     string
	Endpoint   string
}

func (a *AzureTranslator) Name() string {
	return "Azure Translator"
}

// Translate translates text using Azure Translator API
func (a *AzureTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	if a.Key == "" {
		return "", fmt.Errorf("azure Translator key not set. Set AZURE_TRANSLATOR_KEY environment variable")
	}
	
	if a.Region == "" {
		return "", fmt.Errorf("azure Translator region not set. Set AZURE_TRANSLATOR_REGION environment variable")
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
	panic("No translation service configured. Please set either DEEPL_API_KEY or both AZURE_TRANSLATOR_KEY and AZURE_TRANSLATOR_REGION environment variables.")
}


func AutoTranslate(ts *models.TranslationSet, service TranslationService) (int, error) {
	sourceLang := "en"
	translatedCount := 0
	requestCount := 0 //counter for rate limiting

	fmt.Printf("Using %s for translation\n", service.Name())

	baseLanguages := []string{}
	for _, lang := range ts.Languages {
		if !strings.Contains(lang, "-") && lang != sourceLang {
			baseLanguages = append(baseLanguages, lang)
		}
	}

	if len(baseLanguages) == 0 {
		return 0, fmt.Errorf("no target languages to translate to")
	}

	for i, trans := range ts.Translations {
		sourceTranslations, ok := trans.Translations[sourceLang]
		if !ok || len(sourceTranslations) == 0 {
			continue
		}

		for _, targetLang := range baseLanguages {
			// Rate limit: pause after every 70 requests
			if requestCount > 0 && requestCount%70 == 0 {
				fmt.Println("Rate limit reached, sleeping for 60 seconds...")
				time.Sleep(60 * time.Second)
			}

			if trans.Type == models.TypeSingular {
				if _, ok := trans.Translations[targetLang]; ok && len(trans.Translations[targetLang]) > 0 {
					continue
				}

				sourceText := sourceTranslations[models.QuantityOne]
				if sourceText == "" {
					continue
				}

				fmt.Printf("Translating [%s]: %s\n", targetLang, sourceText)
				translatedText, err := service.Translate(sourceText, sourceLang, targetLang)
				requestCount++

				if err != nil {
					fmt.Printf("Error translating to %s: %v\n", targetLang, err)
					continue
				}

				if ts.Translations[i].Translations == nil {
					ts.Translations[i].Translations = make(map[string]map[models.TranslationQuantity]string)
				}
				if ts.Translations[i].Translations[targetLang] == nil {
					ts.Translations[i].Translations[targetLang] = make(map[models.TranslationQuantity]string)
				}
				ts.Translations[i].Translations[targetLang][models.QuantityOne] = translatedText
				translatedCount++

			} else {
				targetTranslations, ok := trans.Translations[targetLang]
				if !ok {
					targetTranslations = make(map[models.TranslationQuantity]string)
					if ts.Translations[i].Translations == nil {
						ts.Translations[i].Translations = make(map[string]map[models.TranslationQuantity]string)
					}
					ts.Translations[i].Translations[targetLang] = targetTranslations
				}

				if sourceOne, ok := sourceTranslations[models.QuantityOne]; ok && sourceOne != "" {
					if _, ok := targetTranslations[models.QuantityOne]; !ok || targetTranslations[models.QuantityOne] == "" {
						fmt.Printf("Translating [%s] singular: %s\n", targetLang, sourceOne)
						translatedOne, err := service.Translate(sourceOne, sourceLang, targetLang)
						requestCount++
						if err != nil {
							fmt.Printf("Error translating singular to %s: %v\n", targetLang, err)
						} else {
							ts.Translations[i].Translations[targetLang][models.QuantityOne] = translatedOne
							translatedCount++
						}
					}
				}

				if sourceOther, ok := sourceTranslations[models.QuantityOther]; ok && sourceOther != "" {
					if _, ok := targetTranslations[models.QuantityOther]; !ok || targetTranslations[models.QuantityOther] == "" {
						fmt.Printf("Translating [%s] plural: %s\n", targetLang, sourceOther)
						translatedOther, err := service.Translate(sourceOther, sourceLang, targetLang)
						requestCount++
						if err != nil {
							fmt.Printf("Error translating plural to %s: %v\n", targetLang, err)
						} else {
							ts.Translations[i].Translations[targetLang][models.QuantityOther] = translatedOther
							translatedCount++
						}
					}
				}
			}
		}
	}

	return translatedCount, nil
}

