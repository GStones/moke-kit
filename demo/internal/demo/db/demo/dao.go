package demo

import (
	"moke-kit/gorm/nosql"
	"moke-kit/gorm/nosql/diface"
)

type Dao struct {
	nosql.DocumentBase `bson:"-"`
	Id                 string `bson:"_id"`
	appId              string `bson:"-"`
	Data               *Data  `bson:"data"`
}

func (dm *Dao) Init(id string, doc diface.ICollection) error {
	key, e := NewDemoKey(id)
	if e != nil {
		return e
	}

	dm.Data = &Data{}
	dm.DocumentBase.Init(&dm.Data, dm.clear, doc, key)
	return nil
}

func (dm *Dao) SetMessage(message string) {
	dm.Data.Message = message
}

func (dm *Dao) clear() {
	dm.Data = nil
}

func (dm *Dao) InitDefault() error {

	return nil
}

func NewDemoModel(id string, doc diface.ICollection) (*Dao, error) {
	dm := &Dao{}
	if err := dm.Init(id, doc); err != nil {
		return nil, err
	}
	return dm, nil
}
