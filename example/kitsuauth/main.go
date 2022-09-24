package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nstratos/go-kitsu/kitsu"
	"golang.org/x/oauth2"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

// Authentication Documentation:
//
// https://kitsu.docs.apiary.io/#introduction/authentication

func run() error {
	var (
		clientID     = flag.String("username", "", "Kitsu account email or slug to use for authentication")
		clientSecret = flag.String("password", "", "Kitsu account password to use for authentication")
	)
	flag.Parse()

	ctx := context.Background()

	tokenClient, err := authenticate(ctx, *clientID, *clientSecret)
	if err != nil {
		return err
	}

	c := demoClient{
		Client: kitsu.NewClient(tokenClient),
	}

	return c.showcase(ctx)
}

func authenticate(ctx context.Context, username, password string) (*http.Client, error) {
	conf := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://kitsu.io/api/oauth/token",
		},
	}

	oauth2Token, err := loadCachedToken()
	if err == nil {
		refreshedToken, err := conf.TokenSource(ctx, oauth2Token).Token()
		if err == nil && (oauth2Token != refreshedToken) {
			fmt.Println("Caching refreshed oauth2 token...")
			if err := cacheToken(*refreshedToken); err != nil {
				return nil, fmt.Errorf("caching refreshed oauth2 token: %s", err)
			}
			return conf.Client(ctx, refreshedToken), nil
		}
		return conf.Client(ctx, oauth2Token), nil
	}

	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("getting password grant token %w", err)
	}

	fmt.Println("Authentication was successful. Caching oauth2 token...")
	if err := cacheToken(*token); err != nil {
		return nil, fmt.Errorf("caching oauth2 token: %s", err)
	}

	return conf.Client(ctx, token), nil
}

const cacheName = "auth-example-token-cache.txt"

func cacheToken(token oauth2.Token) error {
	b, err := json.MarshalIndent(token, "", "   ")
	if err != nil {
		return fmt.Errorf("marshaling token %s: %v", token, err)
	}
	err = os.WriteFile(cacheName, b, 0644)
	if err != nil {
		return fmt.Errorf("writing token %s to file %q: %v", token, cacheName, err)
	}
	return nil
}

func loadCachedToken() (*oauth2.Token, error) {
	b, err := os.ReadFile(cacheName)
	if err != nil {
		return nil, fmt.Errorf("reading oauth2 token from cache file %q: %v", cacheName, err)
	}
	token := new(oauth2.Token)
	if err := json.Unmarshal(b, token); err != nil {
		return nil, fmt.Errorf("unmarshaling oauth2 token: %v", err)
	}
	return token, nil
}
