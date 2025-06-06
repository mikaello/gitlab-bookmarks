module github.io/gitlab-bookmarks

go 1.23.4

require gitlab.com/gitlab-org/api/client-go v0.129.0

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)

replace github.io/gitlab-bookmarks/internal/git => ./internal/git

replace github.io/gitlab-bookmarks/internal/bookmarks => ./internal/bookmarks
