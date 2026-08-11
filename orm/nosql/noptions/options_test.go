package noptions

import (
	"testing"

	"github.com/gstones/moke-kit/orm/nerrors"
)

func TestWithVersionAndAnyVersionConflict(t *testing.T) {
	_, err := NewOptions(WithVersion(1), WithAnyVersion())
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if err != nerrors.ErrAnyVersionConflict {
		t.Fatalf("got %v, want %v", err, nerrors.ErrAnyVersionConflict)
	}
}

func TestDefaultOptions(t *testing.T) {
	o, err := NewOptions()
	if err != nil {
		t.Fatal(err)
	}
	if o.Version != NoVersion || o.AnyVersion {
		t.Fatalf("unexpected defaults: %+v", o)
	}
}
