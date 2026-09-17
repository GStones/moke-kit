package authfx

import "sync"

// unauthTracker holds a thread-safe set of gRPC method paths that bypass authentication.
type unauthTracker struct {
	mu      sync.RWMutex
	methods map[string]struct{}
}

func (u *unauthTracker) isUnauth(method string) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, skip := u.methods[method]
	return skip
}

// AddUnAuthMethod registers a gRPC method path to skip authentication checks.
func (u *unauthTracker) AddUnAuthMethod(method string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.methods == nil {
		u.methods = make(map[string]struct{})
	}
	u.methods[method] = struct{}{}
}

func newUnauthTracker() unauthTracker {
	return unauthTracker{methods: make(map[string]struct{})}
}
