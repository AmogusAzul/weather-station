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

func (sl *StationListener) Listen(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {
		defer wg.Done() // Ensure the WaitGroup is marked as done when the goroutine exits

		for {
			select {
			case <-sl.killChan: // If a shutdown signal is received
				sl.listener.Close()
				close(sl.JobQueue) // Close the job queue to signal workers to finish
				return             // Exit the goroutine

			default:
				conn, err := sl.listener.Accept()
				if err != nil {
					select {
					case <-sl.killChan: // Handle the case where Accept fails due to listener closing
						return
					default:
						log.Printf("Error accepting connection: %v", err)
					}
					continue
				}

				select {
				case sl.JobQueue <- conn: // Send the connection to the job queue
				case <-sl.killChan: // Handle shutdown while waiting to send to the job queue
					conn.Close() // Close the connection if we're shutting down
					return
				}
			}
		}
	}()
}

func (sl *StationListener) Close() error {
	sl.killChan <- true
	return nil
}
