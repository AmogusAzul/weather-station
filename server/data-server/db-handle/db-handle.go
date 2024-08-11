package dbhandle

import (
	"database/sql"
	"fmt"
	"time"
)

type DbHandler struct {
	db *sql.DB
}

func GetDbHandler(dbUser, dbPassword, dbHost, dbPort, dbName string) (*DbHandler, error) {

	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		))

	return &DbHandler{
		db: db,
	}, err

}

func (dh *DbHandler) Exit() {
	dh.db.Close()
}

func (dh *DbHandler) ReadStation(station_id int) (Station, error) {

	station := Station{}

	err := dh.db.QueryRow(
		fmt.Sprintf("SELECT * FROM %s WHERE id = ?",
			StationTableName),
		station_id,
	).Scan(
		&station.StationID,
		&station.StationOwner,
		&station.Latitude,
		&station.Longitude,
	)

	if err != nil {

		//sql.ErrNoRows

		return station, err

	}

	return station, nil

}

func (dh *DbHandler) SendMeasurement(randomNum int) (Measurement, error) {

	result, err := dh.db.Exec(fmt.Sprintf("INSERT INTO %s (random_num) VALUES (?)", MeasurementTableName), randomNum)

	if err != nil {
		return Measurement{}, nil
	}

	measurementID, err := result.LastInsertId()

	return Measurement{
		MeasurementID: int(measurementID),
		RandomNum:     randomNum,
	}, err

}

func (dh *DbHandler) SendEntry(measurement Measurement, station Station, time time.Time) (err error) {

	query := fmt.Sprintf(
		"INSERT INTO %s (station_id, latitude, longitude, measurement_id, entry_time) VALUES (?, ?, ?, ?, ?)",
		EntryTableName)

	_, err = dh.db.Exec(query, station.StationID, station.Latitude, station.Longitude, measurement.MeasurementID, time)

	return err
}
