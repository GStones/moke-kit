package network

import (
	"testing"
)

func newPort() Port {
	return Port(0)
}

func TestPort_ListenAddress(t *testing.T) {
	port := newPort()
	port.ListenAddress()
}

func TestPort_String(t *testing.T) {
	port := newPort()
	_ = port.String()
}

func TestPort_UnmarshalText(t *testing.T) {
	textSlice := []byte("1000")
	port := newPort()
	err := port.UnmarshalText(textSlice)
	if err != nil {
		t.Error(err)
		return
	}
}
