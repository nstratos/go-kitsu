package kitsu

import (
	"fmt"
	"io"
	"reflect"

	"github.com/nstratos/go-kitsu/kitsu/internal/jsonapi"
)

// The possible library entry statuses. They are convenient when creating a
// LibraryEntry or for making comparisons with LibraryEntry.Status.
const (
	LibraryEntryStatusCurrent   = "current"
	LibraryEntryStatusPlanned   = "planned"
	LibraryEntryStatusCompleted = "completed"
	LibraryEntryStatusOnHold    = "on_hold"
	LibraryEntryStatusDropped   = "dropped"
)

// LibraryService handles communication with the library entry related methods
// of the Kitsu API.
type LibraryService service

// LibraryEntry represents a Kitsu user's library entry.
type LibraryEntry struct {
	ID             string `jsonapi:"primary,libraryEntries"`
	Status         string `jsonapi:"attr,status,omitempty"`         // Status for related media. Can be compared with LibraryEntryStatus constants.
	Progress       int    `jsonapi:"attr,progress,omitempty"`       // How many episodes/chapters have been consumed, e.g. 22.
	Reconsuming    bool   `jsonapi:"attr,reconsuming,omitempty"`    // Whether the media is being reconsumed, e.g. false.
	ReconsumeCount int    `jsonapi:"attr,reconsumeCount,omitempty"` // How many times the media has been reconsumed, e.g. 0.
	Notes          string `jsonapi:"attr,notes,omitempty"`          // Note attached to this entry, e.g. Very Interesting!
	Private        bool   `jsonapi:"attr,private,omitempty"`        // Whether this entry is hidden from the public, e.g. false.
	Rating         string `jsonapi:"attr,rating,omitempty"`         // User rating out of 5.0.
	UpdatedAt      string `jsonapi:"attr,updatedAt,omitempty"`      // When the entry was last updated, e.g. 2016-11-12T03:35:00.064Z.
	User           *User  `jsonapi:"relation,user,omitempty"`
	Anime          *Anime `jsonapi:"relation,anime,omitempty"`
	//Media          *Media `jsonapi:"relation,media"`
	Media interface{} `jsonapi:"relation,media,omitempty"`
}

//type Media struct {
//	ID   string `jsonapi:"primary,media"`
//	Type string `jsonapi:"attr,type"`
//}

// Show returns details for a specific LibraryEntry by providing a unique identifier
// of the library entry, e.g. 5269457.
func (s *LibraryService) Show(libraryEntryID string, opts ...URLOption) (*LibraryEntry, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"library-entries/%s", libraryEntryID)

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	e, err := decodeLibraryEntry(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return e, resp, nil
}

func decodeLibraryEntry(r io.Reader) (*LibraryEntry, error) {
	e := new(LibraryEntry)
	err := jsonapi.DecodeOne(r, e)
	return e, err
}

func (s *LibraryService) Create(e *LibraryEntry, opts ...URLOption) ([]*LibraryEntry, *Response, error) {
	u := defaultAPIVersion + "library-entries"

	req, err := s.client.NewRequest("POST", u, e, opts...)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	entries, o, err := decodeLibraryEntryList(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	resp.Offset = o

	return entries, resp, nil
}

func decodeLibraryEntryList(r io.Reader) ([]*LibraryEntry, PageOffset, error) {
	data, o, err := jsonapi.DecodeMany(r, reflect.TypeOf(&LibraryEntry{}))
	if err != nil {
		return nil, PageOffset{}, err
	}

	entries := make([]*LibraryEntry, 0, len(data))
	for _, d := range data {
		if a, ok := d.(*LibraryEntry); ok {
			entries = append(entries, a)
		}
	}

	return entries, makePageOffset(o), nil
}
