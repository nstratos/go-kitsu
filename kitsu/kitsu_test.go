package kitsu

import (
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
	req, _ := c.NewRequest("GET", inURL, nil)

	// Test that the client's base URL is added to the endpoint.
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %q, want %q", inURL, got, want)
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
