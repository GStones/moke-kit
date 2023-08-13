package db_nosql

import (
	"errors"
	"moke-kit/demo/internal/demo/db_nosql/demo"
	"moke-kit/orm/nosql/diface"

	"go.uber.org/zap"

	"moke-kit/orm/nerrors"
)

type Database struct {
	logger *zap.Logger
	coll   diface.ICollection
}

func OpenDatabase(l *zap.Logger, coll diface.ICollection) Database {
	return Database{
		logger: l,
		coll:   coll,
	}
}

func (db *Database) LoadOrCreateDemo(id string) (*demo.Dao, error) {
	if dm, err := demo.NewDemoModel(id, db.coll); err != nil {
		return nil, err
	} else if err = dm.Load(); errors.Is(err, nerrors.ErrNotFound) {
		if dm, err = demo.NewDemoModel(id, db.coll); err != nil {
			return nil, err
		} else if err := dm.InitDefault(); err != nil {
			return nil, err
		} else if err = dm.Create(); err != nil {
			err = dm.Load()
		} else {
			return dm, nil
		}
	} else if err != nil {
		return nil, err
	} else {
		return dm, nil
	}
	return nil, nil
}
