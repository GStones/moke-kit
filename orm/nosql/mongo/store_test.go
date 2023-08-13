package mongo

import (
	"moke-kit/orm/nosql/noptions"
	"testing"

	"moke-kit/orm/nosql/diface"
	"moke-kit/orm/nosql/key"
)

func TestNewDriverProvider(t *testing.T) {
	if mCli, err := NewMongoClient(
		"mongodb://localhost:27017",
		"", "",
	); err != nil {
		t.Error(err)
	} else {
		if provider := NewProvider(mCli, nil); provider == nil {
			t.Error("NewProvider error")
		} else if driver, err := provider.OpenDbDriver("test"); err != nil {
			t.Error(err)
		} else {
			if k, err := key.NewKeyFromParts("test", "10000"); err != nil {
				t.Error(err)
			} else {
				version := diface.Version(0)
				dest := map[string]interface{}{}

				if v, err := driver.Get(k, noptions.WithDestination(&dest)); err != nil {
					t.Error(err)
				} else {
					t.Log(v)
					version = v
				}

				data := map[string]interface{}{
					"test": "test",
					"num":  2,
				}

				if v, err := driver.Set(
					k,
					noptions.WithSource(data),
					noptions.WithVersion(version),
				); err != nil {
					t.Error(err)
				} else {
					t.Log(v)
				}
			}
		}
	}
}
