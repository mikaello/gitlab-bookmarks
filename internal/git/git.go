package git

import (
	"log"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Client creates a GitLab client for the given base URL and token.
func Client(baseurl *string, token string) (*gitlab.Client, error) {
	url := *baseurl
	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(url+"/api/v4"))
	return c, err
}

// WhoAmI returns the user that is logged in to GitLab.
func WhoAmI(c *gitlab.Client) (*gitlab.User, error) {
	user, _, err := c.Users.CurrentUser()
	return user, err
}

// FindAllRepositories returns all repositories that the user has access to,
// up to the given maximum number of pages.
// See also pagination example https://github.com/xanzy/go-gitlab/blob/master/examples/pagination.go
func FindAllRepositories(c *gitlab.Client, maxPages int, groups []string, includeForks bool) ([]*gitlab.Project, error) {
	if len(groups) > 0 {
		return findAllProjectsForGroups(c, maxPages, groups, includeForks)
	}
	return findAllProjects(c, maxPages, includeForks)
}

func findAllProjects(c *gitlab.Client, maxPages int, includeForks bool) ([]*gitlab.Project, error) {
	// exclude forks if includeForks is false
	options := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		//Archived: gitlab.Ptr(true),
		//OrderBy:  gitlab.OrderByID,
		//Sort:     gitlab.SortAsc
	}

	var (
		totalProjects []*gitlab.Project
		err           error
	)

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

	if !includeForks {
		totalProjects = excludeForks(totalProjects)
	}

	return totalProjects, err
}

func findAllProjectsForGroups(c *gitlab.Client, maxPages int, groups []string, includeForks bool) ([]*gitlab.Project, error) {
	options := &gitlab.ListGroupProjectsOptions{
		IncludeSubGroups: gitlab.Ptr(true),
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var (
		projects []*gitlab.Project
		err      error
	)

	for _, groupID := range groups {
		log.Printf("- Fetching projects for group %s", groupID)
		var groupProjects []*gitlab.Project

		for {
			tempGroupProjects, response, err := c.Groups.ListGroupProjects(groupID, options)

			if err != nil {
				return nil, err
			}

			if len(tempGroupProjects) == 0 {
				break
			}

			log.Printf("Page %d of %d (but max %d)", response.CurrentPage, response.TotalPages, maxPages)

			groupProjects = append(groupProjects, tempGroupProjects...)
			options.Page = response.NextPage

			if response.CurrentPage == maxPages {
				break
			} else if response.NextPage == 0 {
				break
			}
		}

		log.Printf("Found %d projects for group %s", len(groupProjects), groupID)
		projects = append(projects, groupProjects...)
	}

	if !includeForks {
		projects = excludeForks(projects)
	}

	return projects, err
}

func excludeForks(projects []*gitlab.Project) []*gitlab.Project {
	var nonForks []*gitlab.Project
	for _, project := range projects {
		if project.ForkedFromProject == nil {
			nonForks = append(nonForks, project)
		}
	}
	return nonForks
}
