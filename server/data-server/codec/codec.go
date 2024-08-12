package codec

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

type Decoder struct {
	WorkerQueue chan net.Conn
	WorkerPool  chan<- chan net.Conn
	killChan    chan bool

	dbHandler *dbhandle.DbHandler
	saver     *safety.Saver
}

const (
	version byte = 1
)

func GetDecoder(workerPool chan chan net.Conn, dbHandler *dbhandle.DbHandler, saver *safety.Saver) *Decoder {

	return &Decoder{
		WorkerQueue: make(chan net.Conn),
		WorkerPool:  workerPool,
		killChan:    make(chan bool),
		dbHandler:   dbHandler,
		saver:       saver,
	}

}

func (d *Decoder) Start(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {
		defer wg.Done()

		killed := false

		for {

			if killed {
				break
			}

			d.WorkerPool <- d.WorkerQueue

			select {
			case conn := <-d.WorkerQueue:

				buffer := make([]byte, 1024)

				n, err := conn.Read(buffer)
				if err != nil {
					log.Println("Error reading from connection: ", err)
					return
				}
				buffer = buffer[:n]

				stationID, token, moment, randomNum, err := d.decode(buffer)

				if err != nil {
					log.Println("Error decoding buffer: ", err)
					// Handle the error appropriately, maybe send an error response and return
					answer(conn, 1, 1, token) // Example response
					return
				}

				valid, newToken := d.saver.Validate(stationID, token)
				if !valid {
					log.Println("Invalid token for station ID: ", stationID)
					answer(conn, 1, 2, newToken) // Example response
					return
				}

				station, err := d.dbHandler.ReadStation(stationID)
				if err != nil {
					log.Println("Error reading station from DB: ", err)
					answer(conn, 1, 3, newToken) // Example response
					return
				}

				measurement, err := d.dbHandler.SendMeasurement(randomNum)
				if err != nil {
					log.Println("Error sending measurement to DB: ", err)
					answer(conn, 1, 4, newToken) // Example response
					return
				}

				err = d.dbHandler.SendEntry(measurement, station, moment)
				if err != nil {
					log.Println("Error sending entry to DB: ", err)
					answer(conn, 1, 5, newToken) // Example response
					return
				}

				answer(conn, 1, 0, newToken)

				d.dbHandler.SendEntry(measurement, station, moment)

			case killed = <-d.killChan:
				close(d.WorkerQueue)
				close(d.killChan)

			}

		}

	}()

}

func (d *Decoder) Kill() {
	d.killChan <- true
}

func answer(conn net.Conn, version byte, status byte, token string) {

	defer conn.Close()

	tokenB := []byte(token)
	response := append([]byte{version, status}, tokenB...)

	_, err := conn.Write(response)
	if err != nil {
		log.Println("Error writing to connection: ", err)
	}
}

func (d *Decoder) decode(buffer []byte) (
	stationID int,
	token string,
	moment time.Time,
	randomNum int,
	err error) {

	if version != buffer[0] {
		return -1, "", time.Now(), 0, fmt.Errorf("wrong version")
	}

	stationID = int(buffer[1])<<24 |
		int(buffer[2])<<16 |
		int(buffer[3])<<8 |
		int(buffer[4])

	token = string(buffer[5:11])

	moment = time.Unix(int64(binary.BigEndian.Uint32(buffer[11:15])), 0)

	randomNum = int(buffer[15])<<24 |
		int(buffer[16])<<16 |
		int(buffer[17])<<8 |
		int(buffer[18])

	return

}
