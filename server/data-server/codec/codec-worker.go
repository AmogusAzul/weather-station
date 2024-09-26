package codec

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/AmogusAzul/weather-station/server/data-server/dispatcher"
	"github.com/AmogusAzul/weather-station/server/data-server/executer"
)

const (
	version         byte = 1
	maxBufferLength uint = 1024

	stationType    byte = 1
	newStationType byte = 2
	closeType      byte = 3
)

var RequestHandlers map[byte]RequestHandler = map[byte]RequestHandler{
	stationType:    StationHandler,
	newStationType: NewStationHandler,
	closeType:      CloseHandler,
}

type Decoder struct {
	WorkerQueue chan net.Conn
	WorkerPool  chan<- chan net.Conn
	killChan    chan bool

	Executer *executer.Executer
}

func GetDecoder(
	workerPool chan chan net.Conn,
	e *executer.Executer,
) dispatcher.Worker {

	return &Decoder{
		WorkerQueue: make(chan net.Conn),
		WorkerPool:  workerPool,
		killChan:    make(chan bool),
		Executer:    e,
	}
}
func (d *Decoder) Close() {
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

				fmt.Println(buffer)

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

				err = requestHandler(conn, content, d.Executer)
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
