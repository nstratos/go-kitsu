package kitsu

import (
	"fmt"
	"reflect"
)

type AnimeType string

const (
	AnimeTypeTV      AnimeType = "TV"
	AnimeTypeSpecial AnimeType = "special"
	AnimeTypeOVA     AnimeType = "OVA"
	AnimeTypeONA     AnimeType = "ONA"
	AnimeTypeMovie   AnimeType = "movie"
	AnimeTypeMusic   AnimeType = "music"
)

type MangaType string

const (
	MangaTypeDrama   MangaType = "drama"
	MangaTypeNovel   MangaType = "novel"
	MangaTypeManhua  MangaType = "manhua"
	MangaTypeOneshot MangaType = "oneshot"
	MangaTypeDoujin  MangaType = "doujin"
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
	Slug     string     `jsonapi:"attr,slug"`
	Genres   []*Genre   `jsonapi:"relation,genres"`
	Castings []*Casting `jsonapi:"relation,castings"`
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

	animeType := reflect.TypeOf(&Anime{})
	data, resp, err := s.client.DoMany(req, animeType)
	if err != nil {
		return nil, resp, err
	}

	anime := make([]*Anime, 0, len(data))
	for _, d := range data {
		a, ok := d.(*Anime)
		if !ok {
			// This should never happen.
			return nil, resp, fmt.Errorf("expected anime type %v but it was %T", animeType, a)
		}
		anime = append(anime, a)
	}

	return anime, resp, nil
}
