// +build ignore

// The purpose of this file is to be manually copied under
// vendor/github.com/google/jsonapi whenever the dep tool
// (https://github.com/golang/dep) modifies the vendor folder. For convenience:
//
//     sed '1,2d' unmarshalwithlinks.go > tmp; mv tmp vendor/github.com/google/jsonapi/unmarshalwithlinks.go
//
// This is a temporary solution for the reason described in
// (https://github.com/nstratos/go-kitsu/issues/2).
//
// The build tag on top should be deleted.

package jsonapi

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// UnmarshalManyPayloadWithLinks is a copy of UnmarshalManyPayload with the
// only difference that it also returns a map of the pagination links which are
// included in the JSON API document. The map is used to parse the offset from
// the links and return it to the user in a convenient way.
func UnmarshalManyPayloadWithLinks(in io.Reader, t reflect.Type) ([]interface{}, map[string]string, error) {
	payload := new(ManyPayload)

	if err := json.NewDecoder(in).Decode(payload); err != nil {
		return nil, nil, err
	}

	links := *payload.Links

	if payload.Included != nil {
		includedMap := make(map[string]*Node)
		for _, included := range payload.Included {
			key := fmt.Sprintf("%s,%s", included.Type, included.ID)
			includedMap[key] = included
		}

		var models []interface{}
		for _, data := range payload.Data {
			model := reflect.New(t.Elem())
			err := unmarshalNode(data, model, &includedMap)
			if err != nil {
				return nil, nil, err
			}
			models = append(models, model.Interface())
		}

		return models, links, nil
	}

	var models []interface{}

	for _, data := range payload.Data {
		model := reflect.New(t.Elem())
		err := unmarshalNode(data, model, nil)
		if err != nil {
			return nil, nil, err
		}
		models = append(models, model.Interface())
	}

	return models, links, nil
}
