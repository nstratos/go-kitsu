package kitsu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestLibraryService_Show(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries/5269457", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		fmt.Fprintf(w, `{"data":{"id":"5269457","type":"libraryEntries","attributes":{"status":"dropped"}}}`)
	})

	got, _, err := client.Library.Show("5269457")
	if err != nil {
		t.Errorf("Library.Show returned error: %v", err)
	}

	want := &LibraryEntry{ID: "5269457", Status: "dropped"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.Show anime mismatch\nhave: %#+v\nwant: %#+v", got, want)
	}
}

func TestLibraryService_Show_decodeAttributes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries/5269457", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"include": "user",
		})
		fmt.Fprintf(w, `{
			"data":{
				"id":"5269457",
				"type":"libraryEntries",
				"attributes":{
					"status":"dropped",
					"progress":3,
					"volumesOwned":0,
					"reconsuming":false,
					"reconsumeCount":0,
					"notes":"",
					"private":false,
					"updatedAt":"2014-05-14T11:54:26.310Z",
					"progressedAt":"2014-05-14T11:54:26.310Z",
					"startedAt":null,
					"finishedAt":null,
					"rating":"0.0",
					"ratingTwenty":null
				},
				"relationships":{
					"user":{
						"links":{
							"self":"https://kitsu.io/api/edge/library-entries/5269457/relationships/user",
							"related":"https://kitsu.io/api/edge/library-entries/5269457/user"
						},
						"data":{
							"type":"users",
							"id":"43133"
						}
					}
				}
			},
			"included":[  
			      {  
					"id":"43133",
					"type":"users",
					"attributes":{  
						"name":"predator914"
					}
				}
			]
		}`)
	})

	got, _, err := client.Library.Show("5269457", Include("user"))
	if err != nil {
		t.Fatalf("Library.Show returned error: %v", err)
	}

	want := &LibraryEntry{
		ID:             "5269457",
		Status:         LibraryEntryStatusDropped,
		Progress:       3,
		Reconsuming:    false,
		ReconsumeCount: 0,
		Notes:          "",
		Private:        false,
		UpdatedAt:      "2014-05-14T11:54:26.310Z",
		Rating:         "0.0",
		User: &User{
			ID:   "43133",
			Name: "predator914",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.Show decode attributes mismatch\nhave: %#+v\nwant: %#+v", got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
}

func TestLibraryService_Show_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		http.Error(w, `{"errors":[{"title":"Record not found","detail":"The record identified by 0 could not be found.","code":"404","status":"404"}]}`, http.StatusNotFound)
	})

	_, resp, err := client.Library.Show("0")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestLibraryService_Show_invalidID(t *testing.T) {
	_, _, err := client.Library.Show("%", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestLibraryService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"filter[userId]": "5554",
		})

		const s = `
		{
		   "data":[
		      {
		         "id":"747296",
		         "type":"libraryEntries",
		         "links":{
		            "self":"https://kitsu.io/api/edge/library-entries/747296"
		         },
		         "attributes":{
		            "status":"completed",
		            "progress":12,
		            "volumesOwned":0,
		            "reconsuming":false,
		            "reconsumeCount":0,
		            "notes":null,
		            "private":false,
		            "updatedAt":"2016-09-06T06:23:05.771Z",
		            "progressedAt":"2016-09-06T06:23:05.771Z",
		            "startedAt":null,
		            "finishedAt":"2016-09-06T06:23:05.771Z",
		            "rating":"3.5",
		            "ratingTwenty":14
		         }
		      },
		      {
		         "id":"747297",
		         "type":"libraryEntries",
		         "links":{
		            "self":"https://kitsu.io/api/edge/library-entries/747297"
		         },
		         "attributes":{
		            "status":"on_hold",
		            "progress":8,
		            "volumesOwned":0,
		            "reconsuming":true,
		            "reconsumeCount":0,
		            "notes":"you should watch it",
		            "private":false,
		            "updatedAt":"2016-04-14T00:56:32.652Z",
		            "progressedAt":"2016-04-14T00:56:32.652Z",
		            "startedAt":null,
		            "finishedAt":null,
		            "rating":"5.0",
		            "ratingTwenty":20
		         }
		      }
		   ],
		   "meta":{
		      "statusCounts":{
		         "current":31,
		         "planned":98,
		         "completed":132,
		         "onHold":10,
		         "dropped":2
		      },
		      "count":273
		   },
		   "links":{
		      "first":"https://kitsu.io/api/edge/library-entries?filter%5BuserId%5D=5554&page%5Blimit%5D=10&page%5Boffset%5D=0",
		      "next":"https://kitsu.io/api/edge/library-entries?filter%5BuserId%5D=5554&page%5Blimit%5D=10&page%5Boffset%5D=10",
		      "last":"https://kitsu.io/api/edge/library-entries?filter%5BuserId%5D=5554&page%5Blimit%5D=10&page%5Boffset%5D=263"
		   }
		}`
		fmt.Fprint(w, s)
	})

	got, resp, err := client.Library.List(
		Filter("userId", "5554"),
	)
	if err != nil {
		t.Errorf("Library.List returned error: %v", err)
	}

	want := []*LibraryEntry{
		{
			ID:             "747296",
			Status:         LibraryEntryStatusCompleted,
			Progress:       12,
			Reconsuming:    false,
			ReconsumeCount: 0,
			Notes:          "",
			Private:        false,
			UpdatedAt:      "2016-09-06T06:23:05.771Z",
			Rating:         "3.5",
		},
		{
			ID:             "747297",
			Status:         LibraryEntryStatusOnHold,
			Progress:       8,
			Reconsuming:    true,
			ReconsumeCount: 0,
			Notes:          "you should watch it",
			Private:        false,
			UpdatedAt:      "2016-04-14T00:56:32.652Z",
			Rating:         "5.0",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.List mismatch\nhave: %#+v\nwant: %#+v", got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
	offset := PageOffset{First: 0, Last: 263, Next: 10, Prev: 0}
	if got, want := resp.Offset, offset; got != want {
		t.Errorf("Library.List response Offset = %+v, want %+v", got, want)
	}
}

func TestLibraryService_List_filterOptionWithUnknownAttribute(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", defaultMediaType)
		testFormValues(t, r, values{
			"filter[unknown_attribute]": "unknown_value",
		})

		w.WriteHeader(http.StatusBadRequest)
		const s = `{"errors":[{"title":"Filter not allowed","detail":"unknown_attribute is not allowed.","code":"102","status":"400"}]}`
		fmt.Fprint(w, s)
	})

	_, _, err := client.Library.List(Filter("unknown_attribute", "unknown_value"))
	if err == nil {
		t.Fatal("Library.List with unknown filter expected to return err")
	}
	want := []Error{{Code: "102", Detail: "unknown_attribute is not allowed.", Status: "400", Title: "Filter not allowed"}}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatal("Library.List with unknown filter expected to return err of type ErrorResponse")
	}
	if got := errResp.Errors; !reflect.DeepEqual(got, want) {
		t.Errorf("Library.List with unknown filter\nhave: %#v\nwant: %#v", got, want)
	}
}

