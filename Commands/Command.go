package commands

type Command interface {
	Execute(message []byte) ([]byte, error)
}
