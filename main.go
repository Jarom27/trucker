package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"trucker/handlers"
)

type DeviceCommand struct {
	protocol string
	message  []byte
}

func main() {
	log.SetPrefix("main(): ")
	file, err := os.OpenFile("logs/errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(file)

	if err != nil {
		log.Fatal("it was not possible to load the log system ", err)
		return
	}

	listener, err := net.Listen("tcp", ":7700")
	if err != nil {
		log.Fatal("it was not possible to start the server, error: ", err)
		return
	}

	defer listener.Close()
	defer file.Close()

	fmt.Println("Server listening on: ", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("It is not possible to connect with the client %s\n", err)
			continue
		}
		fmt.Println("Connection was sucessfull established with: ", conn.RemoteAddr())
		go handleConnection(conn)

	}

}
func handleConnection(conn net.Conn) {
	log.SetPrefix("handleConnection(): ")
	defer conn.Close()

	buffer := make([]byte, 1024)
	raw_message := make(chan []byte, 1024)
	protocol_message := make(chan handlers.Procotol_Message, 1024)
	aproved_messages := make(chan handlers.Procotol_Message, 1024)
	response_gps := make(chan handlers.Response_GPS, 1024)
	location_reports := make(chan handlers.Location_report, 1024)
	leavings := make(chan interface{}, 1024)
	//Identify protocol
	go handlers.IdentifyProtocol(raw_message, protocol_message)

	//Validate Checksu
	go handlers.ValidateMessage(protocol_message, aproved_messages)

	go handlers.Process(protocol_message, response_gps, location_reports, leavings)

	go func() {
		select {
		case response := <-response_gps:
			conn.Write(response.Message)
		case location := <-location_reports:
			fmt.Println(location)
		case <-leavings:
			conn.Close()
		}
	}()

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			log.Printf("There is an error while receiving data: %s\n", err)
			break
		}
		log.Printf("Receiveng this message from GPS %x\n", buffer[:n])
		raw_message <- buffer[:n]
	}
	close(response_gps)
	close(location_reports)
	close(raw_message)
	close(protocol_message)
	close(aproved_messages)
	close(leavings)
}
