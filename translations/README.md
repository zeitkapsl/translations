# translations


## About zeitkapsl

[zeitkapsl](https://zeitkapsl.eu) is a privacy-first end-to-end encrypted photo storage and sharing app for people who care about digital autonomy and data protection. 

We support multiple platforms—iOS, Android, desktop, and web—and keep all data end-to-end encrypted. As we expand to more European markets, localization is becoming a key focus area.

Our ultimate goal is to support all spoken languages in Europe.

---

## Project Overview

Each of our supported platforms uses a different localization format:

- **iOS**: `.xcstrings`
- **Android**: `values/<lang>/strings.xml`
- **Web/Backend**: `<lang>.json`

Managing these translations independently is time-consuming and error-prone. We want to build a **common translation workflow** that centralizes all strings in one editable format and supports easy export back to platform-specific files.

This project is ideal for an internship: it’s well-scoped, technically interesting, and directly impacts the success of our European market expansion.

---

## Project Goals / Requirements

Implement a CLI tool in Golang that helps us translating zeitkapsl into multiple languages for all platforms at once.

The backing storage format will be CSV since it can be read by almost any editor/Excel/Libre-Office

#### Import Translations from

- iOS `.xcstrings`
- Android `strings.xml`
- Web/JSON translation files

Example en.json
```
{
 "some.key": "Some Value"
 "photo_count.singular": "1 Photo"
 "photo_count.plural": "%d Photos"
}
```


### Normalize to common CSV format

 - Context and comments
- Placeholder variables
- Singular/plural rules
   - some.key.singular
   - some.key.plural
- Multiple languages

Sample CSV: 

| key                  | comment                              | de                    | de-AT             | de-DE                 | en               | en-US        | es              | fr                     |
| -------------------- | ------------------------------------ | --------------------- | ----------------- | --------------------- | ---------------- | ------------ | --------------- | ---------------------- |
| welcome_message      | Greeting on app start screen         | Willkommen!           | Guten Tag!        | Willkommen!           | Welcome!         | Hey there!   | ¡Bienvenido!    | Bienvenue !            |
| photo_count.singular | Shown when there's 1 photo           | 1 Foto                | 1 Foto            | 1 Foto                | 1 photo          | 1 picture    | 1 foto          | 1 photo                |
| photo_count.plural   | Shown when there are multiple photos | %d Fotos              | %d Fotos          | %d Fotos              | %d photos        | %d pictures  | %d fotos        | %d photos              |
| upload_success       | Confirmation message after upload    | Upload abgeschlossen. | Hochladen fertig. | Upload abgeschlossen. | Upload complete. | Upload done. | Carga completa. | Téléversement terminé. |


### Export CSV back into platform native formats

- Only keys in the imported source in english are present in the exported files
- Placeholder consistency across all languages
- Fallbacks for regional variants (e.g. `de-AT` → `de`)
- **Export translated CSV** back into:
   - `.xcstrings` for iOS
   - `strings.xml` for Android
   - `.json` for web/backend

### CLI (Command Line Interface)

- Implement a command line interface that supports the following commands
	- import: imports available strings in english from all platforms
	- add-language --lang=es -> adds a new language
		- make sure language columns are sorted asc in the CSV
	- add-region --region=de-AT -> adds a new region to a language
		- make sure regions appear next to the language columns
	- export --platform=ios|android|json|csv: exports all strings/xcstrings/json files again to the native platforms
	- stats: shows 
		- available languages 
		- total number of strings
		- total number of missing strings (compared to english)
	- auto-translate: auto translates using DeepL or Chat GPT all missing language strings (not regions)

### AI Translation Support
- Implement autocomplete support using Chat GPT/DeepL or similar suitable AI Tools to fill in suggestions for missing translations.

### Web UI (optional)
- Implement a locally runnable Web UI in Svelte (Kit)
- It allows to easily fill in missing translations for languages
- All changes should be reflected in the backing CSV file
- Views
	- Give a Dashboard that shows all supported languages and the number of missing items for each language
	- Add a Button that allows to add a new language
	- Add a Button to add a new region to a given language
- Add a Button to Generate all platform files for all languages



## What You'll Learn

- Cross-platform software architecture
- Real-world internationalization (i18n) practices
- Working with structured data formats (XML, JSON, CSV)
- Handling plurals, placeholders, and fallbacks across languages
- Contributing to a real privacy-focused app used in production
- Integrating third part AI services using REST
