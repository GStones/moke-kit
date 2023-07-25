package document

import (
	"go.uber.org/fx"
)

// FxStartingDocumentModule adds any starting documents provided by FxProvideStartingDocuments to the document store
// provider in the dependency graph.
var FxStartingDocumentsModule = fx.Invoke(
	func(in FxStartingDocumentsIn) error {
		for _, d := range in.StartingDocuments {
			if err := in.DocumentStoreProvider.AddStartingDocuments(d); err != nil {
				return err
			}
		}
		return nil
	},
)

// FxStartingDocumentsIn collects StartingDocument instances from the fx dependency graph.
type FxStartingDocumentsIn struct {
	fx.In
	DocumentStoreProvider DocumentStoreProvider `name:"DocumentStoreProvider"`
	StartingDocuments     [][]StartingDocument  `group:"StartingDocuments"`
}

// FxStartingDocumentsOut contains starting documents that should be injected into the fx dependency graph.
type FxStartingDocumentsOut struct {
	fx.Out
	StartingDocuments []StartingDocument `group:"StartingDocuments"`
}
