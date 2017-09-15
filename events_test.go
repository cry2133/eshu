package eshu

import "testing"

func TestNilConnector(t *testing.T) {
	if c := NewConnector(); c != nil {
		t.Fatal("Expected nil connector")
	}

	t.Log("TODO: more tests :)")
}
