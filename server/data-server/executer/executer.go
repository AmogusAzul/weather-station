package executer

import (
	"fmt"
	"sync"

	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

type Executer struct {
	Saver      *safety.Saver
	DBHandler  *dbhandle.DbHandler
	CloseFuncs []CloseFunc
	wg         *sync.WaitGroup
}

func GetExecuter(saver *safety.Saver, dbHandler *dbhandle.DbHandler, wg *sync.WaitGroup) *Executer {

	return &Executer{
		Saver:     saver,
		DBHandler: dbHandler,
		wg:        wg,
	}
}

func (e *Executer) AddToClose(closables ...Closable) {
	for _, c := range closables {
		e.CloseFuncs = append(e.CloseFuncs, c.Close)
	}
}

func (e *Executer) CloseAll() (err error) {
	defer e.wg.Done()
	e.wg.Add(1)
	for _, closeFunc := range e.CloseFuncs {

		if cerr := closeFunc(); cerr != nil {
			fmt.Println(cerr)
			err = cerr // captures the last error encountered
		}

	}
	return
}

type CloseFunc func() error

type Closable interface {
	Close() error
}
