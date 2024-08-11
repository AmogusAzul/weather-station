package dbhandle

import "time"

const (
	StationTableName     string = "station"
	MeasurementTableName string = "measurement"
	EntryTableName       string = "entry"
)

type Station struct {
	StationID    int
	StationOwner string

	Latitude  float32
	Longitude float32
}

type Measurement struct {
	MeasurementID int

	RandomNum int
}

type Entry struct {
	StationID int
	Latitude  float32
	Longitude float32

	MeasurementID int

	EntryTime time.Time
}
