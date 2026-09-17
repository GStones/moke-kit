package utility

import (
	"sync"
	"testing"
)

func TestNamespaceConcurrentSetAndLoad(t *testing.T) {
	t.Cleanup(func() { SetNamespace("") })

	SetNamespace("local")
	if got := Namespace(); got != "local" {
		t.Fatalf("Namespace() = %q, want local", got)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			SetNamespace("prod")
			_ = Namespace()
		}()
	}
	wg.Wait()

	if got := Namespace(); got != "prod" {
		t.Fatalf("Namespace() = %q, want prod", got)
	}
}
