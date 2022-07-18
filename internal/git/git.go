package git

import (
	"log"

	"github.com/xanzy/go-gitlab"
)

// Create a GitLab client
func Client(baseurl *string, token string) (gitlab.Client, error) {
	url := *baseurl
	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(url+"/api/v4"))
	return *c, err
}

func FindAllRepositories(c *gitlab.Client) ([]*gitlab.Project, error) {
	projects, _, err := c.Projects.ListProjects(nil)

	if err != nil {
		return projects, err
	}

	return projects, err

}

// Find a Container Registry by name
func FindNamedRegistry(c *gitlab.Client, pid int, name string) (reg []int, err error) {
	projectRepositories, _, err := c.ContainerRegistry.ListProjectRegistryRepositories(pid, &gitlab.ListRegistryRepositoriesOptions{})
	if err != nil {
		return reg, err
	}
	for _, r := range projectRepositories {
		if r.Name == name {
			reg = append(reg, r.ID)
		}
	}
	return reg, err
}

// Find all Container Registry Repositories without any tags
func FindEmptyRegistries(c *gitlab.Client, pid int) (reg []int, err error) {
	projectRepositories, _, err := c.ContainerRegistry.ListProjectRegistryRepositories(pid, &gitlab.ListRegistryRepositoriesOptions{})
	if err != nil {
		return reg, err
	}
	for _, r := range projectRepositories {
		t, _, err := c.ContainerRegistry.ListRegistryRepositoryTags(pid, r.ID, &gitlab.ListRegistryRepositoryTagsOptions{})
		if err != nil {
			return reg, err
		}
		if len(t) == 0 {
			reg = append(reg, r.ID)
		}
	}
	return reg, err
}

func DeleteRegistries(c *gitlab.Client, pid int, r []int) (err error) {
	log.Print(r)
	for i := range r {
		log.Printf("Deleting repository with ID %d from project %v", r[i], pid)
		c.ContainerRegistry.DeleteRegistryRepository(pid, r[i])
	}
	return err
}
