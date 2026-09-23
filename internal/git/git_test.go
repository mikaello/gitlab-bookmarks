package git

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *gitlab.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := Client(server.URL, "test-token")
	if err != nil {
		t.Fatalf("Client returned error: %s", err)
	}
	return client
}

func writeJSON(t *testing.T, writer http.ResponseWriter, body string) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(writer, body); err != nil {
		t.Fatalf("write response: %s", err)
	}
}

func setPaginationHeaders(writer http.ResponseWriter, page, next, totalPages string) {
	writer.Header().Set("X-Page", page)
	writer.Header().Set("X-Next-Page", next)
	writer.Header().Set("X-Per-Page", "100")
	writer.Header().Set("X-Total-Pages", totalPages)
}

func TestWhoAmI(t *testing.T) {
	client := newTestClient(t, func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v4/user" {
			t.Errorf("path = %q, want /api/v4/user", request.URL.Path)
		}
		if got := request.Header.Get("Private-Token"); got != "test-token" {
			t.Errorf("Private-Token = %q, want test-token", got)
		}
		writeJSON(t, writer, `{"id":1,"username":"mikael"}`)
	})

	user, err := WhoAmI(client)
	if err != nil {
		t.Fatalf("WhoAmI returned error: %s", err)
	}
	if user.Username != "mikael" {
		t.Errorf("username = %q, want mikael", user.Username)
	}
}

func TestFindAllRepositoriesPaginatesFiltersAndSorts(t *testing.T) {
	var requestedPages []string
	client := newTestClient(t, func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v4/projects" {
			t.Errorf("path = %q, want /api/v4/projects", request.URL.Path)
		}
		if got := request.URL.Query().Get("per_page"); got != "100" {
			t.Errorf("per_page = %q, want 100", got)
		}

		page := request.URL.Query().Get("page")
		requestedPages = append(requestedPages, page)
		switch page {
		case "1":
			setPaginationHeaders(writer, "1", "2", "2")
			writeJSON(t, writer, `[
				{"id":2,"name":"Zulu","path_with_namespace":"group/zulu"},
				{"id":3,"name":"Fork","path_with_namespace":"group/fork","forked_from_project":{"id":99}}
			]`)
		case "2":
			setPaginationHeaders(writer, "2", "", "2")
			writeJSON(t, writer, `[{"id":1,"name":"Alpha","path_with_namespace":"group/alpha"}]`)
		default:
			t.Errorf("unexpected page %q", page)
		}
	})

	projects, err := FindAllRepositories(client, 0, nil, false)
	if err != nil {
		t.Fatalf("FindAllRepositories returned error: %s", err)
	}
	if !reflect.DeepEqual(requestedPages, []string{"1", "2"}) {
		t.Errorf("requested pages = %v, want [1 2]", requestedPages)
	}
	if got := projectIDs(projects); !reflect.DeepEqual(got, []int64{1, 2}) {
		t.Errorf("project IDs = %v, want [1 2]", got)
	}
}

func TestFindAllRepositoriesHonorsPageLimit(t *testing.T) {
	requests := 0
	client := newTestClient(t, func(writer http.ResponseWriter, _ *http.Request) {
		requests++
		setPaginationHeaders(writer, "1", "2", "2")
		writeJSON(t, writer, `[{"id":1,"name":"Alpha"}]`)
	})

	projects, err := FindAllRepositories(client, 1, nil, true)
	if err != nil {
		t.Fatalf("FindAllRepositories returned error: %s", err)
	}
	if requests != 1 {
		t.Errorf("requests = %d, want 1", requests)
	}
	if got := projectIDs(projects); !reflect.DeepEqual(got, []int64{1}) {
		t.Errorf("project IDs = %v, want [1]", got)
	}
}

func TestFindAllRepositoriesForGroupsIncludesSubgroupsAndDeduplicates(t *testing.T) {
	requestedGroups := make(map[string]bool)
	client := newTestClient(t, groupHandler(t, requestedGroups))
	projects, err := FindAllRepositories(client, 0, []string{"1", "2"}, true)
	if err != nil {
		t.Fatalf("FindAllRepositories returned error: %s", err)
	}
	if !reflect.DeepEqual(requestedGroups, map[string]bool{"1": true, "2": true}) {
		t.Errorf("requested groups = %v, want groups 1 and 2", requestedGroups)
	}
	if got := projectIDs(projects); !reflect.DeepEqual(got, []int64{1, 2}) {
		t.Errorf("project IDs = %v, want [1 2]", got)
	}
}

func groupHandler(t *testing.T, requestedGroups map[string]bool) http.HandlerFunc {
	t.Helper()
	return func(writer http.ResponseWriter, request *http.Request) {
		if got := request.URL.Query().Get("include_subgroups"); got != "true" {
			t.Errorf("include_subgroups = %q, want true", got)
		}

		group := ""
		switch request.URL.Path {
		case "/api/v4/groups/1/projects":
			group = "1"
		case "/api/v4/groups/2/projects":
			group = "2"
		default:
			t.Errorf("unexpected path %q", request.URL.Path)
		}
		requestedGroups[group] = true
		setPaginationHeaders(writer, "1", "", "1")
		if group == "1" {
			writeJSON(t, writer, `[{"id":2,"name":"Zulu","path_with_namespace":"group/zulu"}]`)
			return
		}
		writeJSON(t, writer, `[
			{"id":2,"name":"Zulu","path_with_namespace":"group/zulu"},
			{"id":1,"name":"Alpha","path_with_namespace":"group/alpha"}
		]`)
	}
}

func TestFindAllRepositoriesReturnsAPIError(t *testing.T) {
	client := newTestClient(t, func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, `{"message":"bad request"}`, http.StatusBadRequest)
	})

	if _, err := FindAllRepositories(client, 0, nil, true); err == nil {
		t.Fatal("FindAllRepositories returned nil error, want API error")
	}
}

func projectIDs(projects []*gitlab.Project) []int64 {
	ids := make([]int64, len(projects))
	for index, project := range projects {
		ids[index] = project.ID
	}
	return ids
}

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
