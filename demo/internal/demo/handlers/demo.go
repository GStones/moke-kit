package handlers

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"moke-kit/demo/internal/demo/db_nosql"
	"moke-kit/demo/internal/demo/db_sql"
	"moke-kit/mq/common"
	"moke-kit/mq/qiface"
)

type Demo struct {
	logger  *zap.Logger
	nosqlDb db_nosql.Database
	mq      qiface.MessageQueue
	gormDb  *gorm.DB
}

func NewDemo(
	logger *zap.Logger,
	database db_nosql.Database,
	mq qiface.MessageQueue,
	gormDb *gorm.DB,
) *Demo {
	return &Demo{
		logger:  logger,
		nosqlDb: database,
		mq:      mq,
		gormDb:  gormDb,
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
	if err := db_sql.FirstOrInit(d.gormDb, uid, message); err != nil {
		return err
	}

	// mq publish
	if err := d.mq.Publish(
		common.NatsHeader.CreateTopic("demo"),
		qiface.WithBytes([]byte(message)),
	); err != nil {
		return err
	}
	return nil
}

func (d *Demo) Watch(ctx context.Context, topic string, callback func(message string) error) error {
	// mq subscribe
	sub, err := d.mq.Subscribe(
		common.NatsHeader.CreateTopic(topic),
		func(msg qiface.Message, err error) common.ConsumptionCode {
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
