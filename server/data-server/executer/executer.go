package executer

import (
	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
)

type Executer struct {
	Saver      *safety.Saver
	DBHandler  *dbhandle.DbHandler
	CloseFuncs []CloseFunc
}

func GetExecuter(saver *safety.Saver, dbHandler *dbhandle.DbHandler) *Executer {

	return &Executer{
		Saver:     saver,
		DBHandler: dbHandler,
	}
}

func (e *Executer) AddToClose(closables ...Closable) {
	for _, c := range closables {
		e.CloseFuncs = append(e.CloseFuncs, c.Close)
	}
}

func (e *Executer) CloseAll() (err error) {
	for _, closeFunc := range e.CloseFuncs {
		if cerr := closeFunc(); cerr != nil {
			err = cerr // captures the last error encountered
		}
	}
	return
}

type CloseFunc func() error

type Closable interface {
	Close() error
}
