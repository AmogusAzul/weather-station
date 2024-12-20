package safety

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	password "github.com/AmogusAzul/weather-station/server/data-server/password"
)

var TOKEN_LENGTH = int(6)

type Saver struct {
	tokens map[int]string

	savePath string
}

func GetSaver(savePath string, dh *dbhandle.DbHandler) *Saver {
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

	s := &Saver{
		tokens:   tokens,
		savePath: savePath,
	}

	stationCount, _ := dh.GetRowCountOf(dbhandle.Station{})

	for id := 0; id < stationCount; id++ {
		if tokens[id] != "" {
			s.CreateToken(id)
		}
	}

	return s
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

func (s *Saver) GetTokenByID(stationID int, rawPassword string, dh *dbhandle.DbHandler) string {
	if password.ValidateStation(dh, stationID, rawPassword) {
		return s.tokens[stationID]
	}

	result := ""
	for i := 0; i < TOKEN_LENGTH; i++ {

		result += "i"

	}
	return result
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
