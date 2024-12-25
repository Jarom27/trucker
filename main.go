package main

import (
	"fmt"
	"os"
	connection "trucker/Connection"
	protocols "trucker/Protocols"
)

func main() {
	host := os.Getenv("TRUCKER_HOST")
	port := os.Getenv("TRUCKER_PORT")

	if port == "" {
		port = "7700"
	}

	strategy := protocols.NewProtocolStrategy()
	service := protocols.NewProtocolService(strategy)

	server := connection.NewTCPServer(fmt.Sprintf("%s:%s", host, port), service)
	server.Start()
}
