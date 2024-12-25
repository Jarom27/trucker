package micodus

type Close struct {
}

func (t *Close) Execute(message []byte) ([]byte, error) {
	return nil, nil
}
