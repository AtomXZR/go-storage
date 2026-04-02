package storage_test

import (
	"testing"

	"github.com/AtomXZR/go-storage"
)

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// --- Valid: no leading slash ---
		{
			name:  "simple key",
			input: "foo.txt",
			want:  "foo.txt",
		},
		{
			name:  "nested path",
			input: "images/avatars/foo.png",
			want:  "images/avatars/foo.png",
		},
		{
			name:  "deep nesting",
			input: "a/b/c/d/e.txt",
			want:  "a/b/c/d/e.txt",
		},
		{
			name:  "file with dots in name",
			input: "archive.tar.gz",
			want:  "archive.tar.gz",
		},
		{
			name:  "path with dots in name",
			input: "releases/v1.2.3/binary.tar.gz",
			want:  "releases/v1.2.3/binary.tar.gz",
		},
		{
			name:  "leading slash simple",
			input: "/foo.txt",
			want:  "foo.txt",
		},
		{
			name:  "leading slash nested",
			input: "/images/avatars/foo.png",
			want:  "images/avatars/foo.png",
		},
		{
			name:  "double slash",
			input: "images//foo.png",
			want:  "images/foo.png",
		},
		{
			name:  "triple slash",
			input: "a///b///c.txt",
			want:  "a/b/c.txt",
		},
		{
			name:  "leading double slash",
			input: "//foo.txt",
			want:  "foo.txt",
		},
		{
			name:  "dot segment in middle",
			input: "images/./foo.png",
			want:  "images/foo.png",
		},
		{
			name:  "multiple dot segments",
			input: "a/./b/./c.txt",
			want:  "a/b/c.txt",
		},
		{
			name:  "parent then child (resolves within root)",
			input: "images/../images/foo.png",
			want:  "images/foo.png",
		},
		{
			name:  "redundant parent traversal that stays in bounds",
			input: "a/b/../b/c.txt",
			want:  "a/b/c.txt",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "root slash only",
			input:   "/",
			wantErr: true,
		},
		{
			name:    "traversal at start",
			input:   "../etc/passwd",
			wantErr: true,
		},
		{
			name:  "traversal with leading slash",
			input: "/../etc/passwd",
			want:  "etc/passwd",
		},
		{
			name:    "traversal escapes deep path",
			input:   "a/b/../../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "just double dot",
			input:   "..",
			wantErr: true,
		},
		{
			name:    "just dot",
			input:   ".",
			wantErr: true,
		},
		{
			name:    "leading slash just dot",
			input:   "/.",
			wantErr: true,
		},
		{
			name:    "leading slash just double dot",
			input:   "/..",
			wantErr: true,
		},
		{
			name:  "trailing slash stripped by path.Clean",
			input: "images/",
			want:  "images",
		},
		{
			name:  "nested trailing slash",
			input: "images/avatars/",
			want:  "images/avatars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.NormalizeKey(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeKey(%q) expected error, got %q", tt.input, got)
				}
				return
			}

			if err != nil {
				t.Errorf("normalizeKey(%q) unexpected error: %v", tt.input, err)
				return
			}

			if got != tt.want {
				t.Errorf("normalizeKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
