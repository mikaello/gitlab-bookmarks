package git

import (
	"testing"

	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

func TestDeduplicateAndSortProjects(t *testing.T) {
	projects := []*gitlab.Project{
		{ID: 2, Name: "Zulu", PathWithNamespace: "group/zulu"},
		{ID: 1, Name: "Alpha", PathWithNamespace: "group/alpha"},
		{ID: 2, Name: "Zulu duplicate", PathWithNamespace: "group/zulu"},
	}

	projects = deduplicateProjects(projects)
	sortProjects(projects)

	if len(projects) != 2 {
		t.Fatalf("got %d projects, want 2", len(projects))
	}
	if projects[0].ID != 1 || projects[1].ID != 2 {
		t.Errorf("project IDs = [%d, %d], want [1, 2]", projects[0].ID, projects[1].ID)
	}
}
