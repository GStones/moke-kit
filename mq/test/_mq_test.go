package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/fxmain"
	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
	"github.com/gstones/moke-kit/mq/pkg/mfx"
)

// TestLocalMQ is a test local(channel) mq.
func TestLocalMQ(t *testing.T) {
	fxmain.Main(
		//inject local(channel) message queue
		mfx.LocalModule,
		//need use local mq app module
		AppLocalMQModule,
	)
}

// TestNatsMQ is a test nats mq.
// Tips: require nats server running.
func TestNatsMQ(t *testing.T) {
	fxmain.Main(
		//inject nats message queue
		mfx.NatsModule,
		//need use nats mq app module
		AppNatsMQModule,
	)
}

// AppLocalMQModule is a test local mq application module.
var AppLocalMQModule = fx.Invoke(
	func(
		l *zap.Logger, //inject logger(default)
		mqParams mfx.MessageQueueParams, //inject message queue
	) error {
		topic := "test"
		//create a local mq subscribe topic
		localTopic := common.LocalHeader.CreateTopic(topic)
		// create a context to control the subscription lifecycle
		ctx := context.Background()
		if _, err := mqParams.MessageQueue.Subscribe(
			ctx, localTopic,
			func(msg miface.Message, error error) common.ConsumptionCode {
				l.Info("local mq consume1", zap.String("msg", string(msg.Data())))
				ctx.Done()
				return common.ConsumeAck
			},
		); err != nil {
			ctx.Done()
			return err
		}

		if _, err := mqParams.MessageQueue.Subscribe(
			ctx, localTopic,
			func(msg miface.Message, error error) common.ConsumptionCode {
				l.Info("local mq consume2", zap.String("msg", string(msg.Data())))
				ctx.Done()
				return common.ConsumeAck
			},
		); err != nil {
			ctx.Done()
			return err
		}

		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			if err := mqParams.MessageQueue.Publish(
				localTopic,
				miface.WithBytes([]byte(fmt.Sprintf("local mq: %s-%d", topic, i))),
			); err != nil {
				return err
			}
		}
		log.Fatalf("local mq test done")
		return nil
	},
)

// AppNatsMQModule is a test nats mq application module.
var AppNatsMQModule = fx.Invoke(
	func(
		l *zap.Logger, //inject logger(default)
		mqParams mfx.MessageQueueParams, //inject message queue
	) error {
		topic := "test"
		//create a nats mq subscribe topic
		natsTopic := common.NatsHeader.CreateTopic(topic)
		// create a context to control the subscription lifecycle
		ctx := context.Background()
		if _, err := mqParams.MessageQueue.Subscribe(
			ctx, natsTopic,
			func(msg miface.Message, error error) common.ConsumptionCode {
				l.Info("nats mq consume", zap.String("msg", string(msg.Data())))
				ctx.Done()
				return common.ConsumeAck
			},
		); err != nil {
			ctx.Done()
			return err
		}

		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			if err := mqParams.MessageQueue.Publish(
				natsTopic,
				miface.WithBytes([]byte(fmt.Sprintf("nats mq: %s-%d", topic, i))),
			); err != nil {
				return err
			}
		}
		log.Fatalf("nats mq test done")
		return nil
	})
