package kitsu

import (
	"encoding/json"
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
		testHeader(t, r, "Accept", defaultMediaType)
		fmt.Fprintf(w, `{"data":{"id":"7442","type":"anime","attributes":{"slug":"attack-on-titan"}}}`)
	})

	got, _, err := client.Anime.Show("7442")
	if err != nil {
		t.Errorf("Anime.Show returned error: %v", err)
	}

	want := &Anime{ID: "7442", Slug: "attack-on-titan"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Show anime mismatch\nhave: %#+v\nwant: %#+v", got, want)
	}
}

func TestAnimeService_Show_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
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
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"page[limit]":    "2",
			"page[offset]":   "0",
			"filter[genres]": "sports,sci-fi",
			"sort":           "-followersCount,-followingCount",
			"include":        "media.genres,media.installments",
		})

		const s = `
		{
			"data": [{
				"id": "7442",
				"type": "anime",
				"attributes": {
					"slug": "attack-on-titan"
				}
			}, {
				"id": "7442",
				"type": "anime",
				"attributes": {
					"slug": "attack-on-titan"
				}
			}],
			"links": {
				"first": "https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=0",
				"next": "https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=50",
				"last": "https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=498"
			}
		}`
		fmt.Fprint(w, s)
	})

	opt := &Options{
		PageLimit:  2,
		PageOffset: 0,
		Filter:     "genres",
		FilterVal:  []string{"sports", "sci-fi"},
		Sort:       []string{"-followersCount", "-followingCount"},
		Include:    []string{"media.genres", "media.installments"},
	}

	got, resp, err := client.Anime.List(opt)
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}

	want := []*Anime{
		{ID: "7442", Slug: "attack-on-titan"},
		{ID: "7442", Slug: "attack-on-titan"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List mismatch\nhave: %#+v\nwant: %#+v", got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
	if got, want := resp.FirstOffset, 0; got != want {
		t.Errorf("Anime.List response FirstOffset = %d, want %d", got, want)
	}
	if got, want := resp.LastOffset, 498; got != want {
		t.Errorf("Anime.List response LastOffset = %d, want %d", got, want)
	}
	if got, want := resp.NextOffset, 50; got != want {
		t.Errorf("Anime.List response NextOffset = %d, want %d", got, want)
	}
	if got, want := resp.PrevOffset, 0; got != want {
		t.Errorf("Anime.List response PrevOffset = %d, want %d", got, want)
	}
}

func TestAnimeService_List_addOptionsUnknownFilter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"filter[unknown_filter]": "unknown_value",
		})

		w.WriteHeader(http.StatusBadRequest)
		const s = `{"errors":[{"title":"Filter not allowed","detail":"unknown_filter is not allowed.","code":"102","status":"400"}]}`
		fmt.Fprint(w, s)
	})

	opt := &Options{
		Filter:    "unknown_filter",
		FilterVal: []string{"unknown_value"},
	}
	_, _, err := client.Anime.List(opt)
	if err == nil {
		t.Fatal("Anime.List with unknown filter expected to return err")
	}
	want := []Error{{Code: "102", Detail: "unknown_filter is not allowed.", Status: "400", Title: "Filter not allowed"}}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatal("Anime.List with unknown filter expected to return err of type ErrorResponse")
	}
	if got := errResp.Errors; !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List with unknown filter\nhave: %#v\nwant: %#v", got, want)
	}
}

// TODO: This test will fail once we ask google/jsonapi to unmashal a struct
// field like Charater.Image thus such cases will remain unsupported for now.
func TestAnimeService_List_include(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"include": "castings.character,castings.person",
		})

		const s = `
        {
           "data":[
              {
                 "id":"1",
                 "type":"anime",
                 "attributes":{
                    "slug":"cowboy-bebop"
                 },
                 "relationships":{
                    "castings":{
                       "data":[
                          {
                             "type":"castings",
                             "id":"47"
                          }
                       ]
                    }
                 }
              },
              {
                 "id":"1",
                 "type":"anime",
                 "attributes":{
                    "slug":"cowboy-bebop"
                 },
                 "relationships":{
                    "castings":{
                       "data":[
                          {
                             "type":"castings",
                             "id":"47"
                          }
                       ]
                    }
                 }
              }
           ],
           "included":[
              {
                 "id":"47",
                 "type":"castings",
                 "attributes":{
                    "role":"Voice Actor",
                    "voiceActor":true,
                    "featured":true,
                    "language":"Japanese"
                 },
                 "relationships":{
                    "character":{
                       "data":{
                          "type":"characters",
                          "id":"2"
                       }
                    },
                    "person":{
                       "data":{
                          "type":"people",
                          "id":"47"
                       }
                    }
                 }
              },
              {
                 "id":"47",
                 "type":"people",
                 "attributes":{
                    "name":"Kouichi Yamadera",
                    "malId":11
                 }
              },
              {
                 "id":"2",
                 "type":"characters",
                 "attributes":{
                    "name":"Spike Spiegel",
                    "malId":1,
                    "image":{
                       "original":"https://media.kitsu.io/characters/images/2/original.jpg?1483096805"
                    }
                 }
              }
           ],
           "links":{
              "first":"https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=0",
              "next":"https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=50",
              "last":"https://kitsu.io/api/17/anime?page%5Blimit%5D=50&page%5Boffset%5D=498"
           }
        }`
		fmt.Fprint(w, s)
	})

	opt := &Options{
		Include: []string{"castings.character", "castings.person"},
	}

	got, resp, err := client.Anime.List(opt)
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}

	want := []*Anime{
		{ID: "1", Slug: "cowboy-bebop",
			Castings: []*Casting{
				{
					ID:   "47",
					Role: "Voice Actor", VoiceActor: true, Featured: true,
					Language:  "Japanese",
					Person:    &Person{ID: "47", Name: "Kouichi Yamadera", MALID: 11},
					Character: &Character{ID: "2", Name: "Spike Spiegel", MALID: 1},
				},
			},
		},
		{ID: "1", Slug: "cowboy-bebop",
			Castings: []*Casting{
				{
					ID:   "47",
					Role: "Voice Actor", VoiceActor: true, Featured: true,
					Language:  "Japanese",
					Person:    &Person{ID: "47", Name: "Kouichi Yamadera", MALID: 11},
					Character: &Character{ID: "2", Name: "Spike Spiegel", MALID: 1},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List mismatch\nhave: %#+v\nwant: %#+v", got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
	if got, want := resp.FirstOffset, 0; got != want {
		t.Errorf("Anime.List response FirstOffset = %d, want %d", got, want)
	}
	if got, want := resp.LastOffset, 498; got != want {
		t.Errorf("Anime.List response LastOffset = %d, want %d", got, want)
	}
	if got, want := resp.NextOffset, 50; got != want {
		t.Errorf("Anime.List response NextOffset = %d, want %d", got, want)
	}
	if got, want := resp.PrevOffset, 0; got != want {
		t.Errorf("Anime.List response PrevOffset = %d, want %d", got, want)
	}
}
