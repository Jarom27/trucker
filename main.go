package main

import (
	"fmt"
	"os"
	"trucker/connection"
	"trucker/messaging"
	"trucker/protocols"
)

func main() {
	host := os.Getenv("TRUCKER_HOST")
	port := os.Getenv("TRUCKER_PORT")

	if port == "" {
		port = "7700"
	}
	queue_host := os.Getenv("RABBIT_HOST")
	queue_port := os.Getenv("RABBIT_PORT")
	queue_name := os.Getenv("RABBIT_QUEUE_NAME")
	queue_user := os.Getenv("RABBIT_DEFAULT_USER")
	queue_pass := os.Getenv("RABBIT_DEFAULT_PASS")

	err, sender := messaging.NewQueueMessenger(queue_name, queue_user, queue_pass, queue_host, queue_port)

	if err != nil {
		fmt.Println("There was an error during the registration of queue")
	}
	strategy := protocols.NewProtocolStrategy()
	service := protocols.NewProtocolService(strategy, sender)

	server := connection.NewTCPServer(fmt.Sprintf("%s:%s", host, port), service)
	server.Start()
}
