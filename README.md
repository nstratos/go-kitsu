go-kitsu
========

go-kitsu is a Go client library for accessing the [kitsu.io API](http://docs.kitsu17.apiary.io).

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/nstratos/go-kitsu/kitsu?status.svg)](https://godoc.org/github.com/nstratos/go-kitsu/kitsu)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstratos/go-kitsu)](https://goreportcard.com/report/github.com/nstratos/go-kitsu)
[![Coverage Status](https://coveralls.io/repos/github/nstratos/go-kitsu/badge.svg?branch=jsonapi)](https://coveralls.io/github/nstratos/go-kitsu?branch=jsonapi)
[![Build Status](https://travis-ci.org/nstratos/go-kitsu.svg?branch=master)](https://travis-ci.org/nstratos/go-kitsu)

Installation
------------

This package can be installed using:

	go get github.com/nstratos/go-kitsu/kitsu

Usage
-----

Import the package using:

```go
import "github.com/nstratos/go-kitsu/kitsu"
```

Project Status
-------------

This project is currently under heavy development in
[jsonapi](https://github.com/nstratos/go-kitsu/tree/jsonapi) branch and
unstable.

A [x] means that the method works. The progress bar on working methods, shows
roughly the percentage of attributes and relationships implemented for that
resource.

[x] List Anime ![Progress](http://progressed.io/bar/65)
[x] Show Anime ![Progress](http://progressed.io/bar/65)
[ ] List Manga
[ ] Show Manga
[ ] List Drama
[ ] Show Drama
[x] List Users ![Progress](http://progressed.io/bar/75)
[ ] Create User
[x] Show User ![Progress](http://progressed.io/bar/75)
[ ] Update User
[ ] List LibraryEntries
[ ] Create LibraryEntry ![Progress](http://progressed.io/bar/50)
[x] Show LibraryEntry ![Progress](http://progressed.io/bar/80)
[ ] Update LibraryEntry
[ ] Delete LibraryEntry

License
-------

MIT
