package logic

import (
	"moke-kit/mq/common"
)

type SubResponseHandler = func(msg Message, err error) common.ConsumptionCode
