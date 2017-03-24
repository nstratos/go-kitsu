package kitsu

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nstratos/jsonapi"
)

type offset struct {
	first int
	last  int
	prev  int
	next  int
}

func parseOffset(links jsonapi.Links) (*offset, error) {
	var first, last, prev, next int
	var err error

	firstVal, ok := links["first"]
	if ok {
		firstStr, isString := firstVal.(string)
		if !isString {
			return nil, fmt.Errorf("first link is not a string")
		}
		first, err = parseOffsetFromLink(firstStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse first link: %v", err)
		}
	}

	lastVal, ok := links["last"]
	if ok {
		lastStr, isString := lastVal.(string)
		if !isString {
			return nil, fmt.Errorf("last link is not a string")
		}
		last, err = parseOffsetFromLink(lastStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last link: %v", err)
		}
	}

	prevVal, ok := links["prev"]
	if ok {
		prevStr, isString := prevVal.(string)
		if !isString {
			return nil, fmt.Errorf("prev link is not a string")
		}
		prev, err = parseOffsetFromLink(prevStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse prev link: %v", err)
		}
	}

	nextVal, ok := links["next"]
	if ok {
		nextStr, isString := nextVal.(string)
		if !isString {
			return nil, fmt.Errorf("next link is not a string")
		}
		next, err = parseOffsetFromLink(nextStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse next link: %v", err)
		}
	}

	o := &offset{
		first: first,
		last:  last,
		prev:  prev,
		next:  next,
	}
	return o, nil
}

func parseOffsetFromLink(link string) (int, error) {
	var offset int
	u, err := url.Parse(link)
	if err != nil {
		return offset, err
	}
	v := u.Query()
	s := v.Get("page[offset]")
	if s == "" {
		return offset, nil
	}
	return strconv.Atoi(s)
}
