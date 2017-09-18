package eshu

import (
	"errors"
	"time"

	"github.com/hyperledger/fabric/events/consumer"
	"github.com/hyperledger/fabric/protos/peer"
)

type eventHandler struct {
	interests []*peer.Interest
}

func newEventHandler(events map[string][]string) *eventHandler {
	e := eventHandler{
		interests: []*peer.Interest{
			{EventType: peer.EventType_REGISTER},
			{EventType: peer.EventType_BLOCK},
			{EventType: peer.EventType_REJECTION},
		},
	}

	for chID, eventNames := range events {
		for _, name := range eventNames {
			e.interests = append(e.interests, &peer.Interest{
				EventType: peer.EventType_CHAINCODE,
				RegInfo: &peer.Interest_ChaincodeRegInfo{
					ChaincodeRegInfo: &peer.ChaincodeReg{ChaincodeId: chID, EventName: name},
				},
			})
		}
	}

	return &e
}

func (e *eventHandler) GetInterestedEvents() ([]*peer.Interest, error) {
	return e.interests, nil
}

func (e *eventHandler) Recv(msg *peer.Event) (bool, error) {
	return true, nil
}

func (e *eventHandler) Disconnected(err error) {}

// Connector represents a link between event listener and a queue-system client
type Connector struct {
	hlfClient *consumer.EventsClient
	qClient   QueueClient
}

// NewConnector returns a pointer to Connector otherwise returns nil and the reason as an error
func NewConnector(peerAddress string, regTimeout time.Duration, qClient QueueClient, events map[string][]string) (*Connector, error) {
	if qClient == nil {
		return nil, errors.New("empty queue system client")
	}

	hlfc, err := consumer.NewEventsClient(peerAddress, regTimeout, newEventHandler(events))
	if err != nil {
		return nil, err
	}

	return &Connector{qClient: qClient, hlfClient: hlfc}, nil
}

// StartListening starts listening events from Hyperledger Fabric and it'll start to send them through the message system
func (c *Connector) StartListening() {
}

// Stop stops the listener
func (c *Connector) Stop() {
}
