package network

import "net"

type Network string

const (
	TcpNetwork Network = "tcp"
)

func (n Network) String() string {
	return (string)(n)
}

func NewTcpListener(port Port) (net.Listener, error) {
	return net.Listen(TcpNetwork.String(), port.ListenAddress())
}
