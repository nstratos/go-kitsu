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
//
// Additional filters: text, season, streamers
type Anime struct {
	ID                  string  `jsonapi:"primary,anime"`
	Slug                string  `jsonapi:"attr,slug,omitempty"`                // Unique slug used for page URLs, e.g. attack-on-titan.
	Synopsis            string  `jsonapi:"attr,synopsis,omitempty"`            // Synopsis of the anime, e.g. Several hundred years ago, humans were...
	CoverImageTopOffset int     `jsonapi:"attr,coverImageTopOffset,omitempty"` // e.g. 263
	CanonicalTitle      string  `jsonapi:"attr,canonical_title,omitempty"`     // Canonical title for the anime, e.g. Attack on Titan
	AverageRating       float64 `jsonapi:"attr,averageRating,omitempty"`       // The average of all user ratings for the anime, e.g. 4.26984658306698
	StartDate           string  `jsonapi:"attr,startDate,omitempty"`           // Date the anime started airing/was released, e.g. 2013-04-07
	EndDate             string  `jsonapi:"attr,endDate,omitempty"`             // Date the anime finished airing, e.g. 2013-09-28
	EpisodeCount        int     `jsonapi:"attr,episodeCount,omitempty"`        // How many episodes the anime has, e.g. 25
	EpisodeLength       int     `jsonapi:"attr,episodeLength,omitempty"`       // How many minutes long each episode is, e.g. 24
	ShowType            string  `jsonapi:"attr,showType,omitempty"`            // Show format of the anime. Can be compared with AnimeType constants.
	YoutubeVideoID      string  `jsonapi:"attr,youtubeVideoId,omitempty"`      // YouTube video id for Promotional Video, e.g. n4Nj6Y_SNYI
	AgeRating           string  `jsonapi:"attr,ageRating,omitempty"`           // Age rating for the anime, e.g. R
	AgeRatingGuide      string  `jsonapi:"attr,ageRatingGuide,omitempty"`      // Description of the age rating, e.g. Violence, Profanity

	// The titles of the anime which include:
	// English title of the anime, e.g. "en": "Attack on Titan"
	// The romaji title of the anime, e.g. "en_jp": "Shingeki no Kyojin"
	// Japanese title of the anime, e.g.  "ja_jp": "進撃の巨人"
	Titles map[string]interface{} `jsonapi:"attr,titles,omitempty"`

	// Shortened nicknames for the anime.
	AbbreviatedTitles []string `jsonapi:"attr,abbreviatedTitles,omitempty"`

	// The URL template for the poster, e.g. "original": "https://static.hummingbird.me/anime/7442/poster/$1.png"
	PosterImage map[string]interface{} `jsonapi:"attr,posterImage,omitempty"`

	// The URL template for the cover, e.g. "original": "https://static.hummingbird.me/anime/7442/cover/$1.png"
	CoverImage map[string]interface{} `jsonapi:"attr,coverImage,omitempty"`

	// How many times each rating has been given to the anime, e.g.
	// "0.5": "114",
	// "1.0": "279",
	// "1.5": "146",
	// "2.0": "359",
	// "2.5": "763",
	// "3.0": "2331",
	// "3.5": "3034",
	// "4.0": "5619",
	// "4.5": "5951",
	// "5.0": "12878"
	RatingFrequencies map[string]interface{} `jsonapi:"attr,ratingFrequencies,omitempty"`

	// Relationships.

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
