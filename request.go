package deer

import (
	"encoding/json"
	"encoding/xml"
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

func WrapRequest(r *http.Request) *RequestAdapter {
	return &RequestAdapter{Raw: r}
}

type RequestAdapter struct {
	pathParams map[string]string
	Raw        *http.Request
}

func (r *RequestAdapter) Method() string {
	return r.Raw.Method
}

func (r *RequestAdapter) Path() string {
	return r.Raw.URL.Path
}

func (r *RequestAdapter) Header(key string) string {
	return r.Raw.Header.Get(key)
}

func (r *RequestAdapter) SetHeader(key string, value string) {
	r.Raw.Header.Set(key, value)
}

func (r *RequestAdapter) PathParams() map[string]string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	return r.pathParams
}

func (r *RequestAdapter) PathParam(key string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	return r.pathParams[key]
}

func (r *RequestAdapter) PathParamInt(key string) (int, error) {
	return strconv.Atoi(r.PathParam(key))
}

func (r *RequestAdapter) PathParamInt64(key string) (int64, error) {
	return strconv.ParseInt(r.PathParam(key), 10, 64)
}

func (r *RequestAdapter) PathParamFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.PathParam(key),  64)
}

func (r *RequestAdapter) PathParamBool(key string) (bool, error) {
	return strconv.ParseBool(r.PathParam(key))
}

func (r *RequestAdapter) PathParamTime(key string, layout ...string) (t time.Time, err error) {
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

func (r *RequestAdapter) PathParamTimeUnix(key string) (t time.Time, err error) {
	sec, err := r.PathParamInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *RequestAdapter) QueryExists(key string) bool {
	values := r.Raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *RequestAdapter) Query(key string) string {
	return r.Raw.URL.Query().Get(key)
}

func (r *RequestAdapter) QuerySlice(key string, sep ...string) []string {
	if !r.QueryExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(r.Query(key), aSep)
}

func (r *RequestAdapter) QuerySliceTrim(key string, sep ...string) []string {
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

func (r *RequestAdapter) QueryInt(key string) (int, error) {
	return strconv.Atoi(r.Query(key))
}

func (r *RequestAdapter) QueryIntSlice(key string, sep ...string) ([]int, error) {
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

func (r *RequestAdapter) QueryInt64(key string) (int64, error) {
	return strconv.ParseInt(r.Query(key), 10, 64)
}

func (r *RequestAdapter) QueryInt64Slice(key string, sep ...string) ([]int64, error) {
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

func (r *RequestAdapter) QueryFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.Query(key), 64)
}

func (r *RequestAdapter) QueryBool(key string) (bool, error) {
	return strconv.ParseBool(r.Query(key))
}

func (r *RequestAdapter) QueryTime(key string, layout ...string) (t time.Time, err error) {
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

func (r *RequestAdapter) QueryTimeUnix(key string) (time.Time, error) {
	sec, err := r.QueryInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *RequestAdapter) FormExists(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.PostForm == nil {
		return false
	}
	return len(r.Raw.PostForm[key]) > 0
}

func (r *RequestAdapter) Form(key string) string {
	return r.Raw.PostFormValue(key)
}

func (r *RequestAdapter) FormSlice(key string, sep ...string) []string {
	if !r.FormExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(r.Form(key), aSep)
}

func (r *RequestAdapter) FormSliceTrim(key string, sep ...string) []string {
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

func (r *RequestAdapter) FormInt(key string) (int, error) {
	return strconv.Atoi(r.Form(key))
}

func (r *RequestAdapter) FormIntSlice(key string, sep ...string) ([]int, error) {
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

func (r *RequestAdapter) FormInt64(key string) (int64, error) {
	return strconv.ParseInt(r.Form(key), 10, 64)
}

func (r *RequestAdapter) FormInt64Slice(key string, sep ...string) ([]int64, error) {
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

func (r *RequestAdapter) FormFloat64(key string) (float64, error) {
	return strconv.ParseFloat(r.Form(key), 64)
}

func (r *RequestAdapter) FormBool(key string) (bool, error) {
	return strconv.ParseBool(r.Form(key))
}

func (r *RequestAdapter) FormTime(key string, layout ...string) (t time.Time, err error) {
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

func (r *RequestAdapter) FormTimeUnix(key string) (time.Time, error) {
	sec, err := r.FormInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (r *RequestAdapter) BindWithJSON(value interface{}) error {
	return json.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *RequestAdapter) BindWithXML(value interface{}) error {
	return xml.NewDecoder(r.Raw.Body).Decode(value)
}
