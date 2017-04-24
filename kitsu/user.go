package kitsu

import (
	"fmt"
	"io"
	"reflect"

	"github.com/nstratos/go-kitsu/kitsu/internal/jsonapi"
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

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	user, err := decodeUser(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return user, resp, nil
}

func decodeUser(r io.Reader) (*User, error) {
	u := new(User)
	err := jsonapi.DecodeOne(r, u)
	return u, err
}

// List returns a list of Users. Optional parameters can be specified to filter
// the search results and control pagination, sorting etc.
func (s *UserService) List(opts ...URLOption) ([]*User, *Response, error) {
	u := defaultAPIVersion + "users"

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	users, o, err := decodeUserList(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	resp.Offset = o

	return users, resp, nil
}

func decodeUserList(r io.Reader) ([]*User, PageOffset, error) {
	data, o, err := jsonapi.DecodeMany(r, reflect.TypeOf(&User{}))
	if err != nil {
		return nil, PageOffset{}, err
	}

	users := make([]*User, 0, len(data))
	for _, d := range data {
		if a, ok := d.(*User); ok {
			users = append(users, a)
		}
	}

	return users, makePageOffset(o), nil
}
