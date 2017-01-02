package kitsu

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestAnimeService_Show(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime/7442", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"data":{"id":"7442","type":"anime","attributes":{"slug":"attack-on-titan"}}}`)
	})

	got, _, err := client.Anime.Show("7442")
	if err != nil {
		t.Errorf("Anime.Show returned error: %v", err)
	}

	want := &AnimeShowResponse{Data: &Anime{Resource: Resource{ID: "7442", Type: "anime"}, Attributes: AnimeAttributes{Slug: "attack-on-titan"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Show anime is \n%+v, want \n%+v", got, want)
	}
}

func TestAnimeService_Show_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, `{"errors":[{"title":"Record not found","detail":"The record identified by 0 could not be found.","code":"404","status":"404"}]}`, http.StatusNotFound)
	})

	_, resp, err := client.Anime.Show("0")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestAnimeService_Show_invalidID(t *testing.T) {
	_, _, err := client.Anime.Show("%")
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestAnimeService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"data":[{"id":"7442","type":"anime","attributes":{"slug":"attack-on-titan"}},{"id":"7442","type":"anime","attributes":{"slug":"attack-on-titan"}}]}`)
	})

	got, _, err := client.Anime.List()
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}

	want := &AnimeListResponse{
		Data: []*Anime{
			{Resource: Resource{ID: "7442", Type: "anime"}, Attributes: AnimeAttributes{Slug: "attack-on-titan"}},
			{Resource: Resource{ID: "7442", Type: "anime"}, Attributes: AnimeAttributes{Slug: "attack-on-titan"}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List returns \n%+v, want \n%+v", got, want)
	}
}
