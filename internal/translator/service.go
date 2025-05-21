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

// TranslationService defines the interface for translation services
type TranslationService interface {
	Translate(text, sourceLang, targetLang string) (string, error)
	Name() string
}

// DeepLTranslator implements the TranslationService interface using DeepL API
type DeepLTranslator struct {
	APIKey string
}

// Name returns the name of the translation service
func (d *DeepLTranslator) Name() string {
	return "DeepL"
}

// Translate translates text using DeepL API
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

// GetTranslationService returns the DeepL translation service if the API key is set
func GetTranslationService() TranslationService {
	apiKey := os.Getenv("DEEPL_API_KEY")
	if apiKey == "" {
		panic("DEEPL_API_KEY is not set. Please set the environment variable for translation to work.")
	}
	return &DeepLTranslator{APIKey: apiKey}
}

// AutoTranslate translates missing strings in the translation set
func AutoTranslate(ts *models.TranslationSet, service TranslationService) (int, error) {
	sourceLang := "en"
	translatedCount := 0

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
