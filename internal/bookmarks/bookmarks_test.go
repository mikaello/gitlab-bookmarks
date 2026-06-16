package bookmarks

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func TestCreateBookmarkHTMLFlat(t *testing.T) {
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

	html, err := CreateBookmarkHTML(projects, FolderModeFlat)
	if err != nil {
		t.Fatalf("CreateBookmarkHTML returned error: %s", err)
	}

	for _, want := range []string{
		"<!DOCTYPE NETSCAPE-Bookmark-file-1>",
		"https://gitlab.example.com/group/alpha",
		"ADD_DATE=\"1700000000\">alpha</A>",
		">alpha</A>",
		"https://gitlab.example.com/group/beta",
		">beta</A>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, html)
		}
	}
}

func TestCreateBookmarkHTMLNamespaceFolders(t *testing.T) {
	activity := time.Unix(1700000000, 0)
	projects := []*gitlab.Project{
		{
			Name:           "alpha",
			WebURL:         "https://gitlab.example.com/group/team/alpha",
			LastActivityAt: &activity,
			Namespace:      &gitlab.ProjectNamespace{FullPath: "group/team"},
		},
		{
			Name:           "beta",
			WebURL:         "https://gitlab.example.com/group/platform/beta",
			LastActivityAt: &activity,
			Namespace:      &gitlab.ProjectNamespace{FullPath: "group/platform"},
		},
		{
			Name:           "gamma",
			WebURL:         "https://gitlab.example.com/group/team/gamma",
			LastActivityAt: &activity,
			Namespace:      &gitlab.ProjectNamespace{FullPath: "group/team"},
		},
	}

	html, err := CreateBookmarkHTML(projects, FolderModeNamespace)
	if err != nil {
		t.Fatalf("CreateBookmarkHTML returned error: %s", err)
	}

	for _, want := range []string{
		">group/team</H3>",
		">group/platform</H3>",
		"https://gitlab.example.com/group/team/alpha",
		"https://gitlab.example.com/group/team/gamma",
		"https://gitlab.example.com/group/platform/beta",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, html)
		}
	}
}

func TestCreateBookmarkHTMLWithoutLastActivityAt(t *testing.T) {
	projects := []*gitlab.Project{
		{
			Name:   "alpha",
			WebURL: "https://gitlab.example.com/group/alpha",
		},
	}

	html, err := CreateBookmarkHTML(projects, FolderModeFlat)
	if err != nil {
		t.Fatalf("CreateBookmarkHTML returned error: %s", err)
	}

	if !strings.Contains(html, "ADD_DATE=\"0\">alpha</A>") {
		t.Errorf("expected output to use a zero add date for missing activity, got:\n%s", html)
	}
}

func TestParseFolderMode(t *testing.T) {
	for _, value := range []string{"flat", "namespace"} {
		if _, err := ParseFolderMode(value); err != nil {
			t.Fatalf("ParseFolderMode(%q) returned error: %s", value, err)
		}
	}

	if _, err := ParseFolderMode("group"); err == nil {
		t.Fatal("expected ParseFolderMode to reject invalid mode")
	}
}

func TestWriteBookmarkFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.html")
	if err := WriteBookmarkFile(path, "hello"); err != nil {
		t.Fatalf("WriteBookmarkFile returned error: %s", err)
	}
}
