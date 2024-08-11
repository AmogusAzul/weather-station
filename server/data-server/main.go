package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting TCP server: ", err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn) // Handle each connection concurrently
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to store incoming data
	buffer := make([]byte, 1024)

	// Read data from the connection
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}

	// Process the received data
	receivedData := string(buffer[:n])
	fmt.Println("Received data: ", receivedData)

	// Send a response back to the client
	response := "Message received: " + receivedData
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Println("Error writing to connection: ", err)
		return
	}
}
