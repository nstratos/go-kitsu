# go-kitsu

go-kitsu is a Go client library for accessing the [kitsu.io API](http://docs.kitsu.apiary.io).

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/nstratos/go-kitsu/kitsu?status.svg)](https://godoc.org/github.com/nstratos/go-kitsu/kitsu)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstratos/go-kitsu)](https://goreportcard.com/report/github.com/nstratos/go-kitsu)
[![Coverage Status](https://coveralls.io/repos/github/nstratos/go-kitsu/badge.svg)](https://coveralls.io/github/nstratos/go-kitsu)
[![Build Status](https://travis-ci.org/nstratos/go-kitsu.svg?branch=master)](https://travis-ci.org/nstratos/go-kitsu)

## Installation

This package can be installed using:

	go get github.com/nstratos/go-kitsu/kitsu

## Usage

Import the package using:

```go
import "github.com/nstratos/go-kitsu/kitsu"
```


## Project Status

This project is currently under development. Below are all the currently
documented Kitsu API resources. The marked ones are currently implemented by
this package. The rest will be implemented over time. Contributions are
welcome.

### Characters & People

- [ ] Anime Characters
- [ ] Anime Productions
- [ ] Anime Staff
- [ ] Castings
- [ ] Characters
- [ ] Manga Characters
- [ ] Manga Staff
- [ ] People
- [ ] Producers

### Groups

- [ ] Group Categories
- [ ] Group Members
- [ ] Group Neighbors
- [ ] Group Permissions
- [ ] Groups

### Media

- [ ] Anime
  - [x] Show
  - [x] List
  - [ ] Create
  - [ ] Update
  - [ ] Delete
- [ ] Categories
- [ ] Category Favorites
- [ ] Chapters
- [ ] Drama
- [ ] Episodes
- [ ] Franchises
- [ ] Genres
- [ ] Installments
- [ ] Manga
- [ ] Mappings
- [ ] Media Follows
- [ ] Media Relationships
- [ ] Streamers
- [ ] Streaming Links

### Posts
- [ ] Comments
- [ ] Post Likes
- [ ] Post Follows
- [ ] Posts

### Reactions
- [ ] Media Reactions
- [ ] Review Likes
- [ ] Reviews

### Site Announcements
- [ ] Site Announcements

### User Libraries
- [ ] Library Entries
  - [x] Show
  - [x] List
  - [x] Create
  - [ ] Update
  - [x] Delete
- [ ] Library Entry Logs

### Users
- [ ] Favorites
- [ ] Follows
- [ ] Profile Link Sites
- [ ] Profile Links
- [ ] Roles
- [ ] Stats
- [ ] User Roles
- [ ] Users
  - [x] Show
  - [x] List
  - [ ] Create
  - [ ] Update
  - [ ] Delete

## Stability

The Kitsu API does not currently provide endpoint versioning. The only
available endpoint is the appropriately named "edge" endpoint
(https://kitsu.io/api/edge/) which "offers no guarantees: anything could change
at any time" (to quote the Kitsu API docs).

As a result, this package provides as many guarantees as the edge endpoint.
Nevertheless there is effort to keep the package as stable as possible through
integration tests that exercise the package against the live Kitsu API.

## Unit testing

To run all unit tests:

    go test

To see test coverage:

    go test -cover

For an HTML presentation of the coverage information:

    go test -coverprofile=cover.out && go tool cover -html=cover.out

And for heat maps:

    go test -coverprofile=cover.out -covermode=count && go tool cover -html=cover.out

## Integration testing

The integration tests will exercise the package against the live Kitsu API and
will hopefully reveal incompatible changes. Since these tests are using live
data, they take much longer to run and there is a chance for false positives.

The tests need a dedicated test account for authentication. A newly created
account has no slug set by default. To set the slug, open the profile settings
and set Profile URL.

To run the integration tests:

    go test -tags=integration -slug="<test account slug>" -password="<test account password>"

## License

MIT
