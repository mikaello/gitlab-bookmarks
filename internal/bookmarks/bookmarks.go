package bookmarks

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"

	"gitlab.com/gitlab-org/api/client-go/v2"
)

//go:embed bookmarks.tmpl
var bookmarksTmpl embed.FS

// Credit: https://betterprogramming.pub/how-to-generate-html-with-golang-templates-5fad0d91252
// and https://pkg.go.dev/html/template

// CreateBookmarkHTML creates a bookmark content for the given repositories.
func CreateBookmarkHTML(projects []*gitlab.Project) (string, error) {
	templates := template.Must(template.New("").ParseFS(bookmarksTmpl, "bookmarks.tmpl"))

	var processed bytes.Buffer
	if err := templates.ExecuteTemplate(&processed, "bookmarks", projects); err != nil {
		return "", fmt.Errorf("execute bookmarks template: %w", err)
	}

	return processed.String(), nil
}

// WriteBookmarkFile writes 'filename' to disk with content 'htmlContent'.
func WriteBookmarkFile(filename string, htmlContent string) error {
	if err := os.WriteFile(filename, []byte(htmlContent), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}
	return nil
}
