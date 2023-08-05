package mongo

import (
	"testing"

	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/document/key"
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

				if v, err := driver.Get(k, diface.WithDestination(&dest)); err != nil {
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
					diface.WithSource(data),
					diface.WithVersion(version),
				); err != nil {
					t.Error(err)
				} else {
					t.Log(v)
				}
			}
		}
	}
}
