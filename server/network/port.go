package network

import (
	"fmt"
	"strconv"
)

type Port int

func (p Port) String() string {
	return fmt.Sprintf("%d", p)
}

func (p Port) Value() int {
	return int(p)
}

func (p Port) ListenAddress() string {
	return fmt.Sprintf(":%d", p)
}

func (p *Port) UnmarshalText(text []byte) error {
	if port, err := strconv.ParseInt(string(text), 10, 64); err != nil {
		return err
	} else {
		*p = Port(int(port))
		return nil
	}
}
