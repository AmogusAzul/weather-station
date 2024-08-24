package codec

import (
	"math"
	"net"
	"time"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/executer"
)

const (
	errorType        byte = 1
	validationError  byte = 1
	databaseError    byte = 2
	serverError      byte = 3
	versionError     byte = 4
	typeError        byte = 5
	noStationIDError byte = 6
	formatError      byte = 7

	okType              byte = 2
	noReturnOk          byte = 1
	noReturnValidatedOk byte = 2
	idReturnOk          byte = 3
)

type RequestHandler func(net.Conn, []byte, *executer.Executer) error

func StationHandler(conn net.Conn, content []byte, e *executer.Executer) (err error) {

	if len(content) < 4+e.Saver.TokenLength+4+4 {
		return ErrorAnswer(conn, formatError)
	}

	bytesUsed := 0
	stationID := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])
	bytesUsed += 4
	token := string(content[bytesUsed : bytesUsed+e.Saver.TokenLength])
	bytesUsed += e.Saver.TokenLength

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
func NewStationHandler(conn net.Conn, content []byte, e *executer.Executer) (err error) {

	bytesUsed := 0
	ownerNameLength := content[bytesUsed]
	bytesUsed += 1
	ownerName := string(content[bytesUsed : bytesUsed+int(ownerNameLength)])
	bytesUsed += int(ownerNameLength)
	latitude := math.Float32frombits(uint32(BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])))
	bytesUsed += 4
	longitude := math.Float32frombits(uint32(BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])))
	bytesUsed += 4

	stationID, err := e.DBHandler.SendRow(
		dbhandle.Station{
			StationOwner: ownerName,
			Latitude:     latitude,
			Longitude:    longitude,
		},
	)

	if err != nil {
		return ErrorAnswer(conn, databaseError)
	}

	return OkReturnAnswer(conn, idReturnOk, IntToBigEndianInt32Hex(stationID))
}

func CloseHandler(conn net.Conn, content []byte, e *executer.Executer) (err error) {

	go func() error {
		if err := e.CloseAll(); err != nil {
			return ErrorAnswer(conn, serverError)
		}

		return OkNoReturnAnswer(conn)
	}()

	return nil

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
