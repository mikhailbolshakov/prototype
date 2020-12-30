package queue

type Queue interface {

	Open(clientId string) error
	Close() error
	Publish(topic string, msg []byte) error
	Subscribe(topic string, receiverChan chan<- []byte) error

}
