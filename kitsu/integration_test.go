// +build integration

package kitsu_test

import (
	"context"
	"flag"
	"sync"
	"testing"

	"github.com/nstratos/go-kitsu/kitsu"
	"golang.org/x/oauth2"
)

var (
	// testAccountID is set by function setup.
	testAccountID = ""

	// Kitsu supports authentication using email+password or slug+password. In
	// order for a new account to have a slug, it should set the Profile URL
	// field in the profile settings.

	testAccountSlug     = flag.String("slug", "testgopher", "Kitsu test account slug to use for authentication")
	testAccountPassword = flag.String("password", "", "Kitsu test account password to use for authentication")
)

// setup creates a new Kitsu client for tests. It needs a test account for
// authentication.
func setup(t *testing.T) *kitsu.Client {
	if *testAccountPassword == "" {
		t.Errorf("No password provided for account with slug %q.", *testAccountSlug)
		t.Error("These tests are meant to be run with a dedicated test account.")
		t.Fatal("You might want to use: go test -tags=integration -slug '<test account slug>' -password '<test account password>'")
	}

	conf := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://kitsu.io/api/oauth/token",
		},
	}

	ctx := context.Background()
	tok, err := conf.PasswordCredentialsToken(ctx, *testAccountSlug, *testAccountPassword)
	if err != nil {
		t.Fatal("could not get token:", err)
	}

	httpClient := conf.Client(ctx, tok)
	kitsuClient := kitsu.NewClient(httpClient)

	var once sync.Once
	once.Do(func() {
		users, _, err := kitsuClient.User.List(kitsu.Filter("slug", *testAccountSlug))
		if err != nil {
			t.Fatal("searching users by slug failed:", err)
		}
		if len(users) != 1 {
			t.Fatalf("could not find 1 user with slug %q", *testAccountSlug)
		}
		testAccountID = users[0].ID
	})

	return kitsuClient
}

func TestAnimeServiceIntegration(t *testing.T) {
	c := setup(t)
	const results = 5

	// Get anime list with options to include specific limit and includes.
	list, resp, err := c.Anime.List(
		kitsu.Limit(results),
		kitsu.Include("castings.character", "castings.person"),
	)
	if err != nil {
		t.Fatal("client.Anime.List returned err:", err)
	}

	// Check page offsets in Response.
	if got, want := resp.Offset.Next, results; got != want {
		t.Fatalf("client.Anime.List Offset.Next = %d, want %d", got, want)
	}
	if got, want := resp.Offset.Prev, 0; got != want {
		t.Fatalf("client.Anime.List Offset.Prev = %d, want %d", got, want)
	}
	if got, want := resp.Offset.First, 0; got != want {
		t.Fatalf("client.Anime.List Offset.First = %d, want %d", got, want)
	}
	if resp.Offset.Last == 0 {
		t.Fatalf("client.Anime.List Offset.Last must not be 0")
	}

	// Test that the number of results is the same as we asked in the options.
	if len(list) != results {
		t.Fatalf("client.Anime.List results = %d, want %d", len(list), results)
	}

	// Check that all anime include their castings.
	for _, a := range list {
		if a.Castings == nil {
			t.Fatalf("client.Anime.List expected to include castings. %s %s castings is nil", a.ID, a.Slug)
		}
	}

	// Get details for the first anime in the list.
	bebop, _, err := c.Anime.Show(
		list[0].ID,
		kitsu.Include("castings.character", "castings.person"),
	)
	if err != nil {
		t.Fatal("client.Anime.Show returned err:", err)
	}

	// First result in kitsu database is Cowboy Bebop.
	const firstResultSlug = "cowboy-bebop"
	if bebop.Slug != firstResultSlug {
		t.Fatalf("client.Anime.Show first result slug = %s, want %s", bebop.Slug, firstResultSlug)
	}

	// Check that the anime includes castings.
	if bebop.Castings == nil {
		t.Fatalf("client.Anime.Show expected to include castings. %s %s castings is nil", bebop.ID, bebop.Slug)
	}
}

func TestLibraryServiceIntegration(t *testing.T) {
	c := setup(t)

	// Get all library entries for test account.
	entries, _, err := c.Library.List(
		kitsu.Filter("userId", testAccountID),
	)
	if err != nil {
		t.Fatal("client.Library.List returned err:", err)
	}

	// Account must have no entries for easier testing.
	if got, want := len(entries), 0; got != want {
		t.Fatalf("Account %q has %d entries but should have %d.", *testAccountSlug, got, want)
	}

	// Add a new library entry.
	newEntry := &kitsu.LibraryEntry{
		Status:   kitsu.LibraryEntryStatusCurrent,
		Progress: 4,
		Rating:   "0.5",
		User: &kitsu.User{
			ID: testAccountID,
		},
		Anime: &kitsu.Anime{
			ID: "1",
		},
	}

	e, _, err := c.Library.Create(newEntry)
	if err != nil {
		t.Fatal("could not create library:", err)
	}

	// Clean up at the end.
	defer func() {
		if _, derr := c.Library.Delete(e.ID); derr != nil {
			t.Errorf("deleting entry with ID %q returned err: %v", e.ID, derr)
		}
	}()

	// Get all library entries again.
	entries, _, err = c.Library.List(
		kitsu.Filter("userId", testAccountID),
	)
	if err != nil {
		t.Fatal("client.Library.List returned err:", err)
	}

	// Check account has the correct number of entries.
	if got, want := len(entries), 1; got != want {
		t.Fatalf("Account %q has %d entries but should have %d.", *testAccountSlug, got, want)
	}
}
