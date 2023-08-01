package mq

type SubResponseHandler = func(msg Message, err error) ConsumptionCode
