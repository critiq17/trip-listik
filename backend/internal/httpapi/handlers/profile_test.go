package handlers

import "testing"

func TestComputeWorldPercent(t *testing.T) {
	tests := []struct {
		countries int64
		want      float64
	}{
		{0, 3.5},
		{195, 100.0},
		{39, 20.0},
	}

	for _, tt := range tests {
		got := computeWorldPercent(tt.countries)
		if got != tt.want {
			t.Fatalf("computeWorldPercent(%d) = %v, want %v", tt.countries, got, tt.want)
		}
	}
}
