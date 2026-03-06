package server

import (
	"testing"
)

func TestIsSameBaseDomain(t *testing.T) {
	tests := []struct {
		origin     string
		baseDomain string
		want       bool
	}{
		{"https://app.example.com", "example.com", true},
		{"http://localhost:3000", "example.com", false},
		{"https://example.com", "example.com", true},
		{"https://sub.sub.example.com", "example.com", true},
		{"https://otherdomain.com", "example.com", false},
		{"https://example.com.other.com", "example.com", false},
		{"", "example.com", false},
		{"https://app.example.com", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.origin+"_"+tt.baseDomain, func(t *testing.T) {
			if got := isSameBaseDomain(tt.origin, tt.baseDomain); got != tt.want {
				t.Errorf("isSameBaseDomain(%q, %q) = %v, want %v", tt.origin, tt.baseDomain, got, tt.want)
			}
		})
	}
}
