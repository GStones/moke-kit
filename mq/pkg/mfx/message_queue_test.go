package mfx

import (
	"testing"

	"github.com/gstones/moke-kit/utility"
)

func TestCreateMessageQueueModuleDoesNotOverwriteNamespace(t *testing.T) {
	t.Cleanup(func() { utility.SetNamespace("") })
	utility.SetNamespace("PROD")

	if _, err := CreateMessageQueueModule(MQImplementations{}); err != nil {
		t.Fatal(err)
	}
	if got := utility.Namespace(); got != "PROD" {
		t.Fatalf("Namespace() = %q after MQ init, want PROD", got)
	}
}
