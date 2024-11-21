package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/AmogusAzul/weather-station/server/data-server/codec"
	dbhandle "github.com/AmogusAzul/weather-station/server/data-server/db-handle"
	"github.com/AmogusAzul/weather-station/server/data-server/dispatcher"
	"github.com/AmogusAzul/weather-station/server/data-server/executer"
	"github.com/AmogusAzul/weather-station/server/data-server/safety"
	"github.com/AmogusAzul/weather-station/server/data-server/tcp"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var wg sync.WaitGroup
	defer wg.Wait()

	//waiting for db to start
	time.Sleep(10 * time.Second)

	//db env variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	//token dir
	tokenPath := os.Getenv("TOKEN_PATH")

	// Concurrency parameters
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	jobBuffer, err0 := strconv.Atoi(os.Getenv("JOB_BUFFER"))

	if err != nil || err0 != nil {
		fmt.Println("MAX_WORKERS or JOB_BUFFER isn't defined correctly in .env, they both have to be numbers")
		return
	}

	jobQueue := make(chan net.Conn, jobBuffer)

	sl := tcp.GetStationLister("8080", jobQueue)

	jobDispatcher := dispatcher.NewJobDispatcher[net.Conn](maxWorkers, codec.GetDecoder, jobQueue)

	dbHandler, err := dbhandle.GetDbHandler(dbUser, dbPassword, dbHost, dbPort, dbName)

	saver := safety.GetSaver(tokenPath)

	if err != nil {
		fmt.Println("dbHandler fucked")
		return
	}

	e := executer.GetExecuter(saver, dbHandler, &wg)
	e.AddToClose(sl, jobDispatcher, dbHandler, saver)
	jobDispatcher.Dispatch(&wg, e)
	sl.Listen(&wg)

}
