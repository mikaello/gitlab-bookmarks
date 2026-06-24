package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"github.com/mikaello/gitlab-bookmarks/internal/bookmarks"
	"github.com/mikaello/gitlab-bookmarks/internal/git"
)

var (
	token        *string
	baseurl      *string
	output       *string
	maxPages     *int
	folderBy     *string
	includeForks *bool
	showVersion  *bool
	groups       groupFlags
)

type groupFlags []string

func (i *groupFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *groupFlags) Set(group string) error {
	*i = append(*i, group)
	return nil
}

func init() {
	token = flag.String("token", "", "a token with API read permissions; without it only public repositories will be fetched")
	baseurl = flag.String("baseurl", "https://gitlab.com", "the base url of your GitLab instance, including protocol scheme")
	output = flag.String("output", "bookmarks.html", "path to the bookmarks file to write")
	maxPages = flag.Int("maxpages", 5, "the maximum number of pages to fetch, GitLab API is paginated")
	folderBy = flag.String("folderby", string(bookmarks.FolderModeFlat), "how to arrange projects in the bookmarks file: flat or namespace")
	includeForks = flag.Bool("includeforks", false, "if forks should be included (default is false)")
	showVersion = flag.Bool("version", false, "print version information and exit")
	flag.Var(&groups, "group", "group to search for projects (use multiple flags for more groups), if not set all groups will be searched")
}

// version is overridden via -ldflags="-X main.version=..." at release build
// time. When unset, fall back to the VCS information embedded by `go build`.
var version = ""

func versionString() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			return s.Value
		}
	}
	return "unknown"
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(versionString())
		os.Exit(0)
	}

	_, err := url.ParseRequestURI(*baseurl)
	if err != nil {
		log.Fatalf("baseURL is invalid, '%s': %s", *baseurl, err)
	}
	folderMode, err := bookmarks.ParseFolderMode(*folderBy)
	if err != nil {
		log.Fatal(err)
	}

	// create a GitLab client
	client, err := git.Client(*baseurl, *token)
	if err != nil {
		log.Fatalf("Error creating GitLab client: %s", err)
	}

	user, err := git.WhoAmI(client)
	if err != nil {
		log.Printf("You are not logged in to GitLab, only public repositories will be fetched")
	} else {
		log.Printf("You are using token of the user: %s", user.Username)
	}

	repos, err := git.FindAllRepositories(client, *maxPages, groups, *includeForks)
	if err != nil {
		log.Fatalf("Error fetching repositories: %s", err)
	}

	log.Printf("Total: Found %d repositories", len(repos))

	htmlContent, err := bookmarks.CreateBookmarkHTML(repos, folderMode)
	if err != nil {
		log.Fatalf("Error creating bookmark HTML: %s", err)
	}
	if err := bookmarks.WriteBookmarkFile(*output, htmlContent); err != nil {
		log.Fatalf("Error writing bookmark file: %s", err)
	}
}
