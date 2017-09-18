package eshu

// QueueClient represents a queue-system client
type QueueClient interface {
	Send(data []byte) error
	Close() error
}
