package jsonapi

import (
	"io"
	"reflect"

	"github.com/nstratos/jsonapi"
)

func EncodeOne(w io.Writer, v interface{}) error {
	return jsonapi.MarshalOnePayload(w, v)
}

func EncodeMany(w io.Writer, v interface{}) error {
	return jsonapi.MarshalManyPayload(w, v)
}

func DecodeOne(r io.Reader, v interface{}) error {
	return jsonapi.UnmarshalPayload(r, v)
}

func DecodeMany(r io.Reader, t reflect.Type) ([]interface{}, Offset, error) {
	v, links, err := jsonapi.UnmarshalManyPayloadWithLinks(r, t)
	if err != nil {
		return nil, Offset{}, err
	}

	o := Offset{}
	var perr error
	if links != nil {
		o, perr = parseOffset(*links)
		if perr != nil {
			return nil, Offset{}, perr
		}
	}
	return v, o, err
}