func TestLibraryService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", defaultMediaType)
		testHeader(t, r, "Content-Type", defaultMediaType)

		const serverResponse = `
        {
           "data":{
              "id":"20181115",
              "type":"libraryEntries",
              "links":{
                 "self":"https://kitsu.io/api/edge/library-entries/20181115"
              },
              "attributes":{
                 "createdAt":"2018-02-19T17:44:36.911Z",
                 "updatedAt":"2018-02-19T17:44:36.911Z",
                 "status":"current",
                 "progress":4,
                 "volumesOwned":0,
                 "reconsuming":false,
                 "reconsumeCount":0,
                 "notes":null,
                 "private":false,
                 "reactionSkipped":"unskipped",
                 "progressedAt":"2018-02-19T17:44:36.911Z",
                 "startedAt":"2018-02-19T17:44:36.911Z",
                 "finishedAt":null,
                 "rating":"0.5",
                 "ratingTwenty":2
              }
           }
        }`
		fmt.Fprint(w, serverResponse)
	})

	newEntry := &LibraryEntry{
		Status:   LibraryEntryStatusCurrent,
		Progress: 4,
		Rating:   "0.5",
		User: &User{
			ID: "183388",
		},
		Media: &Anime{
			ID: "1",
		},
	}

	got, _, err := client.Library.Create(newEntry) //kitsu.Include("anime", "user"),
	if err != nil {
		t.Fatal("could not create library:", err)
	}

	want := &LibraryEntry{
		ID:        "20181115",
		Status:    LibraryEntryStatusCurrent,
		Progress:  4,
		Rating:    "0.5",
		UpdatedAt: "2018-02-19T17:44:36.911Z",
	}
	deepEqual(t, got, want, "create library return mismatch")
}

func deepEqual(t *testing.T, got, want interface{}, message string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s\nhave: %#+v\nwant: %#+v", message, got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
}

func TestLibraryService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries/1644", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", defaultMediaType)
		w.WriteHeader(202)
	})

	resp, err := client.Library.Delete("1644")
	if err != nil {
		t.Errorf("Library.Delete returned error: %v", err)
	}

	if got, want := resp.StatusCode, 202; got != want {
		t.Errorf("Library.Delete response code = %d, want %d", got, want)
	}
}

func TestLibraryService_Delete_404(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/"+defaultAPIVersion+"library-entries/1644", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", defaultMediaType)
		w.WriteHeader(404)
		w.Header().Set("Content-Type", defaultMediaType)
		json.NewEncoder(w).Encode(&struct{ Errors []Error }{[]Error{{Title: "Record not found", Code: "404", Status: "404"}}})
	})

	resp, err := client.Library.Delete("1644")
	if err == nil {
		t.Error("Library.Delete for 404 expected to return error")
	}

	if _, ok := err.(*ErrorResponse); !ok {
		t.Error("expected error to be *ErrorResponse")
	}

	if got, want := resp.StatusCode, 404; got != want {
		t.Errorf("Library.Delete response code = %d, want %d", got, want)
	}
}
