package kitsu

import (
	"fmt"
	"net/http"
)

// UserService handles communication with the user related methods of the
// Kitsu API.
//
// Kitsu API docs:
// http://docs.kitsu17.apiary.io/#reference/users/library/show-user
type UserService service

// UserShowResponse is the response returnes by UserService.Show which
// contains one User.
type UserShowResponse struct {
	Data *UserData `json:"data,omitempty"`
}

type UserData struct {
	Resource
	Attributes *UserAttributes `json:"attributes,omitempty"`
}

func toUser(id string, attr *UserAttributes) *User {
	if attr == nil {
		return &User{ID: id}
	}
	u := &User{
		ID:        id,
		Name:      attr.Name,
		About:     attr.About,
		Avatar:    attr.Avatar,
		LifeSpent: attr.LifeSpent,
	}
	return u
}

// UserAttributes represent the attributes of an User object.
type UserAttributes struct {
	Name      string            `json:"name,omitempty"`
	About     string            `json:"about,omitempty"`
	Avatar    map[string]string `json:"avatar,omitempty"`
	LifeSpent int64             `json:"lifeSpentOnAnime"`
}

type User struct {
	ID        string            `jsonapi:"primary,users"`
	Name      string            `jsonapi:"attr,name"`
	About     string            `jsonapi:"attr,about"`
	Avatar    map[string]string `jsonapi:"attr,avatar"`
	LifeSpent int64             `jsonapi:"attr,lifeSpentOnAnime"`
	// Library   []*LibraryEntries `jsonapi:"relation,libraryEntries"`
}

// Show returns details for a specific User by providing the ID of the user
// e.g. 29745
func (s *UserService) Show(userID string) (*User, *Response, error) {
	u := fmt.Sprintf(defaultAPIVersion+"users/%s", userID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	usr := new(UserShowResponse)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, newResponse(resp), err
	}
	return userFromShowResponse(usr), newResponse(resp), nil
}

func userFromShowResponse(usr *UserShowResponse) *User {
	if usr == nil || usr.Data == nil {
		return nil
	}
	return toUser(usr.Data.ID, usr.Data.Attributes)
}

// UserListResponse is the response returned by UserService.List which
// contains many Users.
type UserListResponse struct {
	Data  []*UserData     `json:"data"`
	Links PaginationLinks `json:"links"`
}

// List returns a list of Users. Optional parameters can be specified to filter
// the search results and control pagination, sorting etc.
func (s *UserService) List(opt *Options) ([]*User, *Response, error) {
	u := defaultAPIVersion + "users"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alr := new(UserListResponse)
	resp, err := s.client.Do(req, alr)
	if err != nil {
		return nil, newResponse(resp), err
	}

	return returnUserListResponse(alr, resp)
}

func userFromListResponse(ulr *UserListResponse) []*User {
	var user []*User
	if ulr != nil && ulr.Data != nil {
		user = make([]*User, 0, len(ulr.Data))
		for _, d := range ulr.Data {
			a := toUser(d.ID, d.Attributes)
			user = append(user, a)
		}
	}
	return user
}

func returnUserListResponse(ulr *UserListResponse, r *http.Response) ([]*User, *Response, error) {
	var user []*User
	var resp = newResponse(r)
	if ulr != nil && ulr.Data != nil {
		user = make([]*User, 0, len(ulr.Data))
		for _, d := range ulr.Data {
			a := toUser(d.ID, d.Attributes)
			user = append(user, a)
		}

		firstOffset, err := parseOffsetFromLink(ulr.Links.First)
		if err != nil {
			return user, resp, fmt.Errorf("failed to parse first link: %v", err)
		}
		resp.FirstOffset = firstOffset

		lastOffset, err := parseOffsetFromLink(ulr.Links.Last)
		if err != nil {
			return user, resp, fmt.Errorf("failed to parse last link: %v", err)
		}
		resp.LastOffset = lastOffset

		prevOffset, err := parseOffsetFromLink(ulr.Links.Prev)
		if err != nil {
			return user, resp, fmt.Errorf("failed to parse prev link: %v", err)
		}
		resp.PrevOffset = prevOffset

		nextOffset, err := parseOffsetFromLink(ulr.Links.Next)
		if err != nil {
			return user, resp, fmt.Errorf("failed to parse next link: %v", err)
		}
		resp.NextOffset = nextOffset
	}
	return user, resp, nil
}
