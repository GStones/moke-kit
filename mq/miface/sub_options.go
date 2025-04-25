package miface

import (
	"time"

	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/internal/qerrors"
)

// DefaultSubOptions 返回默认的订阅选项
func DefaultSubOptions() SubOptions {
	return SubOptions{
		DeliverySemantics: common.Unset,
		GroupId:           "",
		Concurrency:       1,           // 默认单并发
		MaxRetries:        3,           // 默认最多重试3次
		RetryDelay:        time.Second, // 默认重试延迟1秒
		Timeout:           time.Minute, // 默认处理超时1分钟
		AutoAck:           true,        // 默认自动确认
		DLQEnabled:        false,       // 默认不启用死信队列
		DLQTopic:          "",
	}
}

// SubOptions 包含由 WithXyz 函数构造的各种订阅选项
type SubOptions struct {
	DeliverySemantics common.DeliverySemantics
	GroupId           string
	Concurrency       int           // 订阅的并发数
	MaxRetries        int           // 处理失败时的最大重试次数
	RetryDelay        time.Duration // 处理失败时的重试延迟
	Timeout           time.Duration // 处理超时时间
	AutoAck           bool          // 是否自动确认消息
	DLQEnabled        bool          // 是否启用死信队列
	DLQTopic          string        // 死信队列主题
}

// SubOption 是更新 SubOptions 的闭包
type SubOption func(o *SubOptions) error

// NewSubOptions 从提供的 SubOption 闭包构造 SubOptions 结构并返回
func NewSubOptions(opts ...SubOption) (options SubOptions, err error) {
	options = DefaultSubOptions()
	o := &options

	for _, opt := range opts {
		if err = opt(o); err != nil {
			return options, err
		}
	}

	// 可以在这里添加验证逻辑
	if err = o.validate(); err != nil {
		return options, err
	}

	return options, nil
}

// validate 验证选项是否有效
func (o *SubOptions) validate() error {
	// 添加必要的验证逻辑
	// 例如，如果 DeliverySemantics 是 AtMostOnce，GroupId 不能为空
	if o.DeliverySemantics == common.AtMostOnce && o.GroupId == "" {
		return qerrors.ErrInvalidGroupId
	}
	return nil
}

// WithAtLeastOnceDelivery 配置订阅的交付语义为至少一次交付
// 如果没有设置语义首选项，mq 实现将使用其默认模式
// 与 WithAtMostOnceDelivery() 互斥
func WithAtLeastOnceDelivery() SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != common.Unset {
			return qerrors.ErrSemanticsAlreadySet
		}
		o.DeliverySemantics = common.AtLeastOnce
		return nil
	}
}

// WithAtMostOnceDelivery 配置订阅的交付语义为最多一次交付
// 如果没有设置语义首选项，mq 实现将使用其默认模式
// groupId 也是可选的。传入 mq.DefaultId 让 mq 实现设置默认 groupId
// 与 WithAtLeastOnceDelivery() 互斥
func WithAtMostOnceDelivery(groupId common.GroupId) SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != common.Unset {
			return qerrors.ErrSemanticsAlreadySet
		}
		o.DeliverySemantics = common.AtMostOnce
		o.GroupId = string(groupId)
		return nil
	}
}

// WithGroup 设置订阅的 GroupId
// 注意：对于 AtMostOnce 语义，GroupId 在 WithAtMostOnceDelivery 中已设置
func WithGroup(groupId common.GroupId) SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics == common.AtMostOnce {
			return qerrors.ErrGroupAlreadySet
		}
		o.GroupId = string(groupId)
		return nil
	}
}

// WithConcurrency 设置订阅的并发数
// 注意：如果设置了并发数，mq 实现可能会使用不同的实现来处理消息
// 例如，Kafka 可能会使用多个消费者组来处理消息
// 这可能会影响消息的顺序性和处理性能
func WithConcurrency(concurrency int) SubOption {
	return func(o *SubOptions) error {
		if concurrency <= 0 {
			return qerrors.ErrInvalidGroupId
		}
		o.Concurrency = concurrency
		return nil
	}
}

// WithMaxRetries 设置处理失败时的最大重试次数
// 注意：如果设置了最大重试次数，mq 实现可能会使用不同的实现来处理消息
// 例如，Kafka 可能会使用重试主题来处理消息
// 这可能会影响消息的顺序性和处理性能
func WithMaxRetries(maxRetries int) SubOption {
	return func(o *SubOptions) error {
		if maxRetries < 0 {
			return qerrors.ErrInvalidGroupId
		}
		o.MaxRetries = maxRetries
		return nil
	}
}

// WithRetryDelay 设置处理失败时的重试延迟
func WithRetryDelay(retryDelay time.Duration) SubOption {
	return func(o *SubOptions) error {
		if retryDelay < 0 {
			return qerrors.ErrInvalidGroupId
		}
		o.RetryDelay = retryDelay
		return nil
	}
}

// WithTimeout 设置处理超时时间
func WithTimeout(timeout time.Duration) SubOption {
	return func(o *SubOptions) error {
		if timeout < 0 {
			return qerrors.ErrInvalidGroupId
		}
		o.Timeout = timeout
		return nil
	}
}

// WithAutoAck 设置是否自动确认消息
func WithAutoAck(autoAck bool) SubOption {
	return func(o *SubOptions) error {
		o.AutoAck = autoAck
		return nil
	}
}

// WithDLQEnabled 设置是否启用死信队列
func WithDLQEnabled(enabled bool) SubOption {
	return func(o *SubOptions) error {
		o.DLQEnabled = enabled
		return nil
	}
}

// WithDLQTopic 设置死信队列主题
func WithDLQTopic(topic string) SubOption {
	return func(o *SubOptions) error {
		if topic == "" {
			return qerrors.ErrEmptyTopic
		}
		o.DLQTopic = topic
		return nil
	}
}
