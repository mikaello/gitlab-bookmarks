module github.io/gitlab-bookmarks

go 1.24.0

require gitlab.com/gitlab-org/api/client-go v1.41.0

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	golang.org/x/exp v0.0.0-20250813145105-42675adae3e6 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.io/gitlab-bookmarks/internal/git => ./internal/git

replace github.io/gitlab-bookmarks/internal/bookmarks => ./internal/bookmarks
