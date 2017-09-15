package eshu

import "github.com/hyperledger/fabric/events/consumer"

type Connector struct {
	hlfClient *consumer.EventsClient
	qClient   QueueClient
}

func NewConnector() *Connector {
	return nil
}

func (c *Connector) StartListening() {
}

func (c *Connector) Stop() {
}
