// +build integration

package kitsu_test

import (
	"testing"

	"github.com/nstratos/go-kitsu/kitsu"
)

const (
	results         = 5
	firstResultSlug = "cowboy-bebop"
)

var client *kitsu.Client

func setup(t *testing.T) {
	// Create kitsu client for tests.
	client = kitsu.NewClient(nil)
}

func TestAnimeServiceIntegration(t *testing.T) {
	setup(t)

	// Get anime list with options to include specific limit and includes.
	opt := &kitsu.Options{
		PageLimit: results,
		Include:   []string{"castings.character", "castings.person"},
	}
	list, resp, err := client.Anime.List(opt)
	if err != nil {
		t.Fatal("client.Anime.List returned err:", err)
	}

	// Check page offsets in Response.
	if got, want := resp.NextOffset, results; got != want {
		t.Fatalf("client.Anime.List NextOffset = %d, want %d", got, want)
	}
	if got, want := resp.PrevOffset, 0; got != want {
		t.Fatalf("client.Anime.List PrevOffset = %d, want %d", got, want)
	}
	if got, want := resp.FirstOffset, 0; got != want {
		t.Fatalf("client.Anime.List FirstOffset = %d, want %d", got, want)
	}
	if resp.LastOffset == 0 {
		t.Fatalf("client.Anime.List LastOffset must not be 0")
	}

	// Test that the number of results is the same as we asked in the options.
	if len(list) != results {
		t.Fatalf("client.Anime.List results = %d, want %d", len(list), results)
	}

	// Check that all anime include their castings.
	for _, a := range list {
		if a.Castings == nil {
			t.Fatalf("client.Anime.List expected to include castings. %d %s castings is nil", a.ID, a.Slug)
		}
	}

	// Get details for first anime in the list using the same options as
	// before. PageLimit is ignored by the kitsu API this time.
	bebop, _, err := client.Anime.Show(list[0].ID, opt)
	if err != nil {
		t.Fatal("client.Anime.Show returned err:", err)
	}

	// First result in kitsu database is Cowboy Bebop.
	if bebop.Slug != firstResultSlug {
		t.Fatalf("client.Anime.Show first result slug = %s, want %s", bebop.Slug, firstResultSlug)
	}

	// Check that the anime includes castings.
	if bebop.Castings == nil {
		t.Fatalf("client.Anime.Show expected to include castings. %d %s castings is nil", bebop.ID, bebop.Slug)
	}
}