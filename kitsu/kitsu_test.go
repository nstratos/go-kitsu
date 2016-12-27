package kitsu

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &struct{ Foo string }{Foo: "bar"}, `{"Foo":"bar"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// Test that the client's base URL is added to the endpoint.
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %q, want %q", inURL, got, want)
	}

	// Test that body gets encoded to JSON.
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%#v) Body is \n%q, want \n%q", inBody, got, want)
	}

	// Test that the correct Content-Type gets added.
	if got, want := req.Header.Get("Content-Type"), defaultMediaType; got != want {
		t.Errorf("NewRequest() Content-Type is %q, want %q", got, want)
	}
}

func TestClient_NewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	inURL := ":"
	_, err := c.NewRequest("GET", inURL, nil)
	if err == nil {
		t.Errorf("NewRequest(%q) should return parse err", inURL)
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestClient_NewRequest_badBody(t *testing.T) {
	c := NewClient(nil)

	type Foo struct {
		Bar map[interface{}]interface{}
	}
	inBody := &Foo{}
	_, err := c.NewRequest("GET", "/", inBody)

	if err == nil {
		t.Errorf("NewRequest(%#v) should return err", inBody)
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected JSON Unsupported type error, got %#v.", err)
	}
}

func TestClient_NewRequest_emptyBody(t *testing.T) {
	c := NewClient(nil)
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest with empty body returned error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("NewRequest with empty body should construct request with nil Body")
	}
}
