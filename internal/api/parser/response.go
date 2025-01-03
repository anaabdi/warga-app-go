package parser

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/anaabdi/warga-app-go/pkg/constant"
	apperr "github.com/anaabdi/warga-app-go/pkg/error"
	"github.com/anaabdi/warga-app-go/pkg/pagination"
)

const (
	errorCode = "9999"
)

type Response struct {
	Code    string      `json:"code"`
	Msg     string      `json:"message"`
	Detail  string      `json:"detail,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success" example:"false"`
	Meta    *Meta       `json:"meta,omitempty" swaggerignore:"true"`
}

type HTTPError struct {
	Code    string `json:"code"`
	Msg     string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Success bool   `json:"success" example:"false"`
}

type Meta struct {
	Next      string `json:"next"`
	Current   string `json:"current"`
	Previous  string `json:"previous"`
	TotalItem int    `json:"total_item"`
	TotalPage int    `json:"total_page"`
}

type JSONResponder interface {
	Write(w http.ResponseWriter, status int, data interface{})
	WriteHTML(w http.ResponseWriter, status int, content []byte)
	SuccessNoData(w http.ResponseWriter, status int, location string)
	Error(w http.ResponseWriter, apiErr error)
	Success(w http.ResponseWriter, msg string)
	SuccessWithData(w http.ResponseWriter, msg string, data interface{})
	SuccessWithDataPagination(w http.ResponseWriter, msg string, data interface{}, meta *Meta)
	PreparePaginationMeta(r *http.Request, pgRes *pagination.Response) *Meta
	SuccessWithFile(w http.ResponseWriter, file []byte, contentDisposition, contentType string)
}

type jsonResponder struct {
	contentType string
}

// SuccessWithFile implements JSONResponder.
func (c *jsonResponder) SuccessWithFile(w http.ResponseWriter, file []byte, contentDisposition string, contentType string) {
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(file)
}

var (
	obj  *jsonResponder
	once sync.Once
)

func NewJSONResponder() JSONResponder {
	return &jsonResponder{
		contentType: fmt.Sprintf("%s; charset=%s", jsonContentType, jsonCharset),
	}
}

func SyncedJSONResponder() *jsonResponder {
	once.Do(func() {
		if obj == nil {
			obj = &jsonResponder{
				contentType: fmt.Sprintf("%s; charset=%s", jsonContentType, jsonCharset),
			}
		}
	})
	return obj
}

func (c *jsonResponder) Write(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", c.contentType)
	w.WriteHeader(status)
	if data == nil {
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		c.Error(w, err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	_, _ = w.Write(content)
}

func (c *jsonResponder) WriteHTML(w http.ResponseWriter, status int, content []byte) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	_, _ = w.Write(content)
}

func (c *jsonResponder) Success(w http.ResponseWriter, msg string) {
	content := Response{
		Code:    constant.SuccessCode,
		Msg:     msg,
		Success: true,
	}
	c.Write(w, http.StatusOK, content)
}

func (c *jsonResponder) SuccessWithData(w http.ResponseWriter, msg string, data interface{}) {
	content := Response{
		Code:    constant.SuccessCode,
		Msg:     msg,
		Data:    data,
		Success: true,
	}
	c.Write(w, http.StatusOK, content)
}

func (c *jsonResponder) SuccessWithDataPagination(w http.ResponseWriter, msg string, data interface{}, meta *Meta) {
	content := Response{
		Code:    constant.SuccessCode,
		Msg:     msg,
		Data:    data,
		Meta:    meta,
		Success: true,
	}
	c.Write(w, http.StatusOK, content)
}

func (c *jsonResponder) SuccessNoData(w http.ResponseWriter, status int, location string) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	c.Write(w, status, nil)
}

func (c *jsonResponder) Error(w http.ResponseWriter, err error) {
	apiErr, ok := err.(apperr.APIError)
	if !ok {
		c.Write(w, http.StatusBadRequest, HTTPError{
			Code:    errorCode,
			Msg:     "We failed to process your request",
			Detail:  err.Error(),
			Success: false,
		})
		return
	}

	c.Write(w, apiErr.StatusCode, HTTPError{
		Code:    apiErr.Code,
		Msg:     apiErr.Message,
		Success: false,
	})
}

func (c *jsonResponder) PreparePaginationMeta(r *http.Request, pgRes *pagination.Response) *Meta {
	newURL, err := url.Parse(r.URL.Host)
	if err != nil {
		return nil
	}
	newURL.Path = r.URL.Path

	newURL.RawQuery = r.URL.Query().Encode()

	meta := &Meta{}
	meta.Current = c.defineNewURL(newURL, pgRes.Page, pgRes.Limit)
	meta.TotalItem = pgRes.TotalItems
	total := math.Ceil(float64(pgRes.TotalItems) / float64(pgRes.Limit))
	meta.TotalPage = int(total)

	if pgRes.TotalItems == 0 {
		return meta
	}

	if pgRes.TotalItems == pgRes.Limit {
		// return empty when nothing to show next paging
		meta.Next = c.defineNewURL(newURL, pgRes.Page+1, pgRes.Limit)
	}

	if pgRes.Page > 1 {
		meta.Previous = c.defineNewURL(newURL, pgRes.Page-1, pgRes.Limit)
	}

	return meta
}

func (c *jsonResponder) defineNewURL(n *url.URL, page, perPage int) string {
	v := n.Query()
	v.Set("page", strconv.Itoa(page))
	v.Set("limit", strconv.Itoa(perPage))

	n.RawQuery = v.Encode()

	return n.String()
}
