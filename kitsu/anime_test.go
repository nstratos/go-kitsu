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
		testFormValues(t, r, values{
			"filter[genres]": "sports,sci-fi",
			"sort":           "-followersCount,-followingCount",
			"include":        "media.genres,media.installments",
		})
		fmt.Fprintf(w, `{"data":{"id":"7442","type":"anime","attributes":{"slug":"attack-on-titan"}}}`)
	})

	got, _, err := client.Anime.Show("7442",
		Filter("genres", "sports", "sci-fi"),
		Sort("-followersCount", "-followingCount"),
		Include("media.genres", "media.installments"),
	)
	if err != nil {
		t.Errorf("Anime.Show returned error: %v", err)
	}

	want := &Anime{ID: "7442", Slug: "attack-on-titan"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Show anime mismatch\nhave: %#+v\nwant: %#+v", got, want)
	}
}

func TestAnimeService_Show_decodeAttributes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime/7442", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		fmt.Fprintf(w, `{
			"data":{
				"id":"7442",
				"type":"anime",
				"attributes":{
					"slug": "attack-on-titan",
					"synopsis": "Several hundred years ago, humans were nearly exterminated by titans...",
					"coverImageTopOffset": 263,
					"titles": {
						"en": "Attack on Titan",
						"en_jp": "Shingeki no Kyojin",
						"ja_jp": "進撃の巨人"
					},
					"canonical_title": "Attack on Titan",
					"abbreviatedTitles": [ "AoT", "AT" ],
					"averageRating": 4.26984658306698,
					"ratingFrequencies": {
						"0.5": "114",
						"1.0": "279",
						"1.5": "146",
						"2.0": "359",
						"2.5": "763",
						"3.0": "2331",
						"3.5": "3034",
						"4.0": "5619",
						"4.5": "5951",
						"5.0": "12878"
					},
					"startDate": "2013-04-07",
					"endDate": "2013-09-28",
					"posterImage": {
						"original": "https://static.hummingbird.me/anime/7442/poster/$1.png"
					},
					"coverImage": {
						"original": "https://static.hummingbird.me/anime/7442/cover/$1.png"
					},
					"episodeCount": 25,
					"episodeLength": 24,
					"showType": "TV",
					"youtubeVideoId": "n4Nj6Y_SNYI",
					"ageRating": "R",
					"ageRatingGuide": "Violence, Profanity"
				}
			}
		}`)
	})

	got, _, err := client.Anime.Show("7442")
	if err != nil {
		t.Fatalf("Anime.Show returned error: %v", err)
	}

	want := &Anime{
		ID:                  "7442",
		Slug:                "attack-on-titan",
		Synopsis:            "Several hundred years ago, humans were nearly exterminated by titans...",
		CoverImageTopOffset: 263,
		Titles: map[string]interface{}{
			"en":    "Attack on Titan",
			"en_jp": "Shingeki no Kyojin",
			"ja_jp": "進撃の巨人",
		},
		CanonicalTitle:    "Attack on Titan",
		AbbreviatedTitles: []string{"AoT", "AT"},
		AverageRating:     4.26984658306698,
		RatingFrequencies: map[string]interface{}{
			"0.5": "114",
			"1.0": "279",
			"1.5": "146",
			"2.0": "359",
			"2.5": "763",
			"3.0": "2331",
			"3.5": "3034",
			"4.0": "5619",
			"4.5": "5951",
			"5.0": "12878",
		},
		StartDate: "2013-04-07",
		EndDate:   "2013-09-28",
		PosterImage: map[string]interface{}{
			"original": "https://static.hummingbird.me/anime/7442/poster/$1.png",
		},
		CoverImage: map[string]interface{}{
			"original": "https://static.hummingbird.me/anime/7442/cover/$1.png",
		},
		EpisodeCount:   25,
		EpisodeLength:  24,
		ShowType:       "TV",
		YoutubeVideoID: "n4Nj6Y_SNYI",
		AgeRating:      "R",
		AgeRatingGuide: "Violence, Profanity",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Show decode attributes mismatch\nhave: %#+v\nwant: %#+v", got, want)
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
	_, _, err := client.Anime.Show("%", nil)
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
					"slug": "attack-on-titan",
					"showType": "TV"
				}
			}, {
				"id": "7442",
				"type": "anime",
				"attributes": {
					"slug": "attack-on-titan",
					"showType": "TV"
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

	got, resp, err := client.Anime.List(
		Pagination(2, 0),
		Filter("genres", "sports", "sci-fi"),
		Sort("-followersCount", "-followingCount"),
		Include("media.genres", "media.installments"),
	)
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}

	want := []*Anime{
		{ID: "7442", Slug: "attack-on-titan", ShowType: AnimeTypeTV},
		{ID: "7442", Slug: "attack-on-titan", ShowType: AnimeTypeTV},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List mismatch\nhave: %#+v\nwant: %#+v", got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
	offset := PageOffset{First: 0, Last: 498, Next: 50, Prev: 0}
	if got, want := resp.Offset, offset; got != want {
		t.Errorf("Anime.List response Offset = %+v, want %+v", got, want)
	}
}

func TestAnimeService_List_filterOptionWithUnknownAttribute(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"filter[unknown_attribute]": "unknown_value",
		})

		w.WriteHeader(http.StatusBadRequest)
		const s = `{"errors":[{"title":"Filter not allowed","detail":"unknown_attribute is not allowed.","code":"102","status":"400"}]}`
		fmt.Fprint(w, s)
	})

	_, _, err := client.Anime.List(Filter("unknown_attribute", "unknown_value"))
	if err == nil {
		t.Fatal("Anime.List with unknown filter expected to return err")
	}
	want := []Error{{Code: "102", Detail: "unknown_attribute is not allowed.", Status: "400", Title: "Filter not allowed"}}
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

	got, resp, err := client.Anime.List(Include("castings.character", "castings.person"))
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}

	want := []*Anime{
		{ID: "1", Slug: "cowboy-bebop",
			Castings: []*Casting{
				{
					ID:   "47",
					Role: "Voice Actor", VoiceActor: true, Featured: true,
					Language: "Japanese",
					Person:   &Person{ID: "47", Name: "Kouichi Yamadera", MALID: 11},
					Character: &Character{ID: "2", Name: "Spike Spiegel", MALID: 1,
						Image: map[string]interface{}{"original": "https://media.kitsu.io/characters/images/2/original.jpg?1483096805"},
					},
				},
			},
		},
		{ID: "1", Slug: "cowboy-bebop",
			Castings: []*Casting{
				{
					ID:   "47",
					Role: "Voice Actor", VoiceActor: true, Featured: true,
					Language: "Japanese",
					Person:   &Person{ID: "47", Name: "Kouichi Yamadera", MALID: 11},
					Character: &Character{ID: "2", Name: "Spike Spiegel", MALID: 1,
						Image: map[string]interface{}{"original": "https://media.kitsu.io/characters/images/2/original.jpg?1483096805"},
					},
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

	offset := PageOffset{First: 0, Last: 498, Next: 50, Prev: 0}
	if got, want := resp.Offset, offset; got != want {
		t.Errorf("Anime.List response Offset = %+v, want %+v", got, want)
	}
}
