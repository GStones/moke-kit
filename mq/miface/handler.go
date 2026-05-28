package miface

import (
	"context"

	"github.com/gstones/moke-kit/mq/common"
)

type SubResponseHandler = func(context context.Context, msg Message, err error) common.ConsumptionCode
