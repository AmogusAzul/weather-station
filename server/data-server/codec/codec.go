package codec

import (
	"log"
	"net"
	"sync"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/dispatcher"
	safety "github.com/AmogusAzul/weather-station/server/data-server/safety"
)

const (
	version         byte = 1
	maxBufferLength uint = 1024

	stationType byte = 1
)

var RequestHandlers map[byte]RequestHandler = map[byte]RequestHandler{
	stationType: StationHandler,
}

type Decoder struct {
	WorkerQueue chan net.Conn
	WorkerPool  chan<- chan net.Conn
	killChan    chan bool

	dbHandler *dbhandle.DbHandler
	saver     *safety.Saver
}

func GetDecoder(
	workerPool chan chan net.Conn,
	dbHandler *dbhandle.DbHandler,
	saver *safety.Saver,
) dispatcher.Worker {

	return &Decoder{
		WorkerQueue: make(chan net.Conn),
		WorkerPool:  workerPool,
		killChan:    make(chan bool),
		dbHandler:   dbHandler,
		saver:       saver,
	}

}
func (d *Decoder) Kill() {
	go func() { d.killChan <- true }()
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

			// Letting the dispatcher send a new job to
			d.WorkerPool <- d.WorkerQueue

			select {
			case conn := <-d.WorkerQueue:

				buffer := make([]byte, maxBufferLength)

				n, err := conn.Read(buffer)
				if err != nil {
					log.Println("Error reading from connection: ", err)
					continue
				}
				buffer = buffer[:n]

				var requestError byte = 0
				requestVersion, requestType, content := d.decodeRequest(buffer)
				requestHandler := RequestHandlers[requestType]

				if requestHandler == nil {
					requestError = typeError
				}
				if requestVersion != version {
					requestError = versionError
				}
				if requestError != 0 {
					ErrorAnswer(conn, requestError)
					continue
				}

				err = requestHandler(conn, content, d.saver, d.dbHandler)
				if err != nil {
					log.Panicf("error (%s) while processing %v with type %d", err, content, requestType)
				}

			case killed = <-d.killChan:
				close(d.WorkerQueue)
				close(d.killChan)
			}
		}
	}()
}

func (d *Decoder) decodeRequest(buffer []byte) (
	requestVersion byte,
	requestType byte,
	requestContent []byte,
) {
	return buffer[0], buffer[1], buffer[2:]
}

func UncompatibleErrorAnswer(conn net.Conn, requestVersion byte, requestError byte) {

	defer conn.Close()

	_, err := conn.Write([]byte{requestVersion, requestError})
	if err != nil {
		log.Println("Error responding to connection: ", err)
	}
}
