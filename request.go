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

func (request *RequestAdapter) Method() string {
	return request.Raw.Method
}

func (request *RequestAdapter) Path() string {
	return request.Raw.URL.Path
}

func (request *RequestAdapter) Header(key string) string {
	return request.Raw.Header.Get(key)
}

func (request *RequestAdapter) SetHeader(key string, value string) {
	request.Raw.Header.Set(key, value)
}

func (request *RequestAdapter) PathParams() map[string]string {
	if request.pathParams == nil {
		request.pathParams = PathParams(request.Raw)
	}
	return request.pathParams
}

func (request *RequestAdapter) PathParam(key string) string {
	if request.pathParams == nil {
		request.pathParams = PathParams(request.Raw)
	}
	return request.pathParams[key]
}

func (request *RequestAdapter) QueryExists(key string) bool {
	values := request.Raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (request *RequestAdapter) Query(key string) string {
	return request.Raw.URL.Query().Get(key)
}

func (request *RequestAdapter) QuerySlice(key string, sep ...string) []string {
	if !request.QueryExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(request.Query(key), aSep)
}

func (request *RequestAdapter) QuerySliceTrim(key string, sep ...string) []string {
	items := request.QuerySlice(key, sep...)
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

func (request *RequestAdapter) QueryInt(key string) (int, error) {
	return strconv.Atoi(request.Query(key))
}

func (request *RequestAdapter) QueryIntSlice(key string, sep ...string) ([]int, error) {
	items := request.QuerySliceTrim(key, sep...)
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

func (request *RequestAdapter) QueryInt64(key string) (int64, error) {
	return strconv.ParseInt(request.Query(key), 10, 64)
}

func (request *RequestAdapter) QueryInt64Slice(key string, sep ...string) ([]int64, error) {
	items := request.QuerySliceTrim(key, sep...)
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

func (request *RequestAdapter) QueryFloat64(key string) (float64, error) {
	return strconv.ParseFloat(request.Query(key), 64)
}

func (request *RequestAdapter) QueryBool(key string) (bool, error) {
	return strconv.ParseBool(request.Query(key))
}

func (request *RequestAdapter) QueryTime(key string, layout ...string) (t time.Time, err error) {
	if len(layout) > 0 {
		return time.Parse(layout[0], request.Query(key))
	}
	for _, format := range defaultTimeFormats {
		t, err = time.Parse(format, request.Query(key))
		if err == nil {
			return
		}
	}
	return
}

func (request *RequestAdapter) QueryTimeUnix(key string) (time.Time, error) {
	sec, err := request.QueryInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (request *RequestAdapter) FormExists(key string) bool {
	_ = request.Raw.ParseForm()
	if request.Raw.PostForm == nil {
		return false
	}
	return len(request.Raw.PostForm[key]) > 0
}

func (request *RequestAdapter) Form(key string) string {
	return request.Raw.PostFormValue(key)
}

func (request *RequestAdapter) FormSlice(key string, sep ...string) []string {
	if !request.FormExists(key) {
		return nil
	}
	aSep := defaultSliceSeparator
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Split(request.Form(key), aSep)
}

func (request *RequestAdapter) FormSliceTrim(key string, sep ...string) []string {
	items := request.FormSlice(key, sep...)
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

func (request *RequestAdapter) FormInt(key string) (int, error) {
	return strconv.Atoi(request.Form(key))
}

func (request *RequestAdapter) FormIntSlice(key string, sep ...string) ([]int, error) {
	items := request.FormSliceTrim(key, sep...)
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

func (request *RequestAdapter) FormInt64(key string) (int64, error) {
	return strconv.ParseInt(request.Form(key), 10, 64)
}

func (request *RequestAdapter) FormInt64Slice(key string, sep ...string) ([]int64, error) {
	items := request.FormSliceTrim(key, sep...)
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

func (request *RequestAdapter) FormFloat64(key string) (float64, error) {
	return strconv.ParseFloat(request.Form(key), 64)
}

func (request *RequestAdapter) FormBool(key string) (bool, error) {
	return strconv.ParseBool(request.Form(key))
}

func (request *RequestAdapter) FormTime(key string, layout ...string) (t time.Time, err error) {
	if len(layout) > 0 {
		return time.Parse(layout[0], request.Form(key))
	}
	for _, format := range defaultTimeFormats {
		t, err = time.Parse(format, request.Form(key))
		if err == nil {
			return
		}
	}
	return
}

func (request *RequestAdapter) FormTimeUnix(key string) (time.Time, error) {
	sec, err := request.FormInt64(key)
	if err != nil {
		return time.Time{}, nil
	}
	return time.Unix(sec, 0), nil
}

func (request *RequestAdapter) BindWithJSON(value interface{}) error {
	return json.NewDecoder(request.Raw.Body).Decode(value)
}

func (request *RequestAdapter) BindWithXML(value interface{}) error {
	return xml.NewDecoder(request.Raw.Body).Decode(value)
}
