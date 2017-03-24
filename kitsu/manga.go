package kitsu

// MangaType represents all the possible manga show types.
type MangaType string

// The possible manga show types. They are convenient for making comparisons
// with Manga.ShowType.
const (
	MangaTypeDrama   MangaType = "drama"
	MangaTypeNovel   MangaType = "novel"
	MangaTypeManhua  MangaType = "manhua"
	MangaTypeOneshot MangaType = "oneshot"
	MangaTypeDoujin  MangaType = "doujin"
)
