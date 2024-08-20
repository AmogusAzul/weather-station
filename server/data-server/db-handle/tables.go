package dbhandle

import (
	"fmt"
	"reflect"
	"time"
)

const (
	StationTableName     string = "station"
	MeasurementTableName string = "measurement"
	EntryTableName       string = "entry"
)

type Table interface {
	GetTableName() string
	GetFieldsNames() []string
}

type Station struct {
	ID           int
	StationOwner string

	Latitude  float32
	Longitude float32
}

func (s Station) GetTableName() string { return StationTableName }
func (s Station) GetFieldsNames() []string {
	return []string{
		"station_id",
		"station_owner",
		"latitude",
		"longitude",
	}
}

type Measurement struct {
	ID int

	RandomNum int
}

func (m Measurement) GetTableName() string { return MeasurementTableName }
func (m Measurement) GetFieldsNames() []string {
	return []string{
		"measurement_id",
		"random_num",
	}
}

type Entry struct {
	ID        int
	StationID int
	Latitude  float32
	Longitude float32

	MeasurementID int

	EntryTime time.Time
}

func (e Entry) GetTableName() string { return EntryTableName }
func (e Entry) GetFieldsNames() []string {
	return []string{
		"entry_id",
		"station_id",
		"latitude",
		"longitude",
		"measurement_id",
		"entry_time",
	}
}

func GetID(object Table) int {
	val := reflect.ValueOf(object).FieldByName("ID")
	return int(val.Int())
}
func GetValues(object Table) []interface{} {
	val := reflect.ValueOf(object)
	values := []interface{}{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		values = append(values, field.Interface())
	}
	return values
}
func GetPointers(object Table) (pointers []interface{}, err error) {
	val := reflect.ValueOf(object)

	// Check if the object is a pointer to a struct
	if valStruct := val.Elem(); val.Kind() != reflect.Ptr ||
		valStruct.Kind() != reflect.Struct {
		return pointers, fmt.Errorf("object is meant to be a pointer of a struct")
	}
	// Select struct
	val = val.Elem()

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanAddr() {
			// Get the pointer to the field and store it in the slice
			pointers = append(pointers, field.Addr().Interface())
			continue
		}
		err = fmt.Errorf("field called %s isn't addreasable", val.Type().Field(i).Name)
	}

	return
}
