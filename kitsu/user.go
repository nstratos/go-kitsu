package kitsu

import (
	"fmt"
)

// UserService handles communication with the user related methods of the
// Kitsu API.
//
// Kitsu API docs:
// http://docs.kitsu17.apiary.io/#reference/users/library/show-user
type UserService service

// User represents a Kitsu user.
type User struct {
	ID        string `jsonapi:"primary,users"`
	Name      string `jsonapi:"attr,name,omitempty"`
	About     string `jsonapi:"attr,about,omitempty"`
	LifeSpent int64  `jsonapi:"attr,lifeSpentOnAnime,omitempty"`
	//Avatar    map[string]interface{} `jsonapi:"attr,avatar,omitempty"`
	//Avatar         Avatar          `jsonapi:"attr,avatar,omitempty"`
	Waifu          *Character      `jsonapi:"relation,waifu,omitempty"`
	LibraryEntries []*LibraryEntry `jsonapi:"relation,libraryEntries,omitempty"`
}

type Avatar struct {
	Original string
}

// Show returns details for a specific User by providing the ID of the user
// e.g. 29745.
func (s *UserService) Show(userID string, opts ...URLOption) (*User, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"users/%s", userID)

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// List returns a list of Users. Optional parameters can be specified to filter
// the search results and control pagination, sorting etc.
func (s *UserService) List(opts ...URLOption) ([]*User, *Response, error) {
	u := defaultAPIVersion + "users"

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}
