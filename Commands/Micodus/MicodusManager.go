package micodus

import (
	"fmt"
	"trucker/commands"
)

type MicodusManager struct {
	commands map[string]commands.Command
}

func NewMicodusManager() *MicodusManager {
	return &MicodusManager{
		commands: map[string]commands.Command{
			"0003": &Close{},
			"0100": &Register{},
			"0102": &Authentication{},
			"0201": &ReceiveLocation{},
		},
	}
}

func (t *MicodusManager) ExecuteCommand(order string, message []byte) (commands.CommandResponse, error) {
	if cmd, exists := t.commands[order]; exists {
		return cmd.Execute(message)
	}
	return nil, fmt.Errorf("command not found: %s", order)
}
func (t *MicodusManager) IdentifyCommand(message []byte) string {
	message_type := fmt.Sprintf("%02x%02x", message[1], message[2])
	return message_type
}
