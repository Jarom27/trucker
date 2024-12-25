package commands

type CommandResponse interface {
	ToMap() (map[string]interface{}, error)
	ToJSON() ([]byte, error)
}
