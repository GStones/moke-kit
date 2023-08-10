package demo

import "moke-kit/gorm/nosql/key"

func NewDemoKey(id string) (key.Key, error) {
	return key.NewKeyFromParts("demo", id)
}
