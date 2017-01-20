package kitsu

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// AnimeShowResponse is the response returned by AnimeService.Show which
// contains one Anime.
type AnimeShowResponse struct {
	Data *AnimeData `json:"data,omitempty"`
}

type AnimeData struct {
	Resource
	Attributes *AnimeAttributes `json:"attributes,omitempty"`
}

func toAnime(id string, attr *AnimeAttributes) *Anime {
	if attr == nil {
		return &Anime{ID: id}
	}
	a := &Anime{
		ID:   id,
		Slug: attr.Slug,
	}
	return a
}

// AnimeAttributes represent the attributes of an Anime object.
type AnimeAttributes struct {
	Slug string `json:"slug,omitempty"` // Unique slug used for page URLs e.g. attack-on-titan.
}

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
	ID          string         `jsonapi:"primary,characters"`
	Slug        string         `jsonapi:"attr,slug"`
	Name        string         `jsonapi:"attr,name"`
	MALID       int            `jsonapi:"attr,malId"`
	Description string         `jsonapi:"attr,description"`
	Image       CharacterImage `jsonapi:"attr,image"`
}

type CharacterImage struct {
	Original string `jsonapi:"attr,original"`
}

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

	asr := new(AnimeShowResponse)
	resp, err := s.client.Do(req, asr)
	if err != nil {
		return nil, newResponse(resp), err
	}
	return animeFromShowResponse(asr), newResponse(resp), nil
}

func animeFromShowResponse(asr *AnimeShowResponse) *Anime {
	if asr == nil || asr.Data == nil {
		return nil
	}
	return toAnime(asr.Data.ID, asr.Data.Attributes)
}

// AnimeListResponse is the response returned by AnimeService.List which
// contains many Anime.
type AnimeListResponse struct {
	Data  []*AnimeData    `json:"data"`
	Links PaginationLinks `json:"links"`
}

type PaginationLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
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

	alr := new(AnimeListResponse)
	resp, err := s.client.Do(req, alr)
	if err != nil {
		return nil, newResponse(resp), err
	}

	return returnAnimeListResponse(alr, resp)
}

func animeFromListResponse(alr *AnimeListResponse) []*Anime {
	var anime []*Anime
	if alr != nil && alr.Data != nil {
		anime = make([]*Anime, 0, len(alr.Data))
		for _, d := range alr.Data {
			a := toAnime(d.ID, d.Attributes)
			anime = append(anime, a)
		}
	}
	return anime
}

func returnAnimeListResponse(alr *AnimeListResponse, r *http.Response) ([]*Anime, *Response, error) {
	var anime []*Anime
	var resp = newResponse(r)
	if alr != nil && alr.Data != nil {
		anime = make([]*Anime, 0, len(alr.Data))
		for _, d := range alr.Data {
			a := toAnime(d.ID, d.Attributes)
			anime = append(anime, a)
		}

		firstOffset, err := parseOffsetFromLink(alr.Links.First)
		if err != nil {
			return anime, resp, fmt.Errorf("failed to parse first link: %v", err)
		}
		resp.FirstOffset = firstOffset

		lastOffset, err := parseOffsetFromLink(alr.Links.Last)
		if err != nil {
			return anime, resp, fmt.Errorf("failed to parse last link: %v", err)
		}
		resp.LastOffset = lastOffset

		prevOffset, err := parseOffsetFromLink(alr.Links.Prev)
		if err != nil {
			return anime, resp, fmt.Errorf("failed to parse prev link: %v", err)
		}
		resp.PrevOffset = prevOffset

		nextOffset, err := parseOffsetFromLink(alr.Links.Next)
		if err != nil {
			return anime, resp, fmt.Errorf("failed to parse next link: %v", err)
		}
		resp.NextOffset = nextOffset
	}
	return anime, resp, nil
}

func parseOffsetFromLink(link string) (int, error) {
	var offset int
	u, err := url.Parse(link)
	if err != nil {
		return offset, err
	}
	v := u.Query()
	s := v.Get("page[offset]")
	if s == "" {
		return offset, nil
	}
	return strconv.Atoi(s)
}
