package internal

import (
	"math/rand"
	"time"

	"github.com/gstones/platform/services/common/network"
)

type Brokers []network.Address

func (b Brokers) Random() network.Address {
	return b[random.Intn(len(b))]
}

func (b Brokers) Strings() []string {
	var as []string
	for _, a := range b {
		as = append(as, a.String())
	}
	return as
}

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().Unix()))
}
