package password

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func ValidateStation(dh *dbhandle.DbHandler, stationID int, password string) bool {

	table, err := dh.ReadRowByID(stationID, &dbhandle.Station{})
	station, ok := table.(*dbhandle.Station)

	if err != nil || !ok {
		return false
	}
	return HashPassword(password) == station.PasswordHash
}

// Compares the provided password's hash to .ENV's one
func ValidateAdmin(password string) bool {
	adminHash := os.Getenv("ADMIN_HASH")
	return HashPassword(password) == adminHash
}
