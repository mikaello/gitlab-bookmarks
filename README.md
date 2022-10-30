# gitlab-bookmarks

Tool to fetch projects from GitLab to provision your bookmarks.

Run this tool with the GitLab instance of your choice, it will use the API and
gather all projects into a file called `bookmarks.html`. This file follows the
bookmarks format expected by browsers for importing. So open your favourite
browser and import the generated `bookmarks.html` file, you can now use the
adress bar to search for projects of the chosen GitLab instance.

## Usage:

```
$ ./gitlab-bookmarks --help
Usage of ./gitlab-bookmarks:
  -baseurl string
        the base url of your GitLab instance, including protocol scheme (default "https://gitlab.com")
  -group value
        group to search for projects (use multiple flags for more groups), if not set all groups will be searched
  -maxpages int
        the maximum number of pages to fetch, GitLab API is paginated (default 5, 100 per page)
  -token string
        a token with API read permissions, not required, but only public repos without
```

Example usage:

```
./gitlab-bookmarks -baseurl https://mycompany.gitlab.com -group some-group -group another-group -maxpages 100
```

## Install

Either download a prebuilt binary from
[the release page](https://github.com/mikaello/gitlab-bookmarks/releases) (if it
exist for your system), or build with go:

```shell
go build -o gitlab-bookmarks cmd/provision/main.go
```

## Development

Compile and run by running:

```shell
go run cmd/provision/main.go
```

With params:

```shell
go run cmd/provision/main.go -baseurl https://your-gitlab.com
```
