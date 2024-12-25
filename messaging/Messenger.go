package messaging

type Messenger interface {
	Send(interface{}) error
}
