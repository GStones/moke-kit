package handlers

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gstones/moke-kit/demo/internal/demo/db_nosql"
	"github.com/gstones/moke-kit/demo/internal/demo/db_sql"
	"github.com/gstones/moke-kit/mq/common"
	"github.com/gstones/moke-kit/mq/miface"
)

type Demo struct {
	logger  *zap.Logger
	nosqlDb db_nosql.Database
	mq      miface.MessageQueue
	gormDb  *gorm.DB
	redis   *redis.Client
}

func NewDemo(
	logger *zap.Logger,
	database db_nosql.Database,
	mq miface.MessageQueue,
	gormDb *gorm.DB,
	redis *redis.Client,
) *Demo {
	return &Demo{
		logger:  logger,
		nosqlDb: database,
		mq:      mq,
		gormDb:  gormDb,
		redis:   redis,
	}
}

func (d *Demo) Hi(uid, message string) error {
	// nosqlDb create
	if data, err := d.nosqlDb.LoadOrCreateDemo(uid); err != nil {
		return err
	} else {
		if err := data.Update(func() bool {
			data.SetMessage(message)
			return true
		}); err != nil {
			return err
		}
	}
	if err := db_sql.FirstOrCreate(d.gormDb, uid, message); err != nil {
		return err
	}
	d.redis.Set("demo", message, time.Minute)
	// mq publish
	if err := d.mq.Publish(
		common.NatsHeader.CreateTopic("demo"),
		miface.WithBytes([]byte(message)),
	); err != nil {
		return err
	}
	return nil
}

func (d *Demo) Watch(ctx context.Context, topic string, callback func(message string) error) error {
	// mq subscribe
	sub, err := d.mq.Subscribe(
		common.NatsHeader.CreateTopic(topic),
		func(msg miface.Message, err error) common.ConsumptionCode {
			if err := callback(string(msg.Data())); err != nil {
				return common.ConsumeNackPersistentFailure
			}
			return common.ConsumeAck
		})
	if err != nil {
		return err
	}
	<-ctx.Done()
	if err := sub.Unsubscribe(); err != nil {
		return err
	}
	return nil
}
