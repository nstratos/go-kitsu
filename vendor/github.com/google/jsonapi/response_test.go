package jsonapi

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type Blog struct {
	ID            int       `jsonapi:"primary,blogs"`
	ClientID      string    `jsonapi:"client-id"`
	Title         string    `jsonapi:"attr,title"`
	Posts         []*Post   `jsonapi:"relation,posts"`
	CurrentPost   *Post     `jsonapi:"relation,current_post"`
	CurrentPostID int       `jsonapi:"attr,current_post_id"`
	CreatedAt     time.Time `jsonapi:"attr,created_at"`
	ViewCount     int       `jsonapi:"attr,view_count"`
}

type Post struct {
	Blog
	ID            uint64     `jsonapi:"primary,posts"`
	BlogID        int        `jsonapi:"attr,blog_id"`
	ClientID      string     `jsonapi:"client-id"`
	Title         string     `jsonapi:"attr,title"`
	Body          string     `jsonapi:"attr,body"`
	Comments      []*Comment `jsonapi:"relation,comments"`
	LatestComment *Comment   `jsonapi:"relation,latest_comment"`
}

type Comment struct {
	ID       int    `jsonapi:"primary,comments"`
	ClientID string `jsonapi:"client-id"`
	PostID   int    `jsonapi:"attr,post_id"`
	Body     string `jsonapi:"attr,body"`
}

type Book struct {
	ID          uint64  `jsonapi:"primary,books"`
	Author      string  `jsonapi:"attr,author"`
	ISBN        string  `jsonapi:"attr,isbn"`
	Title       string  `jsonapi:"attr,title,omitempty"`
	Description *string `jsonapi:"attr,description"`
	Pages       *uint   `jsonapi:"attr,pages,omitempty"`
	PublishedAt time.Time
}

type Timestamp struct {
	ID   int        `jsonapi:"primary,timestamps"`
	Time time.Time  `jsonapi:"attr,timestamp,iso8601"`
	Next *time.Time `jsonapi:"attr,next,iso8601"`
}

type Car struct {
	ID    *string `jsonapi:"primary,cars"`
	Make  *string `jsonapi:"attr,make,omitempty"`
	Model *string `jsonapi:"attr,model,omitempty"`
	Year  *uint   `jsonapi:"attr,year,omitempty"`
}

func TestMarshalIDPtr(t *testing.T) {
	id, make, model := "123e4567-e89b-12d3-a456-426655440000", "Ford", "Mustang"
	car := &Car{
		ID:    &id,
		Make:  &make,
		Model: &model,
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, car); err != nil {
		t.Fatal(err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &jsonData); err != nil {
		t.Fatal(err)
	}
	data := jsonData["data"].(map[string]interface{})
	// attributes := data["attributes"].(map[string]interface{})

	// Verify that the ID was sent
	val, exists := data["id"]
	if !exists {
		t.Fatal("Was expecting the data.id member to exist")
	}
	if val != id {
		t.Fatalf("Was expecting the data.id member to be `%s`, got `%s`", id, val)
	}
}

func TestMarshall_invalidIDType(t *testing.T) {
	type badIDStruct struct {
		ID *bool `jsonapi:"primary,cars"`
	}
	id := true
	o := &badIDStruct{ID: &id}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, o); err != ErrBadJSONAPIID {
		t.Fatalf(
			"Was expecting a `%s` error, got `%s`", ErrBadJSONAPIID, err,
		)
	}
}

