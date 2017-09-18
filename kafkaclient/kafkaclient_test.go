package kafkaclient

import (
	"testing"
)

func TestNewKafkaClient(t *testing.T) {
	_, err := NewKafkaClient([]string{"localhost"}, "abc", nil)
	if err != nil {
		t.Fatal(err)
	}
}
