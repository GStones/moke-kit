package miface

import (
	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
)

// SubOptions contains all the various options that the provided WithXyz functions construct.
type SubOptions struct {
	DeliverySemantics common.DeliverySemantics
	GroupId           string
}

// SubOption is a closure that updates SubOptions.
type SubOption func(o *SubOptions) error

// NewSubOptions constructs an SubOptions struct from the provided SubOption closures and returns it.
func NewSubOptions(opts ...SubOption) (options SubOptions, err error) {
	o := &options

	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

// Configures delivery semantics of subscription to be at-least-once delivery.
// If a semantics preference is not set, mq implementation will use its default mode.
// Mutually exclusive with WithAtMostOnceDelivery()
func WithAtLeastOnceDelivery() SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != common.Unset {
			return qerrors.ErrSemanticsAlreadySet
		} else {
			o.DeliverySemantics = common.AtLeastOnce
			return nil
		}
	}
}

// Configures delivery semantics of subscription to be at-most-once delivery.
// If a semantics preference is not set, mq implementation will use its default mode.
// groupId is also optional. Pass in mq.DefaultId to have the mq implementation set a default groupId.
// Mutually exclusive with WithAtLeastOnceDelivery()
func WithAtMostOnceDelivery(groupId common.GroupId) SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != common.Unset {
			return qerrors.ErrSemanticsAlreadySet
		} else {
			o.DeliverySemantics = common.AtMostOnce
			o.GroupId = string(groupId)
			return nil
		}
	}
}
