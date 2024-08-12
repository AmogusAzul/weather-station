package tcp

import (
	"log"
	"net"
	"sync"
)

type StationListener struct {
	listener net.Listener

	JobQueue chan net.Conn

	killChan chan bool
}

func GetStationLister(port string, jobQueue chan net.Conn) *StationListener {

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Couldn't start listener error: %v", err)
	}

	return &StationListener{
		listener: listener,
		JobQueue: jobQueue,
	}
}

func (sl *StationListener) Listen(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {

		for {

			conn, err := sl.listener.Accept()
			if err != nil {
				break

			}
			sl.JobQueue <- conn
		}
	}()

	go func() {

		<-sl.KillChan
		sl.listener.Close()

		wg.Done()

	}()
}

func (sl *StationListener) Close() {
	sl.KillChan <- true
}
