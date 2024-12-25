package connection

import (
	"fmt"
	"net"
	"sync"
	protocols "trucker/Protocols"
)

type TCPServer struct {
	Address string
	Service *protocols.ProtocolService
}

func NewTCPServer(address string, service *protocols.ProtocolService) *TCPServer {
	return &TCPServer{
		Address: address,
		Service: service,
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	defer listener.Close()

	fmt.Println("TCP Server started on", s.Address)

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		wg.Add(1)
		go s.handleConnection(conn, &wg)
	}

	wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	fmt.Println("New connection from:", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed by client:", conn.RemoteAddr())
			return
		}
		message := buffer[:n]
		fmt.Printf("Message received: %x\n", message)
		response, err := s.Service.ProcessCommand(message)

		if err != nil {
			fmt.Println("There was an error while processing a command")
		}

		if response != nil {
			fmt.Printf("%x", response)
			conn.Write(response)
		}
	}
}
