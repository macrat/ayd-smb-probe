//go:build integration

package main_test

import (
	"reflect"
	"testing"

	"github.com/macrat/ayd-smb-probe"
	"github.com/macrat/ayd/lib-ayd"
)

func TestCheck(t *testing.T) {
	t.Setenv("TZ", "UTC")

	tests := []struct {
		URL     string
		Message string
		Extra   map[string]interface{}
		Error   string
	}{
		{
			"smb://guest@localhost/",
			"server exists",
			map[string]interface{}{
				"shares_count": 2,
				"type":         "server",
			},
			"",
		},
		{
			"smb://guest@localhost/public",
			"directory exists",
			map[string]interface{}{
				"file_count": 1,
				"mtime":      "2023-01-01T00:00:00Z",
				"type":       "directory",
			},
			"",
		},
		{
			"smb://guest@localhost/public/test.txt",
			"file exists",
			map[string]interface{}{
				"file_size": int64(6),
				"mtime":     "2023-01-02T15:04:05Z",
				"type":      "file",
			},
			"",
		},
		{
			"smb://guest@localhost/public/hogefuga",
			"",
			nil,
			"stat hogefuga: file does not exist",
		},
		{
			"smb://guest@localhost/private",
			"",
			nil,
			"permission denied",
		},
		{
			"smb://foo@localhost/private",
			"",
			nil,
			"response error: The attempted logon is invalid. This is either due to a bad username or authentication information.",
		},
		{
			"smb://foo:baz@localhost/private",
			"",
			nil,
			"response error: The attempted logon is invalid. This is either due to a bad username or authentication information.",
		},
		{
			"smb://foo:bar@localhost/private",
			"directory exists",
			map[string]interface{}{
				"file_count": 2,
				"mtime":      "2023-01-01T00:00:00Z",
				"type":       "directory",
			},
			"",
		},
		{
			"smb://foo:bar@localhost/private/my-secret.txt",
			"file exists",
			map[string]interface{}{
				"file_size": int64(0),
				"mtime":     "2023-01-01T01:01:01Z",
				"type":      "file",
			},
			"",
		},
		{
			"smb://foo:bar@localhost/no-such-share",
			"",
			nil,
			"response error: {Network Name Not Found} The specified share name cannot be found on the remote server.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.URL, func(t *testing.T) {
			u, err := ayd.ParseURL(tt.URL)
			if err != nil {
				t.Fatalf("failed to parse URL: %s", err)
			}

			msg, extra, _, _, err := main.Check(u)

			if tt.Error == "" {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			} else {
				if err == nil {
					t.Fatalf("unexpected error\nexpected: %q\n but got: nil", tt.Error)
				} else if err.Error() != tt.Error {
					t.Fatalf("unexpected error\nexpected: %q\n but got: %q", tt.Error, err.Error())
				}
			}

			if msg != tt.Message {
				t.Errorf("unexpected message\nexpected: %q\n but got: %q", tt.Message, msg)
			}
			if !reflect.DeepEqual(extra, tt.Extra) {
				t.Errorf("unexpected extra\nexpected: %#v\n but got: %#v", tt.Extra, extra)
			}
		})
	}
}
