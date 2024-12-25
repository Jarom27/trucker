package commands

type Protocol interface {
	IdentifyCommand(message []byte) string
	ExecuteCommand(order string, message []byte) ([]byte, error)
}
