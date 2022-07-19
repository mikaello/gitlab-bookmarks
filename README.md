# gitlab-bookmarks

Tool to fetch projects from GitLab to provision your bookmarks.

## Usage:

```
$ ./gitlab-bookmarks --help
Usage of ./gitlab-bookmarks:
  -baseurl string
        the base url of your GitLab instance, including protocol scheme (default "https://gitlab.com")
  -maxpages int
        the maximum number of pages to fetch (default 2, 100 per page)
  -token string
        a token with API read permissions (not necessary, but only public projects without)
```

## Install

Either download a prebuilt binary from
[the release page](https://github.com/mikaello/gitlab-bookmarks/releases) (if it
exist for your system), or build with go:

```shell
go build -o gitlab-bookmarks cmd/provision/main.go
```
