package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"trucker/handlers"

	"github.com/streadway/amqp"
)

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
func publishToRabbitMQ(location handlers.Location_report) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare an exchange
	err = ch.ExchangeDeclare(
		"location_exchange", // name
		"fanout",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// Publish the location message
	body := fmt.Sprintf("Device: %s, Lat: %f, Lon: %f, Alt: %f",
		location.Device_id, location.Latitude, location.Longitude, location.Altitude)
	err = ch.Publish(
		"location_exchange", // exchange
		"",                  // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Location published: %s", body)
}
func handleConnection(conn net.Conn) {
	log.SetPrefix("handleConnection(): ")
	defer conn.Close()

	buffer := make([]byte, 1024)

	//Channels
	raw_message := make(chan []byte, 1024)
	protocol_message := make(chan handlers.Procotol_Message, 1024)
	aproved_messages := make(chan handlers.Procotol_Message, 1024)
	response_gps := make(chan handlers.Response_GPS, 1024)
	location_reports := make(chan handlers.Location_report, 1024)
	leavings := make(chan interface{}, 1024)

	go handlers.IdentifyProtocol(raw_message, protocol_message)

	go handlers.ValidateMessage(protocol_message, aproved_messages)

	go handlers.Process(aproved_messages, response_gps, location_reports, leavings)

	//Response manager
	go func() {
		for {
			select {
			case response, ok := <-response_gps:
				if !ok {
					fmt.Println("response_gps channel closed")
					return
				}
				conn.Write(response.Message)
			case location, ok := <-location_reports:
				if !ok {
					fmt.Println("response_gps channel closed")
					return
				}
				log.Println(location)
				publishToRabbitMQ(location)
			case <-leavings:
				fmt.Println("Closing the connection")
				conn.Close()
				return
			}
		}
	}()

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			log.Printf("There is an error while receiving data: %s\n", err)
			break
		}
		fmt.Printf("Receiving this message from GPS %x\n", buffer[:n])
		raw_message <- buffer[:n]
	}
	fmt.Println("Closing the chanels")
	close(response_gps)
	close(location_reports)
	close(raw_message)
	close(protocol_message)
	close(aproved_messages)
	close(leavings)
}
