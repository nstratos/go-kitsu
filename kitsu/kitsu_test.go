package kitsu

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	// client is the kitsu client used for these tests.
	client *Client

	// server is a test HTTP server that is being started on each test with
	// setup() to provide mock API responses.
	server *httptest.Server

	// mux is the HTTP request multiplexer that the test HTTP server uses.
	mux *http.ServeMux
)

func setup() {
	// Starting new test server with mux as its multiplexer.
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// Configuring kitsu.Client to use the test HTTP server URL.
	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	if err := r.ParseForm(); err != nil {
		t.Error("ParseForm returned error:", err)
	}
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: \n%v, want \n%v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	req, err := c.NewRequest("GET", inURL, nil)
	if err != nil {
		t.Fatalf("NewRequest(%q) returned err: %v", inURL, err)
	}

	// Test that the client's base URL is added to the endpoint.
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %q, want %q", inURL, got, want)
	}

}

func TestClient_NewRequest_encode(t *testing.T) {
	var tests = []struct {
		in  interface{}
		out string
	}{
		{
			&Anime{ID: "1"},
			`{"data":{"type":"anime","id":"1"}}` + "\n",
		},
	}
	c := NewClient(nil)
	for _, tt := range tests {
		req, err := c.NewRequest("GET", "/foo", tt.in)
		if err != nil {
			t.Fatalf("NewRequest(%#v) returned err: %v", tt.in, err)
		}

		// Test that body gets encoded to JSON API.
		body, _ := ioutil.ReadAll(req.Body)
		if got, want := string(body), tt.out; got != want {
			t.Errorf("NewRequest(%#v) Body \nhave: %q\nwant: %q", tt.in, got, want)
		}

		// Test that the correct Content-Type gets added.
		if got, want := req.Header.Get("Content-Type"), defaultMediaType; got != want {
			t.Errorf("NewRequest() Content-Type is %q, want %q", got, want)
		}
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

	inBody := int(1)
	_, err := c.NewRequest("GET", "/", inBody)

	if err == nil {
		t.Errorf("NewRequest(%#v) should return err", inBody)
	}
}

func TestClient_NewRequest_emptyBody(t *testing.T) {
	c := NewClient(nil)
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest with empty body returned error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("NewRequest with empty body should construct request with nil Body.")
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar string `jsonapi:"primary,foo"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"data":{"id":"foobar","type":"foo"}}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	got := new(foo)
	_, err := client.Do(req, got)
	if err != nil {
		t.Fatalf("Do(%#v) returned err: %v", got, err)
	}

	want := &foo{Bar: "foobar"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response body = %v, want %v", got, want)
	}
}

func TestClient_Do_badDecodeType(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar string `jsonapi:"primary,not_foo"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"data":{"id":"foobar","type":"foo"}}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	got := new(foo)
	_, err := client.Do(req, got)
	if err == nil {
		t.Fatalf("Do with bad decode type expected to return err")
	}
}

func TestClient_Do_noDecode(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)
	if err != nil {
		t.Fatalf("Do(%v) returned err: %v", nil, err)
	}
}

func TestClient_Do_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"errors":[{"status":"400","title":"Bad Request"}]}`, http.StatusBadRequest)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)
	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func TestClient_Do_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected URL error, got %#v.", err)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	resp := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Response: resp}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func TestError_Error(t *testing.T) {
	err := Error{}
	if err.Error() == "" {
		t.Errorf("Expected non-empty Error.Error()")
	}
}
