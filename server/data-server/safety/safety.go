package safety

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

// #include "../../../common/safety/safety.h"
import "C"

var TOKEN_LENGTH = int(C.TOKEN_LENGTH)

type Saver struct {
	tokens map[int]string

	savePath string
}

func GetSaver(savePath string) *Saver {
	jsonData, err := os.ReadFile(savePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal the JSON data into a map[string]string
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

		// Decode the hex string back to the original binary data
		decodedValue, err := hex.DecodeString(value)
		if err != nil {
			log.Fatalf("Failed to decode hex string for key %d: %v", num, err)
		}

		tokens[num] = string(decodedValue)
	}

	return &Saver{
		tokens:   tokens,
		savePath: savePath,
	}
}

func (s *Saver) CreateToken(id int) (err error) {

	newToken, err := s.newToken()

	s.tokens[id] = newToken

	return

}

func (s *Saver) Validate(id int, token string) (valid bool, newToken string) {

	newToken = token

	// creates a new token if valid
	if s.tokens[id] == token {

		genToken, err := s.newToken()

		if err == nil {

			newToken = genToken
			s.tokens[id] = newToken
			valid = true

		}

	}

	return
}

func (s *Saver) newToken() (newToken string, err error) {

	b := make([]byte, TOKEN_LENGTH)

	_, err = rand.Read(b)

	newToken = string(b)

	return

}

func (s *Saver) Close() error {
	jsonTokens := make(map[string]string)
	for key, value := range s.tokens {
		// Encode the binary token to a hex string
		encodedValue := hex.EncodeToString([]byte(value))
		jsonTokens[fmt.Sprintf("%d", key)] = encodedValue
	}

	jsonData, err := json.MarshalIndent(jsonTokens, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling map to json: %v", err)
	}

	err = os.WriteFile(s.savePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("wasn't able to save tokens %s", err)
	}

	return nil
}
