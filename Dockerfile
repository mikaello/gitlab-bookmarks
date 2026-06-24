# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /gitlab-bookmarks \
    .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /gitlab-bookmarks /usr/local/bin/gitlab-bookmarks

WORKDIR /output
VOLUME ["/output"]

ENTRYPOINT ["/usr/local/bin/gitlab-bookmarks"]
