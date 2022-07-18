package main

import (
	"flag"
	"log"

	"github.io/gitlab-bookmarks/internal/git"
	"github.io/gitlab-bookmarks/internal/bookmarks"
)

var (
	token    *string
	baseurl  *string
	maxpages *int
)

func init() {
	token = flag.String("token", "", "a token with API read permissions")
	baseurl = flag.String("baseurl", "https://gitlab.com", "the base url of your GitLab instance, including protocol scheme")
	maxpages = flag.Int("maxpages", 2, "the maximum number of pages to fetch")
}

func main() {
	flag.Parse()

	// create a GitLab client
	client, err := git.Client(baseurl, *token)
	if err != nil {
		log.Printf("Error creating GitLab client: %s", err)
	}

	user, err := git.WhoAmI(&client)
	if err != nil {
		log.Printf("You are not logged in to GitLab")
	} else {
		log.Printf("You are using token of the user: %s", user.Username)
	}

	repos, err := git.FindAllRepositories(&client, *maxpages)
	if err != nil {
		log.Printf("Error fetching repositories: %s", err)
	}

	log.Printf("Found %d number of repositories", len(repos))

	bookmarks.CreateBookmarkFile(repos)
}
