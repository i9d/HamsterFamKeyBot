package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var translations map[string]map[string]string

func loadTranslations() {
	translations = make(map[string]map[string]string)

	files := []string{"en.json", "ru.json", "uk.json"}

	for _, file := range files {
		lang := file[:2] // assuming file names like en.json, ru.json

		data, err := os.ReadFile(filepath.Join("lang", file))
		if err != nil {
			log.Fatalf("Error reading file %s: %v", file, err)
		}

		var translation map[string]string
		err = json.Unmarshal(data, &translation)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON from file %s: %v", file, err)
		}

		translations[lang] = translation
	}
}

func getTranslation(lang string, key string) string {
	return translations[lang][key]
}
