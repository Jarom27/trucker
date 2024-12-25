package commands

type CommandResponse interface {
	ToMap() map[string]interface{}
	ToJSON() ([]byte, error)
}
