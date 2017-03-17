package kitsu

type LibraryEntryStatus string

const (
	LibraryEntryStatusCurrent   LibraryEntryStatus = "current"
	LibraryEntryStatusPlanned   LibraryEntryStatus = "planned"
	LibraryEntryStatusCompleted LibraryEntryStatus = "completed"
	LibraryEntryStatusOnHold    LibraryEntryStatus = "on_hold"
	LibraryEntryStatusDropped   LibraryEntryStatus = "dropped"
)

type LibraryEntry struct {
	// Unmarshaling of enum types likes LibraryEntryStatus is not supported by
	// google/jsonapi so we are using string instead.
	Status         string `jsonapi:"attr,status"`         // Status for related media. One of: current, planned, completed, on_hold, dropped.
	Progress       int    `jsonapi:"attr,progress"`       // How many episodes/chapters have been consumed, e.g. 22.
	Reconsuming    bool   `jsonapi:"attr,reconsuming"`    // Whether the media is being reconsumed, e.g. false.
	ReconsumeCount int    `jsonapi:"attr,reconsumeCount"` // How many times the media has been reconsumed, e.g. 0.
	Notes          string `jsonapi:"attr,notes"`          // Note attached to this entry, e.g. Very Interesting!
	Private        bool   `jsonapi:"attr,private"`        // Whether this entry is hidden from the public, e.g. false.
	Rating         string `jsonapi:"attr,rating"`         // User rating out of 5.0.
	UpdatedAt      string `jsonapi:"attr,updatedAt"`      // When the entry was last updated, e.g. 2016-11-12T03:35:00.064Z.
}
