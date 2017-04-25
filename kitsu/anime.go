package kitsu

import (
	"fmt"
)

// The possible anime show types. They are convenient for making comparisons
// with Anime.ShowType.
const (
	AnimeTypeTV      = "TV"
	AnimeTypeSpecial = "special"
	AnimeTypeOVA     = "OVA"
	AnimeTypeONA     = "ONA"
	AnimeTypeMovie   = "movie"
	AnimeTypeMusic   = "music"
)

// AnimeService handles communication with the anime related methods of the
// Kitsu API.
//
// Kitsu API docs:
// http://docs.kitsu17.apiary.io/#reference/media/anime/show-anime
type AnimeService service

// Anime represents a Kitsu anime.
type Anime struct {
	ID       string     `jsonapi:"primary,anime"`
	Slug     string     `jsonapi:"attr,slug,omitempty"`     // Unique slug used for page URLs, e.g. attack-on-titan.
	ShowType string     `jsonapi:"attr,showType,omitempty"` // Show format of the anime. Can be compared with AnimeType constants.
	Genres   []*Genre   `jsonapi:"relation,genres,omitempty"`
	Castings []*Casting `jsonapi:"relation,castings,omitempty"`
}

// Genre represents a Kitsu media genre. Genre is a relationship of Kitsu media
// types like Anime, Manga and Drama.
type Genre struct {
	ID          string `jsonapi:"primary,genres"`
	Name        string `jsonapi:"attr,name"`
	Slug        string `jsonapi:"attr,slug"`
	Description string `jsonapi:"attr,description"`
}

// Casting represents a Kitsu media casting. Casting is a relationship of Kitsu
// media types like Anime, Manga and Drama.
type Casting struct {
	ID         string     `jsonapi:"primary,castings"`
	Role       string     `jsonapi:"attr,role"`
	VoiceActor bool       `jsonapi:"attr,voiceActor"`
	Featured   bool       `jsonapi:"attr,featured"`
	Language   string     `jsonapi:"attr,language"`
	Character  *Character `jsonapi:"relation,character"`
	Person     *Person    `jsonapi:"relation,person"`
}

// BUG(google/jsonapi): Unmarshaling of fields which are of type struct or
// map[string]string is not supported by google/jsonapi. A workaround for
// fields such as Character.Image and User.Avatar is to use
// map[string]interface{} instead.
//
// See: https://github.com/google/jsonapi/issues/74
//
// Another limitation is being unable to unmarshal to custom types such as
// "enum" types like AnimeType, MangaType and LibraryEntryStatus. These are
// useful for doing comparisons and working with fields such as Anime.ShowType,
// Manga.ShowType and LibraryEntry.Status.
//
// Because of this limitation the string type is used for those fields instead.
// As such, instead of using those custom types, we keep the possible values as
// untyped string constants to avoid unnecessary conversions when working with
// those fields.

// Character represents a Kitsu character like the fictional characters that
// appear in anime, manga and drama. Character is a relationship of Casting.
type Character struct {
	ID          string                 `jsonapi:"primary,characters"`
	Slug        string                 `jsonapi:"attr,slug"`
	Name        string                 `jsonapi:"attr,name"`
	MALID       int                    `jsonapi:"attr,malId"`
	Description string                 `jsonapi:"attr,description"`
	Image       map[string]interface{} `jsonapi:"attr,image"`
}

// Person represents a person that is involved with a certain media. It can be
// voice actors, animators, etc. Person is a relationship of Casting.
type Person struct {
	ID    string `jsonapi:"primary,people"`
	Name  string `jsonapi:"attr,name"`
	MALID int    `jsonapi:"attr,malId"`
	Image string `jsonapi:"attr,image"`
}

// Show returns details for a specific Anime by providing a unique identifier
// of the anime e.g. 7442.
func (s *AnimeService) Show(animeID string, opts ...URLOption) (*Anime, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"anime/%s", animeID)

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	a := new(Anime)
	resp, err := s.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// List returns a list of Anime. Optional parameters can be specified to filter
// the search results and control pagination, sorting etc.
func (s *AnimeService) List(opts ...URLOption) ([]*Anime, *Response, error) {
	u := defaultAPIVersion + "anime"

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	var anime []*Anime
	resp, err := s.client.Do(req, &anime)
	if err != nil {
		return nil, resp, err
	}

	return anime, resp, nil
}
