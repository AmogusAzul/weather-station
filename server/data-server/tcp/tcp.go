package tcp

import (
	"log"
	"net"
	"sync"
)

type StationListener struct {
	listener net.Listener

	JobQueue chan<- net.Conn

	killChan chan bool
}

func GetStationLister(port string, jobQueue chan<- net.Conn) *StationListener {

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Couldn't start listener error: %v", err)
	}

	return &StationListener{
		listener: listener,
		JobQueue: jobQueue,
		killChan: make(chan bool),
	}
}

// Listen starts the listener in a goroutine and waits for incoming connections.
func (sl *StationListener) Listen(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done() // Mark the WaitGroup as done when this goroutine exits

		for {
			conn, err := sl.listener.Accept()

			// Check if a shutdown signal was received
			select {
			case <-sl.killChan:
				return

			default:
				if err != nil {
					log.Printf("Error accepting connection: %v", err)
					continue
				}

				select {
				case sl.JobQueue <- conn:
					// Connection successfully sent to JobQueue
				case <-sl.killChan:
					conn.Close() // Close the connection if we're shutting down
					return
				}
			}
		}
	}()
}

// Close sends a shutdown signal to the listener.
func (sl *StationListener) Close() error {
	close(sl.killChan)
	return sl.listener.Close() // Ensure the listener is closed when shutting down
}
