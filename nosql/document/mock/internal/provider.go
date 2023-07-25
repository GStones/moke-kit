package internal

import (
	document2 "github.com/gstones/platform/services/common/nosql/document"
	"github.com/gstones/platform/services/common/nosql/errors"
)

type DocumentStoreProvider struct {
	stores map[string]*DocumentStore
}

func NewDocumentStoreProvider() (*DocumentStoreProvider, error) {
	result := &DocumentStoreProvider{
		stores: map[string]*DocumentStore{},
	}
	return result, nil
}

func (d *DocumentStoreProvider) OpenDocumentStore(name string) (document2.DocumentStore, error) {
	if name == "" {
		return nil, errors.ErrInvalidStoreName
	}
	ds := d.stores[name]
	if ds == nil {
		ds = &DocumentStore{
			name:      name,
			documents: map[document2.Key]*document{},
		}
		d.stores[name] = ds
	}
	return ds, nil
}

func (d *DocumentStoreProvider) Shutdown() error {
	return nil
}

func (d *DocumentStoreProvider) AddStartingDocuments(documents []document2.StartingDocument) error {
	for _, doc := range documents {
		if ds, err := d.OpenDocumentStore(doc.Store); err != nil {
			return err
		} else {
			if doc.Key != "" || doc.Data != nil {
				if k, err := document2.NewKeyFromString(doc.Key); err != nil {
					return err
				} else if _, err := ds.Set(
					k,
					document2.WithSource(doc.Data),
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
