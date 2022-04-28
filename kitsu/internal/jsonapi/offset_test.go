package jsonapi

import (
	"reflect"
	"testing"

	"github.com/google/jsonapi"
)

func Test_parseOffset(t *testing.T) {
	links := jsonapi.Links{
		"first": "http://somesite.com/movies?page[limit]=50&page[offset]=50",
		"prev":  "http://somesite.com/movies?page[limit]=50&page[offset]=0",
		"next":  "http://somesite.com/movies?page[limit]=50&page[offset]=100",
		"last":  "http://somesite.com/movies?page[limit]=50&page[offset]=500",
	}
	o, err := parseOffset(links)
	if err != nil {
		t.Errorf("parseOffset returned err: %v", err)
	}
	got, want := o, Offset{First: 50, Prev: 0, Next: 100, Last: 500}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseOffset = %#v, want %#v", got, want)
	}
}

func Test_parseOffset_structLink(t *testing.T) {
	structLinks := []jsonapi.Links{
		{"first": struct{}{}},
	}
	for _, link := range structLinks {
		_, err := parseOffset(link)
		if err == nil {
			t.Errorf("parseOffset with struct links expected to return err")
		}
	}
}

func Test_parseOffset_pageNumberAndSize(t *testing.T) {
	links := jsonapi.Links{
		"first": "http://example.com?page[number]=1&page[size]=50",
		"prev":  "http://example.com?page[number]=13&page[size]=50",
		"next":  "http://example.com?page[number]=15&page[size]=50",
		"last":  "http://example.com?page[number]=34&page[size]=50",
	}
	o, err := parseOffset(links)
	if err != nil {
		t.Errorf("parseOffset returned err: %v", err)
	}
	// The Kitsu API uses offset & limit instead of number & size so we expect
	// nothing.
	got, want := o, Offset{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseOffset = %#v, want %#v", got, want)
	}
}

func Test_parseOffset_badLinks(t *testing.T) {
	badLinks := []jsonapi.Links{
		{"first": ":"},
		{"prev": ":"},
		{"next": ":"},
		{"last": ":"},
	}
	for _, link := range badLinks {
		_, err := parseOffset(link)
		if err == nil {
			t.Errorf("parseOffset expected to return err")
		}
	}
}
