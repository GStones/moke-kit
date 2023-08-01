package mq

type DeliverySemantics string

const (
	AtLeastOnce DeliverySemantics = "at-least-once"
	AtMostOnce                    = "at-most-once"
	Unset                         = ""
)

type GroupId string

const (
	DefaultId GroupId = ""
)
