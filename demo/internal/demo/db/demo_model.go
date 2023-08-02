package db

import (
	"moke-kit/nosql/document"
	"moke-kit/nosql/document/diface"
)

type DemoModel struct {
	document.DocumentBase
	appId   string
	Message string
}

func (dm *DemoModel) Init(id string, doc diface.ICollection) error {
	key, e := NewDemoKey(id)
	if e != nil {
		return e
	}
	dm.DocumentBase.Init(dm, dm.clear, key, doc)
	return nil
}

func (dm *DemoModel) SetMessage(message string) {
	dm.Message = message
}

func (dm *DemoModel) clear() {

}

func (dm *DemoModel) InitDefault() error {
	dm.SetMessage("hello world")
	return nil
}

func NewDemoModel(id string, doc diface.ICollection) (dm *DemoModel, err error) {
	dm = &DemoModel{}
	if err = dm.Init(id, doc); err != nil {
		return
	}
	return
}
