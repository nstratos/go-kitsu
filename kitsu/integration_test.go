// +build integration

package kitsu_test

import (
	"testing"

	"github.com/nstratos/go-kitsu/kitsu"
)

// setup creates a new Kitsu client for tests.
func setup(t *testing.T) *kitsu.Client {
	return kitsu.NewClient(nil)
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
