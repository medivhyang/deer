package deer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/medivhyang/deer/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultSliceSeparator = ","
	defaultTimeFormats    = []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"20060102",
	}
)

func WrapRequest(r *http.Request) *Request {
	return &Request{Raw: r}
}

type Request struct {
	pathParams map[string]string
	Raw        *http.Request
}

func (r *Request) Context() context.Context {
	return r.Raw.Context()
}

func (r *Request) Method() string {
	return r.Raw.Method
}

func (r *Request) Path() string {
	return r.Raw.URL.Path
}

func (r *Request) Header(key string) string {
	return r.Raw.Header.Get(key)
}

func (r *Request) SetHeader(key string, value string) {
	r.Raw.Header.Set(key, value)
}

func (r *Request) PathParams() map[string]string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	return r.pathParams
}

func (r *Request) PathParam(key string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	return r.pathParams[key]
}

func (r *Request) PathParamInt(key string) (int, error) {
	return strconv.Atoi(r.PathParam(key))
}

func (r *Request) PathParamInt64(key string) (int64, error) {
	return strconv.ParseInt(r.PathParam(key), 10, 64)
}

func (r *Request) PathParamFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.PathParam(key), 64)
}

func (r *Request) PathParamBool(key string) (bool, error) {
	return strconv.ParseBool(r.PathParam(key))
}

func (r *Request) PathParamTime(key string, layout ...string) (t time.Time, err error) {
	if len(layout) > 0 {
		return time.Parse(layout[0], r.PathParam(key))
	}
	for _, format := range defaultTimeFormats {
		t, err = time.Parse(format, r.PathParam(key))
		if err == nil {
			return
		}
	}
	return
}

func (r *Request) PathParamTimeUnix(key string) (t time.Time, err error) {
	sec, err := r.PathParamInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *Request) QueryExists(key string) bool {
	values := r.Raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *Request) Query(key string) string {
	return r.Raw.URL.Query().Get(key)
}

func (r *Request) QuerySlice(key string, sep ...string) []string {
	if !r.QueryExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(r.Query(key), aSep)
}

func (r *Request) QuerySliceTrim(key string, sep ...string) []string {
	items := r.QuerySlice(key, sep...)
	if len(items) == 0 {
		return nil
	}
	var result []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (r *Request) QueryInt(key string) (int, error) {
	return strconv.Atoi(r.Query(key))
}

func (r *Request) QueryIntSlice(key string, sep ...string) ([]int, error) {
	items := r.QuerySliceTrim(key, sep...)
	if len(items) == 0 {
		return nil, nil
	}
	var result []int
	for _, item := range items {
		value, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (r *Request) QueryInt64(key string) (int64, error) {
	return strconv.ParseInt(r.Query(key), 10, 64)
}

func (r *Request) QueryInt64Slice(key string, sep ...string) ([]int64, error) {
	items := r.QuerySliceTrim(key, sep...)
	if len(items) == 0 {
		return nil, nil
	}
	var result []int64
	for _, item := range items {
		value, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (r *Request) QueryFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.Query(key), 64)
}

func (r *Request) QueryBool(key string) (bool, error) {
	return strconv.ParseBool(r.Query(key))
}

func (r *Request) QueryTime(key string, layout ...string) (t time.Time, err error) {
	if len(layout) > 0 {
		return time.Parse(layout[0], r.Query(key))
	}
	for _, format := range defaultTimeFormats {
		t, err = time.Parse(format, r.Query(key))
		if err == nil {
			return
		}
	}
	return
}

func (r *Request) QueryTimeUnix(key string) (time.Time, error) {
	sec, err := r.QueryInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *Request) FormExists(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.PostForm == nil {
		return false
	}
	return len(r.Raw.PostForm[key]) > 0
}

func (r *Request) Form(key string) string {
	return r.Raw.PostFormValue(key)
}

func (r *Request) FormSlice(key string, sep ...string) []string {
	if !r.FormExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(r.Form(key), aSep)
}

func (r *Request) FormSliceTrim(key string, sep ...string) []string {
	items := r.FormSlice(key, sep...)
	if len(items) == 0 {
		return nil
	}
	var result []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (r *Request) FormInt(key string) (int, error) {
	return strconv.Atoi(r.Form(key))
}

func (r *Request) FormIntSlice(key string, sep ...string) ([]int, error) {
	items := r.FormSliceTrim(key, sep...)
	if len(items) == 0 {
		return nil, nil
	}
	var result []int
	for _, item := range items {
		value, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (r *Request) FormInt64(key string) (int64, error) {
	return strconv.ParseInt(r.Form(key), 10, 64)
}

func (r *Request) FormInt64Slice(key string, sep ...string) ([]int64, error) {
	items := r.FormSliceTrim(key, sep...)
	if len(items) == 0 {
		return nil, nil
	}
	var result []int64
	for _, item := range items {
		value, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (r *Request) FormFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.Form(key), 64)
}

func (r *Request) FormBool(key string) (bool, error) {
	return strconv.ParseBool(r.Form(key))
}

func (r *Request) FormTime(key string, layout ...string) (t time.Time, err error) {
	if len(layout) > 0 {
		return time.Parse(layout[0], r.Form(key))
	}
	for _, format := range defaultTimeFormats {
		t, err = time.Parse(format, r.Form(key))
		if err == nil {
			return
		}
	}
	return
}

func (r *Request) FormTimeUnix(key string) (time.Time, error) {
	sec, err := r.FormInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *Request) BindWithJSON(value interface{}) error {
	return json.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *Request) BindWithXML(value interface{}) error {
	return xml.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *Request) BindWithQuery(target interface{}) error {
	return binding.BindWithQuery(target, r.Raw.URL.Query())
}

func (r *Request) BindWithPostForm(target interface{}) error {
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	return binding.BindWithPostForm(target, r.Raw.PostForm)
}

func (r *Request) BindWithForm(target interface{}) error {
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	return binding.BindWithForm(target, r.Raw.Form)
}
