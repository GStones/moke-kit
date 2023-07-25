package internal

import (
	"github.com/gstones/platform/services/common/nosql/document"
	"testing"
)

var doc document.DocumentStore

func TestMain(m *testing.M) {
	provider, err := NewDocumentStoreProvider(
		"mongodb://"+
			"192.168.37.135:27017,"+
			"192.168.37.135:27018,"+
			"192.168.37.135:27019/"+
			"?replicaSet=rs001",
		"admin", "wg1q2w3e", nil)
	if err != nil {
		panic(err)
	}
	doc, err = provider.OpenDocumentStore("game")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestNewDocumentStore(t *testing.T) {
	key, err := document.NewKeyFromParts("test", "10000")
	if err != nil {
		t.Fatal(err)
	}
	test := map[string]string{
		"test":  "test1234",
		"test2": "test1234",
		"test3": "test1234",
	}
	_, err = doc.Set(key, document.WithSource(test))
	if err != nil {
		t.Fatal(err)
	}
	dest := make(map[string]string)
	_, err = doc.Get(key, document.WithDestination(&dest))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dest)
}

func TestWatch(t *testing.T) {
	key, err := document.NewKeyFromParts("test", "10000")
	if err != nil {
		t.Fatal(err)
	}

	watcher := document.NewDocWatcher(key)
	go func() {
		watcher.AddCallBack(func(key document.Key, value interface{}) {
			t.Log("watcher", key.String(), value)
		})

		doc.AddWatcher(key, watcher)
	}()

	select {}
}
