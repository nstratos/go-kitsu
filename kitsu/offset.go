package kitsu

import (
	"fmt"
	"net/url"
	"strconv"
)

type offset struct {
	first int
	last  int
	prev  int
	next  int
}

func parseOffset(links map[string]string) (*offset, error) {
	var first, last, prev, next int
	var err error

	firstStr, ok := links["first"]
	if ok {
		first, err = parseOffsetFromLink(firstStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse first link: %v", err)
		}
	}

	lastStr, ok := links["last"]
	if ok {
		last, err = parseOffsetFromLink(lastStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last link: %v", err)
		}
	}

	prevStr, ok := links["prev"]
	if ok {
		prev, err = parseOffsetFromLink(prevStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse prev link: %v", err)
		}
	}

	nextStr, ok := links["next"]
	if ok {
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
