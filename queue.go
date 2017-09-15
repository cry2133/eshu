package eshu

type QueueClient interface {
	Send(data []byte) error
	Close() error
}
