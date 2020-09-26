package deer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/medivhyang/deer/internal/binding"
	"net/http"
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

func (r *Request) HeaderOrDefault(key string, value string) string {
	result := r.Raw.Header.Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ExistsHeader(key string) bool {
	if r.Raw.Header == nil {
		return false
	}
	return len(r.Raw.Header[key]) > 0
}

func (r *Request) SetHeader(key string, value string) {
	r.Raw.Header.Set(key, value)
}

func (r *Request) PathParam(key string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	return r.pathParams[key]
}

func (r *Request) PathParamOrDefault(key string, value string) string {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	result := r.pathParams[key]
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ExistsPathParam(key string) bool {
	if r.pathParams == nil {
		r.pathParams = PathParams(r.Raw)
	}
	_, ok := r.pathParams[key]
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

func (r *Request) ExistsQuery(key string) bool {
	values := r.Raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *Request) PostForm(key string) string {
	return r.Raw.PostFormValue(key)
}

func (r *Request) ExistsPostForm(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.PostForm == nil {
		return false
	}
	return len(r.Raw.PostForm[key]) > 0
}

func (r *Request) Form() map[string][]string {
	return r.Raw.Form
}

func (r *Request) FormValue(key string) string {
	return r.Raw.FormValue(key)
}

func (r *Request) FormExists(key string) bool {
	_ = r.Raw.ParseForm()
	if r.Raw.Form == nil {
		return false
	}
	return len(r.Raw.Form[key]) > 0
}

func (r *Request) DecodeJSONBody(value interface{}) error {
	return json.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *Request) DecodeXMLBody(value interface{}) error {
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
