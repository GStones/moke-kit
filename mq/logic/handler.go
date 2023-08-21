package logic

import (
	"github.com/gstones/moke-kit/mq/common"
)

type SubResponseHandler = func(msg Message, err error) common.ConsumptionCode
