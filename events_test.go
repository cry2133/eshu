package eshu

import (
	"testing"
	"time"
)

type mockQueueClient struct{}

func (*mockQueueClient) Send(d []byte) error {
	return nil
}

func (*mockQueueClient) Close() error {
	return nil
}

func TestCreateConnectorOK(t *testing.T) {
	ech := map[string][]string{
		"ch-01": []string{"event1"},
	}
	if _, err := NewConnector("test.address", time.Second, new(mockQueueClient), ech); err != nil {
		t.Fatalf("Expected nil error, got error: %v", err)
	}
}

func TestCreateConnectorWithError(t *testing.T) {
	if _, err := NewConnector("test.address", time.Hour, nil, nil); err == nil {
		t.Fatal("Expected nil connector")
	}
}
