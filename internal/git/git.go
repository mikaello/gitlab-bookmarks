package git

import (
	"log"

	"github.com/xanzy/go-gitlab"
)

// Client creates a GitLab client for the given base URL and token.
func Client(baseurl *string, token string) (gitlab.Client, error) {
	url := *baseurl
	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(url+"/api/v4"))
	return *c, err
}

// WhoAmI returns the user that is logged in to GitLab.
func WhoAmI(c *gitlab.Client) (*gitlab.User, error) {
	user, _, err := c.Users.CurrentUser()
	return user, err
}

// FindAllRepositories returns all repositories that the user has access to,
// up to the given maximum number of pages.
// See also pagination example https://github.com/xanzy/go-gitlab/blob/master/examples/pagination.go
func FindAllRepositories(c *gitlab.Client, maxPages int) ([]*gitlab.Project, error) {
	options := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		//Archived: gitlab.Bool(true),
		//OrderBy:  gitlab.OrderByID,
		//Sort:     gitlab.SortAsc
	}

	var totalProjects []*gitlab.Project
	var err error

	for {
		projects, response, err := c.Projects.ListProjects(options)
		if err != nil {
			return nil, err
		}

		if len(projects) == 0 {
			break
		}

		log.Printf("Page %d of %d (but max %d)", response.CurrentPage, response.TotalPages, maxPages)

		totalProjects = append(totalProjects, projects...)
		options.Page = response.NextPage

		if response.CurrentPage == maxPages {
			break
		} else if response.NextPage == 0 {
			break
		}
	}

	return totalProjects, err
}
