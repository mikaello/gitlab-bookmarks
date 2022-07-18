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
func FindAllRepositories(c *gitlab.Client, maxPages int) ([]*gitlab.Project, error) {
	listOptions := gitlab.ListOptions{ PerPage: 100, }
	projects, response, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{
		ListOptions: listOptions,
	})

	for {
		if err != nil {
			return nil, err
		}

		log.Printf("Page %d of %d (but max %d)", response.CurrentPage, response.TotalPages, maxPages)

		if len(projects) == 0 {
			break
		} else if response.NextPage > maxPages {
			break
		}

		proj := []*gitlab.Project{}
		proj, response, err = c.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    response.NextPage,
				PerPage: 100,
			},
			Archived: gitlab.Bool(true),
			//OrderBy:  gitlab.OrderByID,
			//Sort:     gitlab.SortAsc
		})
		projects = append(projects, proj...)
	}


	if err != nil {
		return projects, err
	}

	return projects, err
}
