package store

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"plain", "plain"},
		{"100% legit", "100\\% legit"},
		{"under_score", "under\\_score"},
		{`path\to`, `path\\to`},
		{"%_\\", "\\%\\_\\\\"},
	}

	for _, tt := range tests {
		if got := escapeLike(tt.in); got != tt.want {
			t.Fatalf("escapeLike(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
