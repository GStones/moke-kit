package db_nosql

import (
	"errors"

	"github.com/gstones/moke-kit/demo/internal/demo/db_nosql/demo"
	"github.com/gstones/moke-kit/orm/nosql/diface"

	"go.uber.org/zap"

	"github.com/gstones/moke-kit/orm/nerrors"
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
			if err = dm.Load(); err != nil {
				return nil, err
			} else {
				return dm, nil
			}
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
