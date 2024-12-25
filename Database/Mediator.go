package database

type Mediator interface {
	Publish(message []byte)
	Consume()
	Connect()
	Close()
}
