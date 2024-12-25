package protocols

import (
	"bytes"
	"fmt"
	commands "trucker/Commands"
	micodus "trucker/Commands/Micodus"
)

type Protocol = commands.Protocol

type ProtocolStrategy struct {
	protocols map[string]func() Protocol
}

func NewProtocolStrategy() *ProtocolStrategy {
	return &ProtocolStrategy{
		protocols: map[string]func() Protocol{
			"Micodus": func() Protocol { return micodus.NewMicodusManager() },
		},
	}
}

func (ps *ProtocolStrategy) GetProtocol(message []byte) (Protocol, error) {
	switch {
	case bytes.HasPrefix(message, []byte{0x7e}) && bytes.HasSuffix(message, []byte{0x7e}):
		return ps.protocols["Micodus"](), nil
	default:
		return nil, fmt.Errorf("protocol not found for message: %x", message)
	}
}
