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
  -folderby string
        how to arrange projects in the bookmarks file: flat or namespace (default "flat")
  -group value
        group to search for projects (use multiple flags for more groups), if not set all groups will be searched
  -includeforks
        if forks should be included (default is false)
  -maxpages int
        the maximum number of pages to fetch, GitLab API is paginated (default 5, 100 per page)
  -output string
        path to the bookmarks file to write (default "bookmarks.html")
  -token string
        a token with API read permissions; without it only public repositories will be fetched
  -version
        print version information and exit
```

Example usage:

```shell
./gitlab-bookmarks -baseurl https://mycompany.gitlab.com -group some-group -group another-group -maxpages 100
```

Group projects into namespace folders:

```shell
./gitlab-bookmarks -baseurl https://mycompany.gitlab.com -folderby namespace
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

Install the latest release with Go:

```shell
go install github.com/mikaello/gitlab-bookmarks@latest
```

To install a local checkout instead:

```shell
go install .
```

Make sure the Go binary directory is on your `PATH`.
You can find it with `go env GOBIN`, or use `$(go env GOPATH)/bin` when `GOBIN` is empty.

Prebuilt binaries are also available on the [release page](https://github.com/mikaello/gitlab-bookmarks/releases).

## Development

Compile and run by running:

```shell
go run .
```

With params:

```shell
go run . -baseurl https://your-gitlab.com
```
