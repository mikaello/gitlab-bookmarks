package bookmarks

import (
	"bufio"
	"bytes"
	"embed"
	"html/template"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

//go:embed bookmarks.tmpl
var bookmarksTmpl embed.FS

// Credit: https://betterprogramming.pub/how-to-generate-html-with-golang-templates-5fad0d91252
// and https://pkg.go.dev/html/template

// CreateBookmarkHTML creates a bookmark content for the given repositories.
func CreateBookmarkHTML(projects []*gitlab.Project) string {
	templates := template.Must(template.New("").ParseFS(bookmarksTmpl, "bookmarks.tmpl"))

	var processed bytes.Buffer
	err := templates.ExecuteTemplate(&processed, "bookmarks", projects)

	if err != nil {
		log.Fatalf("Could not execute bookmarks template: %s", err)
	}

	return processed.String()
}

// WriteBookmarkFile simply writes 'filename' to disk with content 'htmlContent'
func WriteBookmarkFile(filename string, htmlContent string) {
	outputPath := "./" + filename
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(htmlContent)
	w.Flush()
}
