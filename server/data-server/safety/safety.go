package safety

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Saver struct {
	tokens      map[int]string
	TokenLength int

	savePath string
}

func GetSaver(tokenLength int, savePath string) *Saver {

	jsonData, err := os.ReadFile(savePath)

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Step 2: Unmarshal the JSON data into a map[string]string
	jsonMap := make(map[string]string)
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	tokens := make(map[int]string)

	for key, value := range jsonMap {

		num, err := strconv.Atoi(key)

		if err != nil {
			log.Fatalf("Index \"%s\" in json isn't parseable", key)
		}

		tokens[num] = value
	}

	return &Saver{
		tokens:      tokens,
		TokenLength: tokenLength,
		savePath:    savePath,
	}

}

func (s *Saver) Validate(station_id int, token string) (valid bool, newToken string) {

	newToken = token

	// creates a new token
	if s.tokens[station_id] == token {

		valid = true

		b := make([]byte, s.TokenLength)

		_, err := rand.Read(b)

		if err == nil {
			newToken = string(b)
			s.tokens[station_id] = newToken
		}

	}

	return
}

func (s *Saver) Save() {

	jsonTokens := make(map[string]string)
	for key, value := range s.tokens {
		jsonTokens[fmt.Sprintf("%d", key)] = value
	}

	jsonData, err := json.MarshalIndent(jsonTokens, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling map to JSON: %v", err)
	}

	err = os.WriteFile(s.savePath, jsonData, 0644)

	if err != nil {
		log.Fatal("wasn't able to save tokens", err)
	}

}
