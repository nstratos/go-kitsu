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

type Anime struct {
	ID       string     `jsonapi:"primary,anime"`
	Slug     string     `jsonapi:"attr,slug"`
	Genres   []*Genre   `jsonapi:"relation,genres"`
	Castings []*Casting `jsonapi:"relation,castings"`
}

type Genre struct {
	ID          string `jsonapi:"primary,genres"`
	Name        string `jsonapi:"attr,name"`
	Slug        string `jsonapi:"attr,slug"`
	Description string `jsonapi:"attr,description"`
}

type Casting struct {
	ID         string     `jsonapi:"primary,castings"`
	Role       string     `jsonapi:"attr,role"`
	VoiceActor bool       `jsonapi:"attr,voiceActor"`
	Featured   bool       `jsonapi:"attr,featured"`
	Language   string     `jsonapi:"attr,language"`
	Character  *Character `jsonapi:"relation,character"`
	Person     *Person    `jsonapi:"relation,person"`
}

type Character struct {
	ID          string `jsonapi:"primary,characters"`
	Slug        string `jsonapi:"attr,slug"`
	Name        string `jsonapi:"attr,name"`
	MALID       int    `jsonapi:"attr,malId"`
	Description string `jsonapi:"attr,description"`
	//Image       CharacterImage `json:"image" jsonapi:"attr,image,omitempty"`
}

// BUG(google/jsonapi): Apparently unmarshaling struct fields does not work
// yet. See https://github.com/google/jsonapi/issues/74
//
//type CharacterImage struct {
//	Original string `json:"original" jsonapi:"attr,original,omitempty"`
//}

type Person struct {
	ID    string `jsonapi:"primary,people"`
	Name  string `jsonapi:"attr,name"`
	MALID int    `jsonapi:"attr,malId"`
	Image string `jsonapi:"attr,image"`
}

// Show returns details for a specific Anime by providing a unique identifier
// of the anime e.g. 7442.
func (s *AnimeService) Show(animeID string) (*Anime, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"anime/%s", animeID)

	req, err := s.client.NewRequest("GET", u, nil)
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
func (s *AnimeService) List(opt *Options) ([]*Anime, *Response, error) {
	u := defaultAPIVersion + "anime"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
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
			return nil, resp, fmt.Errorf("expected anime type %v but it was %T", animeType, a)
		}
		anime = append(anime, a)
	}

	return anime, resp, nil
}
