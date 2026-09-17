package key

import "testing"

func TestNamespaceKeyPrefix(t *testing.T) {
	t.Cleanup(func() { SetNamespace("") })

	SetNamespace("")
	if got := NamespaceKeyPrefix(); got != "" {
		t.Fatalf("empty namespace prefix = %q", got)
	}

	SetNamespace("local")
	if got := Namespace(); got != "local" {
		t.Fatalf("Namespace() = %q", got)
	}
	if got := NamespaceKeyPrefix(); got != "/local/" {
		t.Fatalf("NamespaceKeyPrefix() = %q", got)
	}

	key, err := NewKeyFromString("/users/alice")
	if err != nil {
		t.Fatal(err)
	}
	if got := key.String(); got != "/local/users/alice" {
		t.Fatalf("namespaced key = %q", got)
	}
}
