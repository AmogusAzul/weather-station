package main

import (
	"fmt"
	"net"
	"os"
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

	//TEMP
	maxWorkers := 2
	jobBuffer := 8
	tokenLength := 6

	jobQueue := make(chan net.Conn, jobBuffer)

	sl := tcp.GetStationLister("8080", jobQueue)

	jobDispatcher := dispatcher.NewJobDispatcher[net.Conn](maxWorkers, codec.GetDecoder, jobQueue)

	dbHandler, err := dbhandle.GetDbHandler(dbUser, dbPassword, dbHost, dbPort, dbName)

	saver := safety.GetSaver(tokenLength, tokenPath)

	if err != nil {
		fmt.Println("dbHandler fucked")
		return
	}

	e := executer.GetExecuter(saver, dbHandler, &wg)
	e.AddToClose(sl, jobDispatcher, dbHandler, saver)
	jobDispatcher.Dispatch(&wg, e)
	sl.Listen(&wg)

}
