# gitlab-bookmarks

[![Go CI](https://github.com/mikaello/gitlab-bookmarks/actions/workflows/verify-go-build.yml/badge.svg)](https://github.com/mikaello/gitlab-bookmarks/actions/workflows/verify-go-build.yml)
[![Latest release](https://img.shields.io/github/v/release/mikaello/gitlab-bookmarks)](https://github.com/mikaello/gitlab-bookmarks/releases/latest)
[![License](https://img.shields.io/github/license/mikaello/gitlab-bookmarks)](LICENSE)

Tool to fetch projects from GitLab to provision your bookmarks.

Run this tool with the GitLab instance of your choice, it will use the API and
gather all projects into a file called `bookmarks.html`.
This file follows the bookmarks format expected by browsers for importing.
So open your favourite browser and import the generated `bookmarks.html` file, you can now use the address bar to search for projects of the chosen GitLab instance.

## Requirements

Go 1.25 or newer (only required for building from source).

## Usage

```
$ ./gitlab-bookmarks --help
Usage of ./gitlab-bookmarks:
  -baseurl string
        the base url of your GitLab instance, including protocol scheme (default "https://gitlab.com")
  -group value
        group to search for projects (use multiple flags for more groups), if not set all groups will be searched
  -includeforks
        if forks should be included (default is false)
  -maxpages int
        the maximum number of pages to fetch, GitLab API is paginated (default 5, 100 per page)
  -token string
        a token with API read permissions; without it only public repositories will be fetched
```

Example usage:

```shell
./gitlab-bookmarks -baseurl https://mycompany.gitlab.com -group some-group -group another-group -maxpages 100
```

### Creating a token

To access private projects you need a [personal access token](https://docs.gitlab.com/user/profile/personal_access_tokens/) with the `read_api` scope.
Group access tokens or project access tokens with the same scope also work.

### Importing the generated file

After running the tool you will have a `bookmarks.html` in the current directory.

- **Firefox**: `Bookmarks → Manage bookmarks → Import and Backup → Import bookmarks from HTML…`
- **Chrome / Edge**: `Bookmark manager → ⋮ → Import bookmarks`
- **Safari**: `File → Import From → Bookmarks HTML File…`

## Install

Either download a prebuilt binary from
[the release page](https://github.com/mikaello/gitlab-bookmarks/releases) (if it
exists for your system), or build with go:

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
