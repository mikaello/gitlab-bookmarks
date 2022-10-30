package main

import (
	"flag"
	"log"
	"net/url"
	"strings"

	"github.io/gitlab-bookmarks/internal/bookmarks"
	"github.io/gitlab-bookmarks/internal/git"
)

var (
	token    *string
	baseurl  *string
	maxpages *int
	groups   groupFlags
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
	token = flag.String("token", "", "a token with API read permissions, not required, but only public repos without")
	baseurl = flag.String("baseurl", "https://gitlab.com", "the base url of your GitLab instance, including protocol scheme")
	maxpages = flag.Int("maxpages", 5, "the maximum number of pages to fetch, GitLab API is paginated")
	flag.Var(&groups, "group", "group to search for projects (use multiple flags for more groups), if not set all groups will be searched")
}

func main() {
	flag.Parse()

	_, err := url.ParseRequestURI(*baseurl)
	if err != nil {
		log.Fatalf("baseURL is invalid, '%s': %s", *baseurl, err)
	}

	// create a GitLab client
	client, err := git.Client(baseurl, *token)
	if err != nil {
		log.Fatalf("Error creating GitLab client: %s", err)
	}

	user, err := git.WhoAmI(client)
	if err != nil {
		log.Printf("You are not logged in to GitLab, only public repositories will be fetched")
	} else {
		log.Printf("You are using token of the user: %s", user.Username)
	}

	repos, err := git.FindAllRepositories(client, *maxpages, groups)
	if err != nil {
		log.Fatalf("Error fetching repositories: %s", err)
	}

	log.Printf("Total: Found %d repositories", len(repos))

	htmlContent := bookmarks.CreateBookmarkHTML(repos)
	bookmarks.WriteBookmarkFile("bookmarks.html", htmlContent)
}
