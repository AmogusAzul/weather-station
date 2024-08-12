package dispatcher

import "sync"

type Dispatcher struct {
	workerPool   chan chan interface{}
	createWorker func(chan chan interface{}) *Worker
	maxWorkers   int
	JobQueue     <-chan interface{}
	killChan     chan bool
}

func NewDispatcher(
	maxWorkers int,
	createWorker func(chan chan interface{}) *Worker,
	jobQueue chan interface{},
) *Dispatcher {
	return &Dispatcher{

		maxWorkers:   maxWorkers,
		workerPool:   make(chan chan interface{}, maxWorkers),
		createWorker: createWorker,

		JobQueue: jobQueue,

		killChan: make(chan bool),
	}
}

func (d *Dispatcher) Dispatch(wg *sync.WaitGroup) {

	wg.Add(1)

	workers := make([]*Worker, d.maxWorkers)

	for i := 0; i < d.maxWorkers; i++ {

		workers = append(workers, d.createWorker(d.workerPool))
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

			case killed = <-d.killChan:

				close(d.killChan)

				close(d.workerPool)
				for _, worker := range workers {
					(*worker).Kill()
				}
			}
		}

	}()

}

func (d *Dispatcher) Kill() {
	d.killChan <- true
}

type Worker interface {
	Start()
	Kill()
}
