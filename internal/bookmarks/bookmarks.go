package bookmarks

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	"gitlab.com/gitlab-org/api/client-go/v2"
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
	if err := templates.ExecuteTemplate(&processed, "bookmarks", bookmarkData(projects, folderMode)); err != nil {
		return "", fmt.Errorf("execute bookmarks template: %w", err)
	}

	return processed.String(), nil
}

func bookmarkData(projects []*gitlab.Project, folderMode FolderMode) bookmarkPage {
	root := bookmarkFolder{
		Name:    rootFolderName,
		AddDate: 1658268705,
	}

	if folderMode == FolderModeNamespace {
		root.Folders = namespaceFolders(projects)
	} else {
		root.Projects = bookmarkProjects(projects)
	}

	return bookmarkPage{Root: root}
}

func namespaceFolders(projects []*gitlab.Project) []bookmarkFolder {
	foldersByName := make(map[string]int)
	var folders []bookmarkFolder

	for _, project := range projects {
		name := namespaceName(project)
		index, ok := foldersByName[name]
		if !ok {
			index = len(folders)
			foldersByName[name] = index
			folders = append(folders, bookmarkFolder{Name: name, AddDate: 1658268705})
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

// WriteBookmarkFile writes 'filename' to disk with content 'htmlContent'.
func WriteBookmarkFile(filename string, htmlContent string) error {
	if err := os.WriteFile(filename, []byte(htmlContent), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filename, err)
	}
	return nil
}
