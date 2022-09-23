package main

import (
	"flag"
	"log"
	"net/url"

	"github.io/gitlab-bookmarks/internal/bookmarks"
	"github.io/gitlab-bookmarks/internal/git"
)

var (
	token    *string
	baseurl  *string
	maxpages *int
)

func init() {
	token = flag.String("token", "", "a token with API read permissions, not required, but only public repos without")
	baseurl = flag.String("baseurl", "https://gitlab.com", "the base url of your GitLab instance, including protocol scheme")
	maxpages = flag.Int("maxpages", 2, "the maximum number of pages to fetch, GitLab API is paginated")
}

func main() {
	flag.Parse()

	_, err := url.ParseRequestURI(*baseurl)
	if err != nil {
		log.Fatalf("baseURL '%s' is invalid, %s", *baseurl, err)
	}

	// create a GitLab client
	client, err := git.Client(baseurl, *token)
	if err != nil {
		log.Fatalf("Error creating GitLab client: %s", err)
	}

	user, err := git.WhoAmI(&client)
	if err != nil {
		log.Printf("You are not logged in to GitLab, will only fetch public repositories")
	} else {
		log.Printf("You are using token of the user: %s", user.Username)
	}

	repos, err := git.FindAllRepositories(&client, *maxpages)
	if err != nil {
		log.Fatalf("Error fetching repositories: %s", err)
	}

	log.Printf("Found %d repositories", len(repos))

	htmlContent := bookmarks.CreateBookmarkHtml(repos)
	bookmarks.WriteBookmarkFile("bookmarks.html", htmlContent)
}
