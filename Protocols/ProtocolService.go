package protocols

import (
	"fmt"
	"sync"
	micodus "trucker/Commands/Micodus"
)

// Servicio con concurrencia
type ProtocolService struct {
	strategy *ProtocolStrategy
	mu       sync.Mutex // Para proteger recursos compartidos (si los hay)
}

// Constructor
func NewProtocolService(strategy *ProtocolStrategy) *ProtocolService {
	return &ProtocolService{
		strategy: strategy,
	}
}

func (s *ProtocolService) ProcessCommand(message []byte) ([]byte, error) {
	protocol, err := s.strategy.GetProtocol(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol: %w", err)
	}
	fmt.Printf("Protocol was selected\n")
	command := protocol.(*micodus.MicodusManager).IdentifyCommand(message)
	fmt.Printf("Command was selected: %s\n", command)
	response, err := protocol.(*micodus.MicodusManager).ExecuteCommand(command, message)

	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	fmt.Printf("Command Response: %x\n", response)
	return response, nil
}
