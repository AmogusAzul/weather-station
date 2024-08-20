package codec

import (
	"net"
	"time"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

const (
	okValidatedType  byte = 2
	errorType        byte = 1
	validationError  byte = 1
	databaseError    byte = 2
	serverError      byte = 3
	versionError     byte = 4
	typeError        byte = 5
	noStationIDError byte = 6
)

type RequestHandler func(net.Conn, []byte, *safety.Saver, *dbhandle.DbHandler) error

func StationHandler(conn net.Conn, content []byte, saver *safety.Saver, dbHandler *dbhandle.DbHandler) (err error) {

	bytesUsed := 0
	stationID := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])
	bytesUsed += 4
	token := string(content[bytesUsed : bytesUsed+saver.TokenLength])
	bytesUsed += saver.TokenLength

	valid, newToken := saver.Validate(stationID, token)

	if !valid {
		return ErrorValidatedAnswer(conn, validationError, newToken)
	}
	moment := time.Unix(int64(BigEndianInt32HexToInt(content[bytesUsed:bytesUsed+4])), 0)
	bytesUsed += 4
	randomNum := BigEndianInt32HexToInt(content[bytesUsed : bytesUsed+4])

	table, err := dbHandler.ReadRowByID(stationID, &dbhandle.Station{})
	if err != nil {
		return ErrorValidatedAnswer(conn, noStationIDError, newToken)
	}
	station, ok := table.(*dbhandle.Station)
	if !ok {
		return ErrorValidatedAnswer(conn, serverError, newToken)
	}

	measurementID, err := dbHandler.SendRow(dbhandle.Measurement{RandomNum: randomNum})
	if err != nil {
		return ErrorValidatedAnswer(conn, databaseError, newToken)
	}
	_, err = dbHandler.SendRow(dbhandle.Entry{
		StationID:     stationID,
		Latitude:      station.Latitude,
		Longitude:     station.Longitude,
		MeasurementID: measurementID,
		EntryTime:     moment,
	})
	if err != nil {
		return ErrorValidatedAnswer(conn, databaseError, newToken)
	}

	return OkValidatedAnswer(conn, newToken)

}
func NewStationHandler(conn net.Conn, content []byte, saver *safety.Saver, dbHandler *dbhandle.DbHandler) (err error) {

	return
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

func OkValidatedAnswer(conn net.Conn, newToken string) (err error) {

	defer conn.Close()
	answer := []byte{version, okValidatedType}
	for i := 0; i < len(newToken); i++ {
		answer = append(answer, newToken[i])
	}
	_, err = conn.Write(answer)

	return err

}

func BigEndianInt32HexToInt(content []byte) int {

	return int(content[0])<<24 |
		int(content[1])<<16 |
		int(content[2])<<8 |
		int(content[3])
}