func TestOmitsEmptyAnnotation(t *testing.T) {
	book := &Book{
		Author:      "aren55555",
		PublishedAt: time.Now().AddDate(0, -1, 0),
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, book); err != nil {
		t.Fatal(err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &jsonData); err != nil {
		t.Fatal(err)
	}
	attributes := jsonData["data"].(map[string]interface{})["attributes"].(map[string]interface{})

	// Verify that the specifically omitted field were omitted
	if val, exists := attributes["title"]; exists {
		t.Fatalf("Was expecting the data.attributes.title key/value to have been omitted - it was not and had a value of %v", val)
	}
	if val, exists := attributes["pages"]; exists {
		t.Fatalf("Was expecting the data.attributes.pages key/value to have been omitted - it was not and had a value of %v", val)
	}

	// Verify the implicity omitted fields were omitted
	if val, exists := attributes["PublishedAt"]; exists {
		t.Fatalf("Was expecting the data.attributes.PublishedAt key/value to have been implicity omitted - it was not and had a value of %v", val)
	}

	// Verify the unset fields were not omitted
	if _, exists := attributes["isbn"]; !exists {
		t.Fatal("Was expecting the data.attributes.isbn key/value to have NOT been omitted")
	}
}

func TestHasPrimaryAnnotation(t *testing.T) {
	testModel := &Blog{
		ID:        5,
		Title:     "Title 1",
		CreatedAt: time.Now(),
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)

	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	data := resp.Data

	if data.Type != "blogs" {
		t.Fatalf("type should have been blogs, got %s", data.Type)
	}

	if data.ID != "5" {
		t.Fatalf("ID not transfered")
	}
}

func TestSupportsAttributes(t *testing.T) {
	testModel := &Blog{
		ID:        5,
		Title:     "Title 1",
		CreatedAt: time.Now(),
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	data := resp.Data

	if data.Attributes == nil {
		t.Fatalf("Expected attributes")
	}

	if data.Attributes["title"] != "Title 1" {
		t.Fatalf("Attributes hash not populated using tags correctly")
	}
}

func TestOmitsZeroTimes(t *testing.T) {
	testModel := &Blog{
		ID:        5,
		Title:     "Title 1",
		CreatedAt: time.Time{},
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	data := resp.Data

	if data.Attributes == nil {
		t.Fatalf("Expected attributes")
	}

	if data.Attributes["created_at"] != nil {
		t.Fatalf("Created at was serialized even though it was a zero Time")
	}
}

func TestMarshalISO8601Time(t *testing.T) {
	testModel := &Timestamp{
		ID:   5,
		Time: time.Date(2016, 8, 17, 8, 27, 12, 23849, time.UTC),
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	data := resp.Data

	if data.Attributes == nil {
		t.Fatalf("Expected attributes")
	}

	if data.Attributes["timestamp"] != "2016-08-17T08:27:12Z" {
		t.Fatal("Timestamp was not serialised into ISO8601 correctly")
	}
}

func TestMarshalISO8601TimePointer(t *testing.T) {
	tm := time.Date(2016, 8, 17, 8, 27, 12, 23849, time.UTC)
	testModel := &Timestamp{
		ID:   5,
		Next: &tm,
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	data := resp.Data

	if data.Attributes == nil {
		t.Fatalf("Expected attributes")
	}

	if data.Attributes["next"] != "2016-08-17T08:27:12Z" {
		t.Fatal("Next was not serialised into ISO8601 correctly")
	}
}

func TestRelations(t *testing.T) {
	testModel := testBlog()

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	relations := resp.Data.Relationships

	if relations == nil {
		t.Fatalf("Relationships were not materialized")
	}

	if relations["posts"] == nil {
		t.Fatalf("Posts relationship was not materialized")
	}

	if relations["current_post"] == nil {
		t.Fatalf("Current post relationship was not materialized")
	}

	if len(relations["posts"].(map[string]interface{})["data"].([]interface{})) != 2 {
		t.Fatalf("Did not materialize two posts")
	}
}

func TestNoRelations(t *testing.T) {
	testModel := &Blog{ID: 1, Title: "Title 1", CreatedAt: time.Now()}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayload(out, testModel); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	if resp.Included != nil {
		t.Fatalf("Encoding json response did not omit included")
	}
}

func TestMarshalOnePayloadWithoutIncluded(t *testing.T) {
	data := &Post{
		ID:       1,
		BlogID:   2,
		ClientID: "123e4567-e89b-12d3-a456-426655440000",
		Title:    "Foo",
		Body:     "Bar",
		Comments: []*Comment{
			&Comment{
				ID:   20,
				Body: "First",
			},
			&Comment{
				ID:   21,
				Body: "Hello World",
			},
		},
		LatestComment: &Comment{
			ID:   22,
			Body: "Cool!",
		},
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayloadWithoutIncluded(out, data); err != nil {
		t.Fatal(err)
	}

	resp := new(OnePayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	if resp.Included != nil {
		t.Fatalf("Encoding json response did not omit included")
	}
}

func TestMarshalMany(t *testing.T) {
	data := []interface{}{
		&Blog{
			ID:        5,
			Title:     "Title 1",
			CreatedAt: time.Now(),
			Posts: []*Post{
				&Post{
					ID:    1,
					Title: "Foo",
					Body:  "Bar",
				},
				&Post{
					ID:    2,
					Title: "Fuubar",
					Body:  "Bas",
				},
			},
			CurrentPost: &Post{
				ID:    1,
				Title: "Foo",
				Body:  "Bar",
			},
		},
		&Blog{
			ID:        6,
			Title:     "Title 2",
			CreatedAt: time.Now(),
			Posts: []*Post{
				&Post{
					ID:    3,
					Title: "Foo",
					Body:  "Bar",
				},
				&Post{
					ID:    4,
					Title: "Fuubar",
					Body:  "Bas",
				},
			},
			CurrentPost: &Post{
				ID:    4,
				Title: "Foo",
				Body:  "Bar",
			},
		},
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalManyPayload(out, data); err != nil {
		t.Fatal(err)
	}

	resp := new(ManyPayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	d := resp.Data

	if len(d) != 2 {
		t.Fatalf("data should have two elements")
	}
}

func TestMarshalMany_WithSliceOfStructPointers(t *testing.T) {
	var data []*Blog
	for len(data) < 2 {
		data = append(data, testBlog())
	}

	out := bytes.NewBuffer(nil)
	if err := MarshalManyPayload(out, data); err != nil {
		t.Fatal(err)
	}

	resp := new(ManyPayload)
	if err := json.NewDecoder(out).Decode(resp); err != nil {
		t.Fatal(err)
	}

	d := resp.Data

	if len(d) != 2 {
		t.Fatalf("data should have two elements")
	}
}

func TestMarshalMany_SliceOfInterfaceAndSliceOfStructsSameJSON(t *testing.T) {
	structs := []*Book{
		&Book{ID: 1, Author: "aren55555", ISBN: "abc"},
		&Book{ID: 2, Author: "shwoodard", ISBN: "xyz"},
	}
	interfaces := []interface{}{}
	for _, s := range structs {
		interfaces = append(interfaces, s)
	}

	// Perform Marshals
	structsOut := new(bytes.Buffer)
	if err := MarshalManyPayload(structsOut, structs); err != nil {
		t.Fatal(err)
	}
	interfacesOut := new(bytes.Buffer)
	if err := MarshalManyPayload(interfacesOut, interfaces); err != nil {
		t.Fatal(err)
	}

	// Generic JSON Unmarshal
	structsData, interfacesData :=
		make(map[string]interface{}), make(map[string]interface{})
	if err := json.Unmarshal(structsOut.Bytes(), &structsData); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(interfacesOut.Bytes(), &interfacesData); err != nil {
		t.Fatal(err)
	}

	// Compare Result
	if !reflect.DeepEqual(structsData, interfacesData) {
		t.Fatal("Was expecting the JSON API generated to be the same")
	}
}

func TestMarshalMany_InvalidIntefaceArgument(t *testing.T) {
	out := new(bytes.Buffer)
	if err := MarshalManyPayload(out, true); err != ErrExpectedSlice {
		t.Fatal("Was expecting an error")
	}
	if err := MarshalManyPayload(out, 25); err != ErrExpectedSlice {
		t.Fatal("Was expecting an error")
	}
	if err := MarshalManyPayload(out, Book{}); err != ErrExpectedSlice {
		t.Fatal("Was expecting an error")
	}
}

func testBlog() *Blog {
	return &Blog{
		ID:        5,
		Title:     "Title 1",
		CreatedAt: time.Now(),
		Posts: []*Post{
			&Post{
				ID:    1,
				Title: "Foo",
				Body:  "Bar",
				Comments: []*Comment{
					&Comment{
						ID:   1,
						Body: "foo",
					},
					&Comment{
						ID:   2,
						Body: "bar",
					},
				},
				LatestComment: &Comment{
					ID:   1,
					Body: "foo",
				},
			},
			&Post{
				ID:    2,
				Title: "Fuubar",
				Body:  "Bas",
				Comments: []*Comment{
					&Comment{
						ID:   1,
						Body: "foo",
					},
					&Comment{
						ID:   3,
						Body: "bas",
					},
				},
				LatestComment: &Comment{
					ID:   1,
					Body: "foo",
				},
			},
		},
		CurrentPost: &Post{
			ID:    1,
			Title: "Foo",
			Body:  "Bar",
			Comments: []*Comment{
				&Comment{
					ID:   1,
					Body: "foo",
				},
				&Comment{
					ID:   2,
					Body: "bar",
				},
			},
			LatestComment: &Comment{
				ID:   1,
				Body: "foo",
			},
		},
	}
}
