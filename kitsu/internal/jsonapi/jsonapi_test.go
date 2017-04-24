package jsonapi

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

type Anime struct {
	ID   string `jsonapi:"primary,anime"`
	Slug string `jsonapi:"attr,slug"`
}

func TestEncode_one(t *testing.T) {
	in := &Anime{Slug: "bebob"}
	out := `{"data":{"type":"anime","attributes":{"slug":"bebob"}}}` + "\n"

	buf := &bytes.Buffer{}
	if err := Encode(buf, in); err != nil {
		t.Fatalf("Encode returned err: %v", err)
	}
	if got, want := buf.String(), out; got != want {
		t.Errorf("Encode \nhave: %q\nwant: %q", got, want)
	}
}

func TestEncode_many(t *testing.T) {
	in := []*Anime{{Slug: "foo"}, {Slug: "bar"}}
	out := `{"data":[{"type":"anime","attributes":{"slug":"foo"}},{"type":"anime","attributes":{"slug":"bar"}}]}` + "\n"

	buf := &bytes.Buffer{}
	if err := Encode(buf, in); err != nil {
		t.Fatalf("Encode returned err: %v", err)
	}
	if got, want := buf.String(), out; got != want {
		t.Errorf("Encode \nhave: %q\nwant: %q", got, want)
	}
}

// Encoding a type that is not a pointer to struct or a slice of pointers to
// structs should return an error.
func TestEncode_invalidTypes(t *testing.T) {
	buf := &bytes.Buffer{}
	a := 1

	var tests = []struct {
		in interface{}
	}{
		{a},
		{&a},
		{[]int{a, a}},
		{[]*int{&a, &a}},
	}

	for _, tt := range tests {
		if err := Encode(buf, tt.in); err == nil {
			t.Errorf("Encode(%T) expected to return err", tt.in)
		}
	}
}

func TestEncode_nil(t *testing.T) {
	buf := &bytes.Buffer{}
	var a *Anime
	err := Encode(buf, a)
	if err == nil {
		t.Errorf("Encode(%#v) expected to return err", a)
	}

	var anime []*Anime
	err = Encode(buf, anime)
	if err == nil {
		t.Errorf("Encode(%#v) expected to return err", a)
	}
}

func TestDecode_one(t *testing.T) {
	in := `{"data":{"type":"anime","id":"","attributes":{"slug":"bebob"}}}`
	r := strings.NewReader(in)

	a := new(Anime)
	o, err := Decode(r, a)
	if err != nil {
		t.Errorf("Decode returned err: %v", err)
	}
	if got, want := a, (&Anime{Slug: "bebob"}); !reflect.DeepEqual(got, want) {
		t.Errorf("Decode \nhave: %+v\nwant: %+v", got, want)
	}
	if got, want := o, (Offset{}); !reflect.DeepEqual(got, want) {
		t.Errorf("Decode \nhave: %+v\nwant: %+v", got, want)
	}
}

