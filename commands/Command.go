package commands

type Command interface {
	Execute(message []byte) (CommandResponse, error)
}
