package codec

import (
	"math"
	"net"
	"time"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/executer"
	"github.com/AmogusAzul/weather-station/server/data-server/password"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

const (
	errorType byte = 1

	validationError  byte = 1
	databaseError    byte = 2
	serverError      byte = 3
	versionError     byte = 4
	typeError        byte = 5
	noStationIDError byte = 6
	formatError      byte = 7

	okType byte = 2

	noReturnOk          byte = 1
	noReturnValidatedOk byte = 2
	idReturnOk          byte = 3
)

type RequestHandler func(net.Conn, byte, []byte, *executer.Executer) error

func StationHandler(conn net.Conn, specificRequestType byte, content []byte, e *executer.Executer) (err error) {

	if len(content) < 4+safety.TOKEN_LENGTH+4+4 {
		return ErrorAnswer(conn, formatError)
	}

	bytesUsed := 0

	stationID := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])
	bytesUsed += 4
	token := string(content[bytesUsed : bytesUsed+safety.TOKEN_LENGTH])
	bytesUsed += safety.TOKEN_LENGTH

	valid, newToken := e.Saver.Validate(stationID, token)

	if !valid {
		return ErrorValidatedAnswer(conn, validationError, newToken)
	}
	moment := time.Unix(int64(BigEndianInt32HexToInt(content[bytesUsed:bytesUsed+4])), 0)
	bytesUsed += 4
	randomNum := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])

	table, err := e.DBHandler.ReadRowByID(stationID, &dbhandle.Station{})
	if err != nil {
		return ErrorValidatedAnswer(conn, noStationIDError, newToken)
	}
	station, ok := table.(*dbhandle.Station)
	if !ok {
		return ErrorValidatedAnswer(conn, serverError, newToken)
	}

	measurementID, err := e.DBHandler.SendRow(dbhandle.Measurement{RandomNum: randomNum})
	if err != nil {
		return ErrorValidatedAnswer(conn, databaseError, newToken)
	}
	_, err = e.DBHandler.SendRow(dbhandle.Entry{
		StationID:     stationID,
		Latitude:      station.Latitude,
		Longitude:     station.Longitude,
		MeasurementID: measurementID,
		EntryTime:     moment,
	})
	if err != nil {
		return ErrorValidatedAnswer(conn, databaseError, newToken)
	}

	return OkNoReturnValidatedAnswer(conn, newToken)

}
func NewStationHandler(conn net.Conn, specificRequestType byte, content []byte, e *executer.Executer) (err error) {

	bytesUsed := 0
	ownerNameLength := content[bytesUsed]
	bytesUsed += 1
	ownerName := string(content[bytesUsed : bytesUsed+int(ownerNameLength)])
	bytesUsed += int(ownerNameLength)
	latitude := math.Float32frombits(uint32(BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])))
	bytesUsed += 4
	longitude := math.Float32frombits(uint32(BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])))
	bytesUsed += 4
	stationPasswordLength := content[bytesUsed]
	bytesUsed += 1
	stationPassword := string(content[bytesUsed : bytesUsed+int(stationPasswordLength)])
	bytesUsed += int(stationPasswordLength)

	if !password.ValidateAdmin(string(content[bytesUsed+1:])) {
		return ErrorAnswer(conn, validationError)
	}

	stationID, err := e.DBHandler.SendRow(
		dbhandle.Station{
			StationOwner: ownerName,
			Latitude:     latitude,
			Longitude:    longitude,
			PasswordHash: password.HashPassword(stationPassword),
		},
	)

	if err != nil {
		return ErrorAnswer(conn, databaseError)
	}

	e.Saver.CreateToken(stationID)

	return OkReturnAnswer(conn, idReturnOk, IntToBigEndianInt32Hex(stationID))
}

func CloseHandler(conn net.Conn, specificRequestType byte, content []byte, e *executer.Executer) (err error) {

	if !password.ValidateAdmin(string(content[1:])) {
		return ErrorAnswer(conn, validationError)
	}
	go func() error {
		if err := e.CloseAll(); err != nil {
			return ErrorAnswer(conn, serverError)
		}

		return OkNoReturnAnswer(conn)
	}()
	return nil
}

func GetTokenHandler(conn net.Conn, specificRequestType byte, content []byte, e *executer.Executer) (err error) {

	bytesUsed := 0

	stationID := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])
	bytesUsed += 4
	providedPasswordLength := content[bytesUsed]
	bytesUsed += 1
	providedPassword := string(content[bytesUsed : bytesUsed+int(providedPasswordLength)])
	bytesUsed += int(providedPasswordLength)

	token := e.Saver.GetTokenByID(stationID, providedPassword, e.DBHandler)

	return OkNoReturnValidatedAnswer(conn, token)
}

func ErrorAnswer(conn net.Conn, specificErrorType byte) (err error) {

	defer conn.Close()
	_, err = conn.Write([]byte{version, errorType, specificErrorType})
	return err
}

func ErrorValidatedAnswer(conn net.Conn, specificErrorType byte, token string) (err error) {

	defer conn.Close()
	answer := []byte{version, errorType, specificErrorType}
	for i := 0; i < len(token); i++ {
		answer = append(answer, token[i])
	}
	_, err = conn.Write(answer)
	return err
}

func OkNoReturnValidatedAnswer(conn net.Conn, newToken string) (err error) {

	defer conn.Close()
	answer := []byte{version, okType, noReturnValidatedOk}
	for i := 0; i < len(newToken); i++ {
		answer = append(answer, newToken[i])
	}
	_, err = conn.Write(answer)

	return err

}

func OkReturnValidatedAnswer(conn net.Conn, specificOkType byte, returnContent []byte, newToken string) (err error) {

	defer conn.Close()
	answer := []byte{version, okType, specificOkType}
	for i := 0; i < len(newToken); i++ {
		answer = append(answer, newToken[i])
	}
	answer = append(answer, returnContent...)
	_, err = conn.Write(answer)
	return
}

func OkReturnAnswer(conn net.Conn, specificOkType byte, returnContent []byte) (err error) {

	defer conn.Close()
	answer := []byte{version, okType, specificOkType}
	answer = append(answer, returnContent...)
	_, err = conn.Write(answer)
	return
}

func OkNoReturnAnswer(conn net.Conn) (err error) {

	defer conn.Close()

	_, err = conn.Write([]byte{version})

	return
}

func BigEndianInt32HexToInt(content []byte) int {

	return int(content[0])<<24 |
		int(content[1])<<16 |
		int(content[2])<<8 |
		int(content[3])
}

func IntToBigEndianInt32Hex(num int) []byte {
	// Convert the int to an int32 (assuming the input will fit in int32)
	num32 := int32(num)

	// Manually extract each byte using bit manipulation
	bytes := []byte{
		byte((num32 >> 24) & 0xFF), // Most significant byte
		byte((num32 >> 16) & 0xFF),
		byte((num32 >> 8) & 0xFF),
		byte(num32 & 0xFF), // Least significant byte
	}
	return bytes
}
