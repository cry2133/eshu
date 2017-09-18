package eshu

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/events/consumer"
	"github.com/hyperledger/fabric/protos/peer"
)

type eventHandler struct {
	interests []*peer.Interest
	qClient   QueueClient
}

func newEventHandler(events map[string][]string, qClient QueueClient) *eventHandler {
	e := eventHandler{
		qClient: qClient,
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
	data, err := proto.Marshal(msg)
	if err != nil {
		return false, err
	}

	if err := e.qClient.Send(data); err != nil {
		return false, err
	}

	return true, nil
}

func (e *eventHandler) Disconnected(err error) {
	log.Printf("Disconnected error: '%v'", err)
}

// Connector represents a link between event listener and a queue-system client
type Connector struct {
	hlfClient *consumer.EventsClient
	qClient   QueueClient
}

// NewConnector returns a pointer to Connector otherwise returns nil and the reason as an error
func NewConnector(peerAddress string, regTimeout time.Duration, qClient QueueClient, events map[string][]string) (*Connector, error) {
	if qClient == nil {
		return nil, errors.New("nil queue-system client")
	}

	hlfc, err := consumer.NewEventsClient(peerAddress, regTimeout, newEventHandler(events, qClient))
	if err != nil {
		return nil, err
	}

	return &Connector{hlfClient: hlfc, qClient: qClient}, nil
}

// StartListening starts listening events from Hyperledger Fabric and it'll start to send them through the message system
func (c *Connector) StartListening() error {
	return c.hlfClient.Start()
}

// Stop stops the listener and closes the queue-system client
func (c *Connector) Stop() (err error) {
	if errHLF := c.hlfClient.Stop(); errHLF != nil {
		err = errHLF
	}

	if errQ := c.qClient.Close(); errQ != nil {
		if err == nil {
			err = errQ
		} else {
			err = fmt.Errorf("%v - %v", err, errQ)
		}
	}

	return
}
