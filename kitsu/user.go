package kitsu

import (
	"fmt"
	"reflect"
)

// UserService handles communication with the user related methods of the
// Kitsu API.
//
// Kitsu API docs:
// http://docs.kitsu17.apiary.io/#reference/users/library/show-user
type UserService service

// User represents a Kitsu user.
type User struct {
	ID             string                 `jsonapi:"primary,users"`
	Name           string                 `jsonapi:"attr,name"`
	About          string                 `jsonapi:"attr,about"`
	LifeSpent      int64                  `jsonapi:"attr,lifeSpentOnAnime"`
	Avatar         map[string]interface{} `jsonapi:"attr,avatar"`
	Waifu          *Character             `jsonapi:"relation,waifu"`
	LibraryEntries []*LibraryEntry        `jsonapi:"relation,libraryEntries"`
}

// Show returns details for a specific User by providing the ID of the user
// e.g. 29745.
func (s *UserService) Show(userID string, opts ...URLOption) (*User, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"users/%s", userID)

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}
	return usr, resp, nil
}

// List returns a list of Users. Optional parameters can be specified to filter
// the search results and control pagination, sorting etc.
func (s *UserService) List(opts ...URLOption) ([]*User, *Response, error) {
	u := defaultAPIVersion + "users"

	req, err := s.client.NewRequest("GET", u, nil, opts...)
	if err != nil {
		return nil, nil, err
	}

	userType := reflect.TypeOf(&User{})
	data, resp, err := s.client.DoMany(req, userType)
	if err != nil {
		return nil, resp, err
	}

	users := make([]*User, 0, len(data))
	for _, d := range data {
		a, ok := d.(*User)
		if !ok {
			// This should never happen.
			return nil, resp, fmt.Errorf("expected user type %v but it was %T", userType, a)
		}
		users = append(users, a)
	}

	return users, resp, nil
}
