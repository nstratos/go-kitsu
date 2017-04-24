package jsonapi

import (
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/nstratos/jsonapi"
)

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func Encode(w io.Writer, v interface{}) (err error) {
	const errFormat = "cannot encode type %T, need pointer to struct or slice of pointers to structs"
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = fmt.Errorf(errFormat+": %v", v, r.(error))
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
		return jsonapi.MarshalOnePayload(w, v)
	case reflect.Slice:
		s := reflect.ValueOf(v)
		if s.Type().Elem().Kind() != reflect.Ptr {
			return fmt.Errorf(errFormat, v)
		}
		if s.Type().Elem().Elem().Kind() != reflect.Struct {
			return fmt.Errorf(errFormat, v)

		}
		return jsonapi.MarshalManyPayload(w, v)
	}
}

func Decode(r io.Reader, ptr interface{}) (Offset, error) {
	const errFormat = "cannot decode to %T, need pointer to struct or pointer to slice"
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return Offset{}, fmt.Errorf(errFormat, ptr)
	}
	v := reflect.Indirect(reflect.ValueOf(ptr))
	switch v.Kind() {
	default:
		return Offset{}, fmt.Errorf(errFormat, ptr)
	case reflect.Struct:
		return Offset{}, jsonapi.UnmarshalPayload(r, ptr)
	case reflect.Slice:
		data, links, err := jsonapi.UnmarshalManyPayloadWithLinks(r, v.Type().Elem())
		if err != nil {
			return Offset{}, err
		}
		for _, d := range data {
			v.Set(reflect.Append(v, reflect.ValueOf(d)))
		}

		o := Offset{}
		var perr error
		if links != nil {
			o, perr = parseOffset(*links)
			if perr != nil {
				return Offset{}, perr
			}
		}
		return o, nil
	}
}

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
