package mfx

import (
	"testing"

	"github.com/gstones/moke-kit/utility"
)

func TestCreateAppModuleSetsSharedNamespaceFromDeployment(t *testing.T) {
	t.Setenv("DEPLOYMENT", "PROD")
	t.Cleanup(func() { utility.SetNamespace("") })

	out, err := CreateAppModule()
	if err != nil {
		t.Fatal(err)
	}
	if out.Deployment != "PROD" {
		t.Fatalf("Deployment = %q, want PROD", out.Deployment)
	}
	if got := utility.Namespace(); got != "PROD" {
		t.Fatalf("Namespace() = %q, want PROD", got)
	}
}
