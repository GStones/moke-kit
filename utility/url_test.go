package utility

import "testing"

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"no userinfo", "mongodb://localhost:27017/db", "mongodb://localhost:27017/db"},
		{"password", "redis://user:secret@host:6379/0", "redis://user:xxxxx@host:6379/0"},
		{"username only", "nats://user@host:4222", "nats://user@host:4222"},
		{"mongo srv", "mongodb+srv://u:p@cluster.example/db", "mongodb+srv://u:xxxxx@cluster.example/db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactURL(tt.in); got != tt.want {
				t.Fatalf("RedactURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
