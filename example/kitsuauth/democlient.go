package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/nstratos/go-kitsu/kitsu"
)

// demoClient has methods showcasing the usage of the different Kitsu API
// methods. It stores the first error it encounters so error checking only needs
// to be done once.
//
// This pattern is used for convenience and should not be used in concurrent
// code without guarding the error.
type demoClient struct {
	*kitsu.Client
	err error
}

func (c *demoClient) showcase(ctx context.Context) error {
	methods := []func(context.Context){
		// Uncomment the methods you need to see their results. Run or build
		// using -tags=debug to see the full HTTP request and response.
		c.animeList,
	}
	for _, m := range methods {
		m(ctx)
	}
	if c.err != nil {
		return c.err
	}
	return nil
}

func (c *demoClient) animeList(ctx context.Context) {
	if c.err != nil {
		return
	}
	const results = 5

	// Get anime list with options to include specific limit and includes.
	anime, _, err := c.Anime.List(
		kitsu.Limit(results),
		kitsu.Include("genres"),
	)
	if err != nil {
		c.err = err
		return
	}

	for _, a := range anime {
		fmt.Printf("ID: %5s, Rank: %5d, Popularity: %5d %25s (%v)", a.ID, a.RatingRank, a.PopularityRank, a.Titles["en"], a.StartDate)
		genres := make([]string, 0, len(a.Genres))
		for _, g := range a.Genres {
			genres = append(genres, g.Name)
		}
		fmt.Printf(" [%55s]\n", strings.Join(genres, ", "))
	}
}
