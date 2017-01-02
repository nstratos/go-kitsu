package kitsu

import (
	"fmt"
	"net/http"
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
	Data *Anime `json:"data"`
}

// Anime represents a Kitsu Anime.
type Anime struct {
	Resource
	Attributes AnimeAttributes `json:"attributes"`
}

// AnimeAttributes represent the attributes of an Anime object.
type AnimeAttributes struct {
	Slug string `json:"slug,omitempty"` // Unique slug used for page URLs e.g. attack-on-titan.
}

// Show returns details for a specific Anime by providing a unique identifier
// of the anime e.g. 7442.
func (s *AnimeService) Show(animeID string) (*AnimeShowResponse, *http.Response, error) {
	urlStr := fmt.Sprintf(defaultAPIVersion+"anime/%s", animeID)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	anime := new(AnimeShowResponse)
	resp, err := s.client.Do(req, anime)
	if err != nil {
		return nil, resp, err
	}
	return anime, resp, nil
}
