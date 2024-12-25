package messaging

type Messenger interface {
	Send(data []byte) error
}
