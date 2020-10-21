package deer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/medivhyang/deer/internal/binding"
	"net/http"
)

func WrapRequest(r *http.Request) *Request {
	return &Request{raw: r}
}

type Request struct {
	pathParams map[string]string
	raw        *http.Request
}

func (r *Request) Raw() *http.Request {
	return r.raw
}

func (r *Request) Context() context.Context {
	return r.raw.Context()
}

func (r *Request) Method() string {
	return r.raw.Method
}

func (r *Request) Path() string {
	return r.raw.URL.Path
}

func (r *Request) Header(key string) string {
	return r.raw.Header.Get(key)
}

func (r *Request) HeaderOrDefault(key string, value string) string {
	result := r.raw.Header.Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ExistsHeader(key string) bool {
	if r.raw.Header == nil {
		return false
	}
	return len(r.raw.Header[key]) > 0
}

func (r *Request) SetHeader(key string, value string) {
	r.raw.Header.Set(key, value)
}

func (r *Request) Cookie(key string) (string, error) {
	cookie, err := r.raw.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (r *Request) SetCooke(cookie *http.Cookie) {
	r.raw.AddCookie(cookie)
}

func (r *Request) CookieOrDefault(key string, defaultValue ...string) string {
	cookie, err := r.raw.Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return ""
		}
		panic(err)
	}
	return cookie.Value
}

func (r *Request) ExistsCookie(key string) bool {
	_, err := r.raw.Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			return false
		}
		panic(err)
	}
	return true
}

func (r *Request) PathParam(key string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.raw)
	}
	return r.pathParams[key]
}

func (r *Request) PathParamOrDefault(key string, value string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.raw)
	}
	result := r.pathParams[key]
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ExistsPathParam(key string) bool {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.raw)
	}
	_, ok := r.pathParams[key]
	return ok
}

func (r *Request) Query(key string) string {
	return r.raw.URL.Query().Get(key)
}

func (r *Request) QueryOrDefault(key string, value string) string {
	result := r.raw.URL.Query().Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ExistsQuery(key string) bool {
	values := r.raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *Request) PostForm(key string) string {
	return r.raw.PostFormValue(key)
}

func (r *Request) ExistsPostForm(key string) bool {
	_ = r.raw.ParseForm()
	if r.raw.PostForm == nil {
		return false
	}
	return len(r.raw.PostForm[key]) > 0
}

func (r *Request) Form() map[string][]string {
	return r.raw.Form
}

func (r *Request) FormValue(key string) string {
	return r.raw.FormValue(key)
}

func (r *Request) ExistsForm(key string) bool {
	_ = r.raw.ParseForm()
	if r.raw.Form == nil {
		return false
	}
	return len(r.raw.Form[key]) > 0
}

func (r *Request) DecodeJSONBody(value interface{}) error {
	return json.NewDecoder(r.raw.Body).Decode(value)
}

func (r *Request) DecodeXMLBody(value interface{}) error {
	return xml.NewDecoder(r.raw.Body).Decode(value)
}

func (r *Request) BindWithQuery(target interface{}) error {
	return binding.BindWithQuery(target, r.raw.URL.Query())
}

func (r *Request) BindWithPostForm(target interface{}) error {
	if err := r.raw.ParseForm(); err != nil {
		return err
	}
	return binding.BindWithPostForm(target, r.raw.PostForm)
}

func (r *Request) BindWithForm(target interface{}) error {
	if err := r.raw.ParseForm(); err != nil {
		return err
	}
	return binding.BindWithForm(target, r.raw.Form)
}

func (r *Request) BasicAuth() (username string, password string, ok bool) {
	return r.raw.BasicAuth()
}

func (r *Request) SetBasicAuth(username string, password string) {
	r.raw.SetBasicAuth(username, password)
}
