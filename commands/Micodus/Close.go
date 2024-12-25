package micodus

import "trucker/commands"

type Close struct {
}

func (t *Close) Execute(message []byte) (commands.CommandResponse, error) {
	return nil, nil
}
