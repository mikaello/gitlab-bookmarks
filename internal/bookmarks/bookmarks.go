package bookmarks

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/gitlab-org/api/client-go/v3"
)

//go:embed bookmarks.tmpl
var bookmarksTmpl embed.FS

const rootFolderName = "GitLab projects"

// FolderMode controls how projects are arranged in the generated bookmarks.
type FolderMode string

const (
	FolderModeFlat      FolderMode = "flat"
	FolderModeNamespace FolderMode = "namespace"
)

type bookmarkPage struct {
	Root bookmarkFolder
}

type bookmarkFolder struct {
	Name     string
	AddDate  int64
	Projects []bookmarkProject
	Folders  []bookmarkFolder
}

type bookmarkProject struct {
	Name    string
	WebURL  string
	AddDate int64
}

// ParseFolderMode validates the folder mode used for generated bookmarks.
func ParseFolderMode(value string) (FolderMode, error) {
	switch FolderMode(value) {
	case FolderModeFlat, FolderModeNamespace:
		return FolderMode(value), nil
	default:
		return "", fmt.Errorf("invalid folder mode %q: expected %q or %q", value, FolderModeFlat, FolderModeNamespace)
	}
}

// CreateBookmarkHTML creates a bookmark content for the given repositories.
func CreateBookmarkHTML(projects []*gitlab.Project, folderMode FolderMode) (string, error) {
	templates := template.Must(template.New("").ParseFS(bookmarksTmpl, "bookmarks.tmpl"))

	var processed bytes.Buffer
	if err := templates.ExecuteTemplate(&processed, "bookmarks", bookmarkData(projects, folderMode, time.Now().Unix())); err != nil {
		return "", fmt.Errorf("execute bookmarks template: %w", err)
	}

	return processed.String(), nil
}

func bookmarkData(projects []*gitlab.Project, folderMode FolderMode, addDate int64) bookmarkPage {
	root := bookmarkFolder{
		Name:    rootFolderName,
		AddDate: addDate,
	}

	if folderMode == FolderModeNamespace {
		root.Folders = namespaceFolders(projects, addDate)
	} else {
		root.Projects = bookmarkProjects(projects)
	}

	return bookmarkPage{Root: root}
}

func namespaceFolders(projects []*gitlab.Project, addDate int64) []bookmarkFolder {
	foldersByName := make(map[string]int)
	var folders []bookmarkFolder

	for _, project := range projects {
		name := namespaceName(project)
		index, ok := foldersByName[name]
		if !ok {
			index = len(folders)
			foldersByName[name] = index
			folders = append(folders, bookmarkFolder{Name: name, AddDate: addDate})
		}
		folders[index].Projects = append(folders[index].Projects, bookmarkProjectFor(project))
	}

	return folders
}

func bookmarkProjects(projects []*gitlab.Project) []bookmarkProject {
	bookmarks := make([]bookmarkProject, 0, len(projects))
	for _, project := range projects {
		bookmarks = append(bookmarks, bookmarkProjectFor(project))
	}
	return bookmarks
}

func bookmarkProjectFor(project *gitlab.Project) bookmarkProject {
	bookmark := bookmarkProject{
		Name:   project.Name,
		WebURL: project.WebURL,
	}
	if project.LastActivityAt != nil {
		bookmark.AddDate = project.LastActivityAt.Unix()
	}
	return bookmark
}

func namespaceName(project *gitlab.Project) string {
	if project.Namespace != nil {
		if project.Namespace.FullPath != "" {
			return project.Namespace.FullPath
		}
		if project.Namespace.Name != "" {
			return project.Namespace.Name
		}
	}
	if project.PathWithNamespace != "" && project.Path != "" {
		return strings.TrimSuffix(project.PathWithNamespace, "/"+project.Path)
	}
	return "Uncategorized"
}

// WriteBookmarkFile writes htmlContent to filename, or to stdout when filename is "-".
func WriteBookmarkFile(filename string, htmlContent string) error {
	return writeBookmarkFile(filename, htmlContent, os.Stdout)
}

func writeBookmarkFile(filename string, htmlContent string, stdout io.Writer) error {
	if filename == "-" {
		if _, err := io.WriteString(stdout, htmlContent); err != nil {
			return fmt.Errorf("write stdout: %w", err)
		}
		return nil
	}

	temporary, err := os.CreateTemp(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", filename, err)
	}
	temporaryName := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryName)
	}()

	if err := temporary.Chmod(0o644); err != nil {
		return fmt.Errorf("set permissions on temporary file for %s: %w", filename, err)
	}
	if _, err := io.WriteString(temporary, htmlContent); err != nil {
		return fmt.Errorf("write temporary file for %s: %w", filename, err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary file for %s: %w", filename, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file for %s: %w", filename, err)
	}
	if err := os.Rename(temporaryName, filename); err != nil {
		return fmt.Errorf("replace %s: %w", filename, err)
	}
	return nil
}
