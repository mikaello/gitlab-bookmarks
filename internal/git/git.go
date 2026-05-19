package git

import (
	"log"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

const perPage = 100

// Client creates a GitLab client for the given base URL and token.
func Client(baseurl string, token string) (*gitlab.Client, error) {
	return gitlab.NewClient(token, gitlab.WithBaseURL(baseurl+"/api/v4"))
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

// pageFetcher returns the projects for a single page along with the response
// metadata used to drive pagination.
type pageFetcher func(page int64) ([]*gitlab.Project, *gitlab.Response, error)

// paginate walks pages of projects via fetch, stopping when the API runs out
// of pages or maxPages is reached (maxPages <= 0 means no cap).
func paginate(maxPages int, fetch pageFetcher) ([]*gitlab.Project, error) {
	var all []*gitlab.Project
	var page int64 = 1
	for {
		projects, response, err := fetch(page)
		if err != nil {
			return nil, err
		}
		if len(projects) == 0 {
			break
		}

		log.Printf("Page %d of %d (but max %d)", response.CurrentPage, response.TotalPages, maxPages)
		all = append(all, projects...)

		if maxPages > 0 && response.CurrentPage >= int64(maxPages) {
			break
		}
		if response.NextPage == 0 {
			break
		}
		page = response.NextPage
	}
	return all, nil
}

func findAllProjects(c *gitlab.Client, maxPages int, includeForks bool) ([]*gitlab.Project, error) {
	projects, err := paginate(maxPages, func(page int64) ([]*gitlab.Project, *gitlab.Response, error) {
		return c.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: perPage},
		})
	})
	if err != nil {
		return nil, err
	}

	if !includeForks {
		projects = excludeForks(projects)
	}
	return projects, nil
}

func findAllProjectsForGroups(c *gitlab.Client, maxPages int, groups []string, includeForks bool) ([]*gitlab.Project, error) {
	var all []*gitlab.Project

	for _, groupID := range groups {
		log.Printf("- Fetching projects for group %s", groupID)

		groupProjects, err := paginate(maxPages, func(page int64) ([]*gitlab.Project, *gitlab.Response, error) {
			return c.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{
				IncludeSubGroups: gitlab.Ptr(true),
				ListOptions:      gitlab.ListOptions{Page: page, PerPage: perPage},
			})
		})
		if err != nil {
			return nil, err
		}

		log.Printf("Found %d projects for group %s", len(groupProjects), groupID)
		all = append(all, groupProjects...)
	}

	if !includeForks {
		all = excludeForks(all)
	}
	return all, nil
}

func excludeForks(projects []*gitlab.Project) []*gitlab.Project {
	var nonForks []*gitlab.Project
	for _, project := range projects {
		if project.ForkedFromProject == nil {
			nonForks = append(nonForks, project)
		}
	}
	log.Printf("Excluded %d forked projects", len(projects)-len(nonForks))
	return nonForks
}
