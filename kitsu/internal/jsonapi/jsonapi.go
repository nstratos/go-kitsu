package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/google/jsonapi"
)

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Encode returns the JSON API encoding of v. It requires v to be a pointer to
// struct or a slice of pointers to structs.
func Encode(w io.Writer, v interface{}) (err error) {
	const errFormat = "cannot encode type %T, need pointer to struct or slice of pointers to structs"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = fmt.Errorf("cannot encode type %T: %v", v, r)
		}
	}()
	if isZeroOfUnderlyingType(v) {
		return fmt.Errorf("cannot encode nil value of %#v", v)
	}
	t := reflect.TypeOf(v)
	switch t.Kind() {
	default:
		return fmt.Errorf(errFormat, v)
	case reflect.Ptr:
		if t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf(errFormat, v)
		}
		return jsonapi.MarshalPayload(w, v)
	case reflect.Slice:
		s := reflect.ValueOf(v)
		if s.Type().Elem().Kind() != reflect.Ptr {
			return fmt.Errorf(errFormat, v)
		}
		if s.Type().Elem().Elem().Kind() != reflect.Struct {
			return fmt.Errorf(errFormat, v)

		}
		return jsonapi.MarshalPayload(w, v)
	}
}

type extra struct {
	Links *jsonapi.Links `json:"links,omitempty"`
	Meta  *jsonapi.Meta  `json:"meta,omitempty"` // Not returned for now.
}

// Decode parses the JSON API encoded data and stores the result in the value
// pointed to by v. It requires v to be a pointer to struct or pointer to
// slice.
func Decode(r io.Reader, v interface{}) (offset Offset, err error) {
	const errFormat = "cannot decode to %T, need pointer to struct or pointer to slice"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = fmt.Errorf("cannot decode to %T: %v", v, r)
		}
	}()
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return Offset{}, fmt.Errorf(errFormat, v)
	}
	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	default:
		return Offset{}, fmt.Errorf(errFormat, v)
	case reflect.Struct:
		return Offset{}, jsonapi.UnmarshalPayload(r, v)
	case reflect.Slice:
		var buf bytes.Buffer
		tee := io.TeeReader(r, &buf)

		// Decode data.
		data, uerr := jsonapi.UnmarshalManyPayload(tee, val.Type().Elem())
		if uerr != nil {
			return Offset{}, uerr
		}
		for _, d := range data {
			val.Set(reflect.Append(val, reflect.ValueOf(d)))
		}

		// Decode links.
		x := new(extra)
		if err := json.NewDecoder(&buf).Decode(x); err != nil {
			return Offset{}, err
		}

		o := Offset{}
		var perr error
		if x.Links != nil {
			o, perr = parseOffset(*x.Links)
			if perr != nil {
				return Offset{}, perr
			}
		}
		return o, nil
	}
}
