package main_test

import (
	"net/url"
	"testing"

	"github.com/macrat/ayd-smb-probe"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{"smb://hello:world@example.com", "smb://hello:world@example.com/"},
		{"smb://foo:bar@127.0.0.1:1234", "smb://foo:bar@127.0.0.1:1234/"},
		{"smb://example.com", "smb://guest@example.com/"},
		{"smb://example.com/path/to#abc#def=ghi", "smb://guest@example.com/path/to"},
		{"smb://example.com/path", "smb://guest@example.com/path"},
		{"smb://example.com/a/../b/", "smb://guest@example.com/b"},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			u, err := url.Parse(tt.Input)
			if err != nil {
				t.Fatalf("failed to parse input url: %s", err)
			}

			u = main.NormalizeTarget(u)

			if u.String() != tt.Output {
				t.Errorf("expected %s but got %s", tt.Output, u)
			}
		})
	}
}
