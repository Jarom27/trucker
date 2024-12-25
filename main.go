package main

import (
	"fmt"
	"trucker/config"
	"trucker/connection"
	"trucker/messaging"
	"trucker/protocols"
)

func main() {
	config := config.LoadConfig()

	sender, err := messaging.NewQueueMessenger(
		config.RabbitQueue,
		config.RabbitUser,
		config.RabbitPass,
		config.RabbitHost,
		config.RabbitPort,
	)

	if err != nil {
		fmt.Printf("There was an error during the registration of queue: %s", err)
	}
	strategy := protocols.NewProtocolStrategy()
	service := protocols.NewProtocolService(strategy, sender)

	server := connection.NewTCPServer(fmt.Sprintf("%s:%s", config.Host, config.Port), service)
	server.Start()
}
