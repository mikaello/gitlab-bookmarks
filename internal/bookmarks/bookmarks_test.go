package bookmarks

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func TestCreateBookmarkHTML(t *testing.T) {
	activity := time.Unix(1700000000, 0)
	projects := []*gitlab.Project{
		{
			Name:           "alpha",
			WebURL:         "https://gitlab.example.com/group/alpha",
			LastActivityAt: &activity,
		},
		{
			Name:           "beta",
			WebURL:         "https://gitlab.example.com/group/beta",
			LastActivityAt: &activity,
		},
	}

	html, err := CreateBookmarkHTML(projects)
	if err != nil {
		t.Fatalf("CreateBookmarkHTML returned error: %s", err)
	}

	for _, want := range []string{
		"<!DOCTYPE NETSCAPE-Bookmark-file-1>",
		"https://gitlab.example.com/group/alpha",
		">alpha</A>",
		"https://gitlab.example.com/group/beta",
		">beta</A>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, html)
		}
	}
}

func TestWriteBookmarkFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.html")
	if err := WriteBookmarkFile(path, "hello"); err != nil {
		t.Fatalf("WriteBookmarkFile returned error: %s", err)
	}
}
