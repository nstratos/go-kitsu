package kitsu

import (
	"fmt"
)

// The possible age rating values for media types like Anime, Manga and Drama.
const (
	AgeRatingG   = "G"   // General Audiences
	AgeRatingPG  = "PG"  // Parental Guidance Suggested
	AgeRatingR   = "R"   // Restricted
	AgeRatingR18 = "R18" // Explicit
)

// Possible values for Anime.Status.
const (
	AnimeStatusCurrent    = "current"
	AnimeStatusFinished   = "finished"
	AnimeStatusTBA        = "tba"
	AnimeStatusUnreleased = "unreleased"
	AnimeStatusUpcoming   = "upcoming"
)

// The possible anime subtypes. They are convenient for making comparisons
// with Anime.Subtype.
const (
	AnimeSubtypeONA     = "ONA"
	AnimeSubtypeOVA     = "OVA"
	AnimeSubtypeTV      = "TV"
	AnimeSubtypeMovie   = "movie"
	AnimeSubtypeMusic   = "music"
	AnimeSubtypeSpecial = "special"
)

// AnimeService handles communication with the anime related methods of the
// Kitsu API.
//
// Kitsu API docs:
// http://docs.kitsu.apiary.io/#reference/media/anime
type AnimeService service

// Anime represents a Kitsu anime.
//
// Additional filters: text, season, streamers
type Anime struct {
	ID string `jsonapi:"primary,anime"`

	// --- Attributes ---

	// ISO 8601 date and time, e.g. 2017-07-27T22:21:26.824Z
	CreatedAt string `jsonapi:"attr,createdAt,omitempty"`

	// ISO 8601 of last modification, e.g. 2017-07-27T22:47:45.129Z
	UpdatedAt string `jsonapi:"attr,updatedAt,omitempty"`

	// Unique slug used for page URLs, e.g. cowboy-bebop
	Slug string `jsonapi:"attr,slug,omitempty"`

	// Synopsis of the anime, e.g.
	//
	// In the year 2071, humanity has colonoized several of the planets and
	// moons...
	Synopsis string `jsonapi:"attr,synopsis,omitempty"`

	// e.g. 400
	CoverImageTopOffset int `jsonapi:"attr,coverImageTopOffset,omitempty"`

	// Titles in different languages. Other languages will be listed if they
	// exist, e.g.
	//
	// "en": "Attack on Titan"
	//
	// "en_jp": "Shingeki no Kyojin"
	//
	// "ja_jp": "進撃の巨人"
	Titles map[string]interface{} `jsonapi:"attr,titles,omitempty"`

	// Canonical title for the anime, e.g. Attack on Titan
	CanonicalTitle string `jsonapi:"attr,canonicalTitle,omitempty"`

	// Shortened nicknames for the anime, e.g. COWBOY BEBOP
	AbbreviatedTitles []string `jsonapi:"attr,abbreviatedTitles,omitempty"`

	// The average of all user ratings for the anime, e.g. 88.65
	AverageRating string `jsonapi:"attr,averageRating,omitempty"`

	// How many times each rating has been given to the anime, e.g.
	//
	// "2": "72"
	//
	// "3": "0"
	//
	// ...
	//
	// "19": "40"
	//
	// "20": "13607"
	RatingFrequencies map[string]interface{} `jsonapi:"attr,ratingFrequencies,omitempty"`

	// e.g. 40405
	UserCount int `jsonapi:"attr,userCount,omitempty"`

	// e.g. 3277
	FavoritesCount int `jsonapi:"attr,favoritesCount,omitempty"`

	// Date the anime started airing/was released, e.g. 2013-04-07
	StartDate string `jsonapi:"attr,startDate,omitempty"`

	// Date the anime finished airing, e.g. 2013-09-28
	EndDate string `jsonapi:"attr,endDate,omitempty"`

	// e.g. 10
	PopularityRank int `jsonapi:"attr,popularityRank,omitempty"`

	// e.g. 10
	RatingRank int `jsonapi:"attr,ratingRank,omitempty"`

	// Possible values described by the AgeRating constants.
	AgeRating string `jsonapi:"attr,ageRating,omitempty"`

	// Description of the age rating, e.g. 17+ (violence & profanity)
	AgeRatingGuide string `jsonapi:"attr,ageRatingGuide,omitempty"`

	// Show format of the anime. Possible values described by the AnimeSubtype
	// constants.
	Subtype string `jsonapi:"attr,subtype,omitempty"`

	// Possible values described by the AnimeStatus constants.
	Status string `jsonapi:"attr,status,omitempty"`

	// The URL template for the poster, e.g.
	//
	// "tiny": "https://media.kitsu.io/anime/poster_images/1/tiny.jpg?1431697256"
	//
	// "small": "https://media.kitsu.io/anime/poster_images/1/small.jpg?1431697256"
	//
	// "medium": "https://media.kitsu.io/anime/poster_images/1/medium.jpg?1431697256"
	//
	// "large": "https://media.kitsu.io/anime/poster_images/1/large.jpg?1431697256"
	//
	// "original: "https://media.kitsu.io/anime/poster_images/1/original.jpg?1431697256"
	PosterImage map[string]interface{} `jsonapi:"attr,posterImage,omitempty"`

	// The URL template for the cover, e.g.
	//
	// "tiny": "https://media.kitsu.io/anime/cover_images/1/tiny.jpg?1416336000"
	//
	// "small": "https://media.kitsu.io/anime/cover_images/1/small.jpg?1416336000"
	//
	// "large": "https://media.kitsu.io/anime/cover_images/1/large.jpg?1416336000"
	//
	// "original": "https://media.kitsu.io/anime/cover_images/1/original.jpg?1416336000"
	CoverImage map[string]interface{} `jsonapi:"attr,coverImage,omitempty"`

	// How many episodes the anime has, e.g. 25
	EpisodeCount int `jsonapi:"attr,episodeCount,omitempty"`

	// How many minutes long each episode is, e.g. 24
	EpisodeLength int `jsonapi:"attr,episodeLength,omitempty"`

	// YouTube video id for Promotional Video, e.g. n4Nj6Y_SNYI
	YoutubeVideoID string `jsonapi:"attr,youtubeVideoId,omitempty"`

	// --- Relationships ---

	Genres     []*Genre          `jsonapi:"relation,genres,omitempty"`
	Mappings   []*Mapping        `jsonapi:"relation,mappings,omitempty"`
	Staff      []*AnimeStaff     `jsonapi:"relation,animeStaff,omitempty"`
	Characters []*AnimeCharacter `jsonapi:"relation,animeCharacters,omitempty"`

	// Deprecated: Use Staff instead.
	Castings []*Casting `jsonapi:"relation,castings,omitempty"`
}

