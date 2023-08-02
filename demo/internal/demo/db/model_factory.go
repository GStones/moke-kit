package db

import (
	"errors"
	"go.uber.org/zap"
	"moke-kit/nosql/document/diface"
	errors2 "moke-kit/nosql/errors"
)

//
//func  LoadOrCreateBuddyQueue(appId string, id string) (bq *BuddyQueue, err error) {
//	if bq, err = d.NewBuddyQueue(appId, id); err != nil {
//		return
//	} else if err = bq.Load(); errors.Cause(err) == errors2.ErrKeyNotFound {
//		if bq, err = d.NewBuddyQueue(appId, id); err != nil {
//			return
//		} else if err := bq.InitDefault(); err != nil {
//			return nil, err
//		} else if err = bq.Create(); err != nil {
//			err = bq.Load()
//		}
//	}
//	// Even if we know the BuddyQueue is on the latest version, we should still run it through fixups
//	if err == nil {
//		err = d.FixupBuddyQueue(bq)
//	}
//	return
//}
//

type Database struct {
	logger *zap.Logger
	db     diface.ICollection
}

func OpenDatabase(l *zap.Logger, db diface.ICollection) Database {
	return Database{
		logger: l,
		db:     db,
	}
}

func (db *Database) LoadOrCreateDemo(id string) (bq *DemoModel, err error) {
	if bq, err = NewDemoModel(id, db.db); err != nil {
		return
	} else if err = bq.Load(); errors.Is(err, errors2.ErrKeyNotFound) {
		if bq, err = NewDemoModel(id, db.db); err != nil {
			return
		} else if err := bq.InitDefault(); err != nil {
			return nil, err
		} else if err = bq.Create(); err != nil {
			err = bq.Load()
		}
	}
	return
}
