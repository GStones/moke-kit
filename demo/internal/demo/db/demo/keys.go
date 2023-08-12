package demo

import "moke-kit/nsorm/nosql/key"

func NewDemoKey(id string) (key.Key, error) {
	return key.NewKeyFromParts("demo", id)
}
