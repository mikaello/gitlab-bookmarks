package main

import (
	"flag"
	"log"

	"github.io/gitlab-bookmarks/internal/git"
)

var (
	mode    *string
	pid     *int
	token   *string
	name    *string
	baseurl *string
)

func init() {

	mode = flag.String("mode", "empty", "empty\n  deletes empty registries.\nby-name\n  deletes registries with a given name")
	pid = flag.Int("pid", -1, "the GitLab project id")
	token = flag.String("token", "", "a token with permissions in the project")
	name = flag.String("name", "", "the name of a repo to delete when mode=by_name")
	baseurl = flag.String("baseurl", "https://gitlab.com", "the base url of your gitlab instance, including protocol scheme")
}

func main() {

	flag.Parse()

	var r []int
	var err error
	var p int = *pid

	// create a GitLab client
	c, err := git.Client(baseurl, *token)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	switch {
	case *mode == "empty":
		r, err = git.FindEmptyRegistries(&c, p)
	case *mode == "name":
		r, err = git.FindNamedRegistry(&c, p, *name)
	}
	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		err := git.DeleteRegistries(&c, p, r)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
