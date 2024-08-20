package dispatcher

import (
	"sync"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

type JobDispatcher[T any] struct {
	workerPool   chan chan T
	createWorker func(chan chan T, *dbhandle.DbHandler, *safety.Saver) Worker
	maxWorkers   int
	JobQueue     <-chan T
	killChan     chan bool
}

func NewJobDispatcher[T any](
	maxWorkers int,
	createWorker func(chan chan T, *dbhandle.DbHandler, *safety.Saver) Worker,
	jobQueue chan T,
) *JobDispatcher[T] {
	return &JobDispatcher[T]{

		maxWorkers:   maxWorkers,
		workerPool:   make(chan chan T, maxWorkers),
		createWorker: createWorker,

		JobQueue: jobQueue,

		killChan: make(chan bool),
	}
}

func (d *JobDispatcher[T]) Dispatch(wg *sync.WaitGroup, dbHandler *dbhandle.DbHandler, saver *safety.Saver) {

	wg.Add(1)

	workers := make([]Worker, d.maxWorkers)

	for i := 0; i < d.maxWorkers; i++ {
		workers[i] = d.createWorker(d.workerPool, dbHandler, saver)
		workers[i].Start(wg)
	}

	go func() {

		defer wg.Done()
		killed := false

		for {

			if killed {
				break
			}

			select {

			case job := <-d.JobQueue:
				// job enters in the oldest free worker
				<-d.workerPool <- job

			case killed = <-d.killChan:

				close(d.killChan)

				close(d.workerPool)
				for _, worker := range workers {
					(worker).Kill()
				}
			}
		}

	}()

}

func (d *JobDispatcher[T]) Kill() {
	d.killChan <- true
}

type Worker interface {
	Start(*sync.WaitGroup)
	Kill()
}
