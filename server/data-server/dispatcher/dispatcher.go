package dispatcher

import "sync"

type Dispatcher struct {
	workerPool   chan chan interface{}
	createWorker func(int, chan chan interface{}) *Worker
	maxWorkers   int
	JobQueue     <-chan interface{}
	KillChan     chan bool
}

func NewDispatcher(
	maxWorkers int,
	createWorker func(int, chan chan interface{}) *Worker,
	jobQueue chan interface{},
) *Dispatcher {
	return &Dispatcher{

		maxWorkers:   maxWorkers,
		workerPool:   make(chan chan interface{}, maxWorkers),
		createWorker: createWorker,

		JobQueue: jobQueue,

		KillChan: make(chan bool),
	}
}

func (d *Dispatcher) Dispatch(wg *sync.WaitGroup) {

	wg.Add(1)

	workers := make([]*Worker, d.maxWorkers)

	for i := 0; i < d.maxWorkers; i++ {

		workers = append(workers, d.createWorker(i, d.workerPool))
		(*workers[i]).Start()

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

			case killed = <-d.KillChan:

				close(d.KillChan)

				close(d.workerPool)
				for _, worker := range workers {
					(*worker).Kill()
				}
			}
		}

	}()

}

type Worker interface {
	Start()
	Kill()
}
