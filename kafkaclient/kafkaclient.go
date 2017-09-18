package kafkaclient

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var DefaultKafkaConfig *sarama.Config

func init() {
	DefaultKafkaConfig = sarama.NewConfig()
	c.Net.DialTimeout = 5 * time.Second
	c.Net.ReadTimeout = 6 * time.Second
	c.Net.WriteTimeout = 6 * time.Second
	c.Producer.Return.Successes = true
}

type KafkaClient struct {
	prod  sarama.SyncProducer
	addrs []string
	topic string
}

func NewKafkaClient(kafkaAddrs []string, topic string, c *sarama.Config) (*KafkaClient, error) {
	if len(kafkaAddrs) == 0 {
		return nil, fmt.Errorf("empty list of kafka addresses")
	}
	if topic == "" {
		return nil, fmt.Errorf("empty kafka topic")
	}

	prod, err := sarama.NewSyncProducer(kafkaAddrs, c)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{prod: prod, addrs: kafkaAddrs, topic: topic}, nil
}

func (kc *KafkaClient) Send(data []byte) error {
	msg := sarama.ProducerMessage{Topic: kc.topic, Value: sarama.ByteEncoder(data)}

	if _, _, err := kc.prod.SendMessage(&msg); err != nil {
		return fmt.Errorf("Failed to send Kafka message, host: %v, topic: %q, error: %v", kc.addrs, kc.topic, err)
	}

	return nil
}

func (kc *KafkaClient) Close() error {
	return kc.prod.Close()
}
