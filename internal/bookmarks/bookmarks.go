package bookmarks

import (
	"bufio"
	"bytes"
	"html/template"
	"os"

	"github.com/xanzy/go-gitlab"
)

// CreateBookmarkFile creates a bookmark file for the given repositories.
func CreateBookmarkFile(projects []*gitlab.Project) {
	allFiles := []string{"bookmarks.tmpl"}

	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./templates/"+tmpl)
	}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	templates.ExecuteTemplate(&processed, "bookmarks", projects)

	outputPath := "./bookmarks.html"
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(processed.String())
	w.Flush()
}
