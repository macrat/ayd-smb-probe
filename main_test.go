package main_test

import (
	"testing"

	"github.com/macrat/ayd-smb-probe"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
		Error  string
	}{
		{"smb://hello:world@example.com", "smb://hello:world@example.com:445", ""},
		{"smb://foo:bar@127.0.0.1:1234", "smb://foo:bar@127.0.0.1:1234", ""},
		{"smb://example.com", "smb://guest@example.com:445", ""},
		{"smb:", "", "invalid target URI: hostname is required"},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			u, err := main.ParseTarget(tt.Input)
			if err != nil {
				if err.Error() != tt.Error {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			} else if tt.Error != "" {
				t.Fatal("expected error but got nil")
			}

			if u.String() != tt.Output {
				t.Errorf("unexpected output: %s", u)
			}
		})
	}
}
