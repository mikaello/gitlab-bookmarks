module github.io/gitlab-bookmarks

go 1.25.0

require gitlab.com/gitlab-org/api/client-go/v2 v2.13.0

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/time v0.14.0 // indirect
)

replace github.io/gitlab-bookmarks/internal/git => ./internal/git

replace github.io/gitlab-bookmarks/internal/bookmarks => ./internal/bookmarks
