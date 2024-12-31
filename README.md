# go-kitsu

go-kitsu is a Go client library for accessing the [kitsu.io API](http://docs.kitsu.apiary.io).

[![Go Reference](https://pkg.go.dev/badge/github.com/nstratos/go-kitsu/kitsu.svg)](https://pkg.go.dev/github.com/nstratos/go-kitsu/kitsu)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstratos/go-kitsu)](https://goreportcard.com/report/github.com/nstratos/go-kitsu)
[![Coverage Status](https://coveralls.io/repos/github/nstratos/go-kitsu/badge.svg)](https://coveralls.io/github/nstratos/go-kitsu)
[![Test Status](https://github.com/nstratos/go-kitsu/workflows/tests/badge.svg)](https://github.com/nstratos/go-kitsu/actions?query=workflow%3Atests)
[![Integration Status](https://github.com/nstratos/go-kitsu/workflows/integration/badge.svg)](https://github.com/nstratos/go-kitsu/actions?query=workflow%3Aintegration)

## Installation

This package can be installed using:

	go get github.com/nstratos/go-kitsu/kitsu

## Usage

Import the package using:

```go
import "github.com/nstratos/go-kitsu/kitsu"
```


## Project Status

This project is currently under development. Expect things to change. Some
useful methods like getting users, library entries and anime are already
implemented. For a full list of what needs to be implemented please check
[Roadmap.md](Roadmap.md).


## Endpoint Stability

The Kitsu API does not currently provide endpoint versioning. The only
available endpoint is the appropriately named "edge" endpoint
(https://kitsu.io/api/edge/) which (to quote the Kitsu API docs) "offers no
guarantees: anything could change at any time".

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