func TestDecode_many(t *testing.T) {
	in := `{"data":[{"type":"anime","id":"","attributes":{"slug":"foo"}},{"type":"anime","id":"","attributes":{"slug":"bar"}}]}`

	r := strings.NewReader(in)
	var anime []*Anime
	_, err := Decode(r, &anime)
	if err != nil {
		t.Errorf("Decode returned err: %v", err)
	}

	want := []*Anime{{Slug: "foo"}, {Slug: "bar"}}
	if got := anime; !reflect.DeepEqual(got, want) {
		t.Errorf("Decode \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestDecode_manyWithLinks(t *testing.T) {
	in := `{
  "data":[{"type":"anime","id":"1"},{"type":"anime","id":"2"}],
  "links": {
    "first": "http://somesite.com/movies?page[limit]=50&page[offset]=50",
    "prev": "http://somesite.com/movies?page[limit]=50&page[offset]=0",
    "next": "http://somesite.com/movies?page[limit]=50&page[offset]=100",
    "last": "http://somesite.com/movies?page[limit]=50&page[offset]=500"
  }
}`

	r := strings.NewReader(in)
	var anime []*Anime
	o, err := Decode(r, &anime)
	if err != nil {
		t.Errorf("Decode returned err: %v", err)
	}

	want := []*Anime{{ID: "1"}, {ID: "2"}}
	if got := anime; !reflect.DeepEqual(got, want) {
		t.Errorf("Decode \nhave: %#v\nwant: %#v", got, want)
	}

	if got, want := o, (Offset{First: 50, Last: 500, Next: 100}); !reflect.DeepEqual(got, want) {
		t.Errorf("Decode offset \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestDecode_manyWithBadLinks(t *testing.T) {
	in := `{ "data":[],"links":{ "first": ":"}}`

	r := strings.NewReader(in)
	var anime []*Anime
	_, err := Decode(r, &anime)
	if err == nil {
		t.Errorf("Decode with bad links expected to return err")
	}
}

// Decoding to a type that is not a pointer to struct or a pointer to slice of
// pointers to structs should return an error.
func TestDecode_toBadType(t *testing.T) {
	in := `{"data":[{"type":"anime","id":"1"},{"type":"anime","id":"2"}]}`

	r := strings.NewReader(in)
	var anime []*Anime

	if _, err := Decode(r, anime); err == nil {
		t.Errorf("Decode(%T) expected to return err", anime)
	}

	var i *int
	if _, err := Decode(r, i); err == nil {
		t.Errorf("Decode(%T) expected to return err", i)
	}
}

func TestDecode_toSliceForSinglePayload(t *testing.T) {
	in := `{"data":{"type":"anime","id":"","attributes":{"slug":"bebob"}}}`

	r := strings.NewReader(in)
	var anime []*Anime

	_, err := Decode(r, &anime)
	if err == nil {
		t.Errorf("Decode(%v, %T) expected to return err", in, &anime)
	}
}

func TestEncodeOne(t *testing.T) {
	in := &Anime{Slug: "bebob"}
	out := `{"data":{"type":"anime","attributes":{"slug":"bebob"}}}` + "\n"
	buf := &bytes.Buffer{}
	if err := EncodeOne(buf, in); err != nil {
		t.Errorf("EncodeOne returned err: %v", err)
	}
	if got, want := buf.String(), out; got != want {
		t.Errorf("EncodeOne \nhave: %q\nwant: %q", got, want)
	}
}

func TestEncodeMany(t *testing.T) {
	in := []*Anime{{Slug: "foo"}, {Slug: "bar"}}
	out := `{"data":[{"type":"anime","attributes":{"slug":"foo"}},{"type":"anime","attributes":{"slug":"bar"}}]}` + "\n"
	buf := &bytes.Buffer{}
	if err := EncodeMany(buf, in); err != nil {
		t.Errorf("EncodeMany returned err: %v", err)
	}
	if got, want := buf.String(), out; got != want {
		t.Errorf("EncodeMany \nhave: %q\nwant: %q", got, want)
	}
}

func TestDecodeOne(t *testing.T) {
	in := `{"data":{"type":"anime","id":"","attributes":{"slug":"bebob"}}}`
	r := strings.NewReader(in)

	a := new(Anime)
	if err := DecodeOne(r, a); err != nil {
		t.Errorf("DecodeOne returned err: %v", err)
	}
	want := &Anime{Slug: "bebob"}
	if got := a; !reflect.DeepEqual(got, want) {
		t.Errorf("DecodeOne \nhave: %+v\nwant: %+v", got, want)
	}
}

func TestDecodeMany(t *testing.T) {
	in := `{"data":[{"type":"anime","id":"","attributes":{"slug":"foo"}},{"type":"anime","id":"","attributes":{"slug":"bar"}}]}`

	r := strings.NewReader(in)
	anime, _, err := DecodeMany(r, reflect.TypeOf(&Anime{}))
	if err != nil {
		t.Errorf("DecodeMany returned err: %v", err)
	}

	want := []interface{}{&Anime{Slug: "foo"}, &Anime{Slug: "bar"}}
	if got := anime; !reflect.DeepEqual(got, want) {
		t.Errorf("DecodeMany \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestDecodeMany_badBody(t *testing.T) {
	in := `{"data":{"type":"anime","id":"","attributes":{"slug":"bebob"}}}`

	r := strings.NewReader(in)
	_, _, err := DecodeMany(r, reflect.TypeOf(&Anime{}))
	if err == nil {
		t.Errorf("DecodeMany with bad body expected to return err")
	}
}

func TestDecodeMany_withLinks(t *testing.T) {
	in := `{
  "data":[{"type":"anime","id":"1"},{"type":"anime","id":"2"}],
  "links": {
    "first": "http://somesite.com/movies?page[limit]=50&page[offset]=50",
    "prev": "http://somesite.com/movies?page[limit]=50&page[offset]=0",
    "next": "http://somesite.com/movies?page[limit]=50&page[offset]=100",
    "last": "http://somesite.com/movies?page[limit]=50&page[offset]=500"
  }
}`

	r := strings.NewReader(in)
	anime, o, err := DecodeMany(r, reflect.TypeOf(&Anime{}))
	if err != nil {
		t.Errorf("DecodeMany returned err: %v", err)
	}

	want := []interface{}{&Anime{ID: "1"}, &Anime{ID: "2"}}
	if got := anime; !reflect.DeepEqual(got, want) {
		t.Errorf("DecodeMany \nhave: %#v\nwant: %#v", got, want)
	}

	if got, want := o, (Offset{First: 50, Last: 500, Next: 100}); !reflect.DeepEqual(got, want) {
		t.Errorf("DecodeMany offset \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestDecodeMany_withBadLinks(t *testing.T) {
	in := `{ "data":[],"links":{ "first": ":"}}`

	r := strings.NewReader(in)
	_, _, err := DecodeMany(r, reflect.TypeOf(&Anime{}))
	if err == nil {
		t.Errorf("DecodeMany with bad links expected to return err")
	}
}
