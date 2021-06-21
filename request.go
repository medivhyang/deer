package deer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/medivhyang/duck/naming"
	"github.com/medivhyang/duck/reflectutil"
)

const (
	bindingTagKey = "binding"
)

func WrapRequest(r *http.Request) *Request {
	return &Request{Raw: r}
}

type Request struct {
	Raw    *http.Request
	params map[string]string
}

func (r *Request) Context() context.Context {
	return r.Raw.Context()
}

func (r *Request) SetContext(ctx context.Context) {
	r.Raw = r.Raw.WithContext(ctx)
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

func (r *Request) HeaderOrDefault(key string, value string) string {
	result := r.Raw.Header.Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) HeaderExists(key string) bool {
	if r.Raw.Header == nil {
		return false
	}
	return len(r.Raw.Header[key]) > 0
}

func (r *Request) SetHeader(key string, value string) {
	r.Raw.Header.Set(key, value)
}

func (r *Request) Cookie(key string) (string, error) {
	cookie, err := r.Raw.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (r *Request) AddCooke(cookie *http.Cookie) {
	r.Raw.AddCookie(cookie)
}

func (r *Request) CookieOrDefault(key string, defaultValue ...string) string {
	cookie, err := r.Raw.Cookie(key)
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

func (r *Request) CookieExists(key string) bool {
	_, err := r.Raw.Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			return false
		}
		panic(err)
	}
	return true
}

func (r *Request) Param(key string) string {
	if r.params == nil {
		r.params = Params(r.Raw)
	}
	return r.params[key]
}

func (r *Request) ParamOrDefault(key string, value string) string {
	if r.params == nil {
		r.params = Params(r.Raw)
	}
	result := r.params[key]
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ParamExists(key string) bool {
	if r.params == nil {
		r.params = Params(r.Raw)
	}
	_, ok := r.params[key]
	return ok
}

func (r *Request) Query(key string) string {
	return r.Raw.URL.Query().Get(key)
}

func (r *Request) QueryOrDefault(key string, value string) string {
	result := r.Raw.URL.Query().Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) QueryExists(key string) bool {
	values := r.Raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *Request) PostForm(key string) string {
	return r.Raw.PostFormValue(key)
}

func (r *Request) PostFormExists(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.PostForm == nil {
		return false
	}
	return len(r.Raw.PostForm[key]) > 0
}

func (r *Request) Form(key string) string {
	return r.Raw.FormValue(key)
}

func (r *Request) FormExists(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.Form == nil {
		return false
	}
	return len(r.Raw.Form[key]) > 0
}

func (r *Request) BindJSON(i interface{}) error {
	return json.NewDecoder(r.Raw.Body).Decode(i)
}

func (r *Request) BindXML(i interface{}) error {
	return xml.NewDecoder(r.Raw.Body).Decode(i)
}

func (r *Request) BindQuery(i interface{}) error {
	m := reflectutil.ParseStructTag(i, bindingTagKey)
	return reflectutil.BindStructFunc(i, func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = naming.ToSnake(s)
		}
		return r.Raw.URL.Query()[s]
	})
}

func (r *Request) BindPostForm(i interface{}) error {
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	m := reflectutil.ParseStructTag(i, bindingTagKey)
	return reflectutil.BindStructFunc(i, func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = naming.ToSnake(s)
		}
		return r.Raw.PostForm[s]
	})
}

func (r *Request) BindForm(i interface{}) error {
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	m := reflectutil.ParseStructTag(i, bindingTagKey)
	return reflectutil.BindStructFunc(i, func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = naming.ToSnake(s)
		}
		return r.Raw.Form[s]
	})
}

func (r *Request) BasicAuth() (username string, password string, ok bool) {
	return r.Raw.BasicAuth()
}

func (r *Request) SetBasicAuth(username string, password string) {
	r.Raw.SetBasicAuth(username, password)
}
