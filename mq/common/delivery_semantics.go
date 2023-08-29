package common

type DeliverySemantics string

const (
	AtLeastOnce DeliverySemantics = "at-least-once"
	AtMostOnce  DeliverySemantics = "at-most-once"
	Unset       DeliverySemantics = ""
)

type GroupId string

const (
	DefaultId GroupId = ""
)
