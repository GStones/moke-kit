package internal

import (
	k "github.com/segmentio/kafka-go"
)

type Balancer = k.Balancer

type RoundRobin = k.RoundRobin
type LeastBytes = k.LeastBytes
type Hash = k.Hash

type BalancerCode int32

const (
	BalancerCodeRoundRobin BalancerCode = iota
	BalancerCodeLeastBytes
	BalancerCodeHash
)
