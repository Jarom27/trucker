package connection

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"trucker/protocols"
)

type TCPServer struct {
	Address string
	Service *protocols.ProtocolService
	done    chan bool
}

func NewTCPServer(address string, service *protocols.ProtocolService) *TCPServer {
	return &TCPServer{
		Address: address,
		Service: service,
		done:    make(chan bool),
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine para escuchar la señal de interrupción
	go func() {
		<-signalChan
		fmt.Println("\nGraceful shutdown initiated...")
		close(s.done) // Cerrar el canal para indicar el final
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				// Si el canal está cerrado, salimos del bucle
				fmt.Println("Server is shutting down, stopped accepting connections.")
				wg.Wait()
				return nil
			default:
				fmt.Println("Failed to accept connection:", err)
				continue
			}
		}

		wg.Add(1)
		go s.handleConnection(conn, &wg)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in handleConnection:", r)
		}
	}()

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
		} else {
			break
		}
	}
	fmt.Println("Connection was successful close")

}
