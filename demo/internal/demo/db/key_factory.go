package db

import "moke-kit/nosql/document/key"

func NewDemoKey(id string) (key.Key, error) {
	return key.NewKeyFromParts("demo", id)
}
