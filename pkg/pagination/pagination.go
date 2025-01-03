package pagination

import (
	"fmt"
	"net/http"
	"strconv"
)

type Request struct {
	Page   int
	Offset int
	Limit  int
}

type Response struct {
	Page           int
	Limit          int
	TotalItems     int
	TotalRetrieved int
}

func FromParam(r *http.Request) (*Request, *Response, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return nil, nil, err
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return nil, nil, err
	}

	if limit < 1 || page < 1 {
		return nil, nil, fmt.Errorf("bad pagination parameter")
	}

	req := &Request{Page: page, Limit: limit, Offset: (page - 1) * limit}
	res := &Response{Page: page, Limit: limit}

	return req, res, nil
}

const (
	PaginationDefaultMaxLimit = 50
)

func MustFromParam(r *http.Request, maxLimit int) (*Request, *Response, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	if maxLimit == 0 {
		maxLimit = PaginationDefaultMaxLimit
	}

	if maxLimit > 0 && limit > maxLimit {
		limit = maxLimit
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	if limit < 1 || page < 1 {
		return nil, nil, fmt.Errorf("bad pagination parameter")
	}

	req := &Request{Page: page, Limit: limit, Offset: (page - 1) * limit}
	res := &Response{Page: page, Limit: limit}

	return req, res, nil
}