// Genre represents a Kitsu media genre. Genre is a relationship of Kitsu media
// types like Anime, Manga and Drama.
type Genre struct {
	ID          string `jsonapi:"primary,genres"`
	Name        string `jsonapi:"attr,name"`
	Slug        string `jsonapi:"attr,slug"`
	Description string `jsonapi:"attr,description"`
	CreatedAt   string `jsonapi:"attr,createdAt,omitempty"`
	UpdatedAt   string `jsonapi:"attr,updatedAt,omitempty"`
}

// Casting represents a Kitsu media casting. Casting is a relationship of Kitsu
// media types like Anime, Manga and Drama.
//
// Deprecated: Use relationship animeStaff.person instead. At the time of
// writing, the Kitsu API returns 500 when the client asks for castings.
// According to their documentation, castings is being replaced by Anime Staff.
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

// AnimeCharacter represents the characters of an Anime entry.
type AnimeCharacter struct {
	ID   string `jsonapi:"primary,animeCharacters"`
	Role string `jsonapi:"attr,role"`

	Character *Character `jsonapi:"relation,character,omitempty"`
}

// Character represents a Kitsu character like the fictional characters that
// appear in anime, manga and drama. Character is a relationship of Casting.
type Character struct {
	ID          string                 `jsonapi:"primary,characters"`
	Slug        string                 `jsonapi:"attr,slug"`
	Name        string                 `jsonapi:"attr,name"`
	MALID       float64                `jsonapi:"attr,malId"`
	Description string                 `jsonapi:"attr,description"`
	Image       map[string]interface{} `jsonapi:"attr,image"`
}

// AnimeStaff represents the staff of an Anime entry.
type AnimeStaff struct {
	ID        string `jsonapi:"primary,animeStaff"`
	Role      string `jsonapi:"attr,role,omitempty"`
	CreatedAt string `jsonapi:"attr,createdAt,omitempty"` // ISO 8601 date and time e.g. 2017-07-27T22:21:26.824Z
	UpdatedAt string `jsonapi:"attr,updatedAt,omitempty"` // ISO 8601 of last modification e.g. 2017-07-27T22:47:45.129Z

	Person *Person `jsonapi:"relation,person,omitempty"`
}

// Person represents a person that is involved with a certain media. It can be
// voice actors, animators, etc. Person is a relationship of Casting.
type Person struct {
	ID          string                 `jsonapi:"primary,people"`
	Name        string                 `jsonapi:"attr,name"`
	MALID       string                 `jsonapi:"attr,malId"`
	Description string                 `jsonapi:"attr,description"`
	Image       map[string]interface{} `jsonapi:"attr,image"`
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
