package deer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RequestBuilder struct {
	prefix string
	method string
	path   string
	query  map[string]string
	header map[string]string
	body   io.Reader
	err    error
}

func Get(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodGet,
	}
}

func Post(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodPost,
	}
}

func Put(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodPut,
	}
}

func Delete(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodDelete,
	}
}

func Patch(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodPatch,
	}
}

func Options(path string) *RequestBuilder {
	return &RequestBuilder{
		path:   path,
		method: http.MethodOptions,
	}
}

func NewBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

func (c *RequestBuilder) Prefix(p string) *RequestBuilder {
	c.prefix = p
	return c
}

func (c *RequestBuilder) Get(path string) *RequestBuilder {
	c.method = http.MethodGet
	c.path = path
	return c
}

func (c *RequestBuilder) Post(path string) *RequestBuilder {
	c.method = http.MethodPost
	c.path = path
	return c
}

func (c *RequestBuilder) Put(path string) *RequestBuilder {
	c.method = http.MethodPut
	c.path = path
	return c
}

func (c *RequestBuilder) Delete(path string) *RequestBuilder {
	c.method = http.MethodDelete
	c.path = path
	return c
}

func (c *RequestBuilder) Patch(path string) *RequestBuilder {
	c.method = http.MethodPatch
	c.path = path
	return c
}

func (c *RequestBuilder) Options(path string) *RequestBuilder {
	c.method = http.MethodOptions
	c.path = path
	return c
}

func (c *RequestBuilder) Query(q map[string]string) *RequestBuilder {
	if c.query == nil {
		c.query = map[string]string{}
	}
	for k, v := range q {
		c.query[k] = v
	}
	return c
}

func (c *RequestBuilder) Header(h map[string]string) *RequestBuilder {
	if c.header == nil {
		c.header = map[string]string{}
	}
	for k, v := range h {
		c.header[k] = v
	}
	return c
}

func (c *RequestBuilder) WithTextBody(text string) *RequestBuilder {
	c.body = strings.NewReader(text)
	return c
}

func (c *RequestBuilder) WithJSONBody(v interface{}) *RequestBuilder {
	bs, err := json.Marshal(v)
	if err != nil {
		c.err = err
		return c
	}
	c.body = bytes.NewReader(bs)
	return c
}

func (c *RequestBuilder) WithXMLBody(v interface{}) *RequestBuilder {
	bs, err := xml.Marshal(v)
	if err != nil {
		c.err = err
		return c
	}
	c.body = bytes.NewReader(bs)
	return c
}

func (c *RequestBuilder) WithFile(filename string) *RequestBuilder {
	file, err := os.Open(filename)
	if err != nil {
		c.err = err
		return c
	}
	c.body = file
	return c
}

func (c *RequestBuilder) WithBodyReader(r io.Reader) *RequestBuilder {
	c.body = r
	return c
}

func (c *RequestBuilder) Do(client ...*http.Client) (*Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	aURL := c.url()
	if aURL == "" {
		return nil, fmt.Errorf("deer: http client: require url")
	}
	if c.method == "" {
		return nil, errors.New("deer: http client: require method")
	}

	request, err := http.NewRequest(c.method, aURL, c.body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.header {
		request.Header.Set(k, v)
	}

	var aClient *http.Client
	if len(client) > 0 {
		aClient = client[0]
	} else {
		aClient = http.DefaultClient
	}

	response, err := aClient.Do(request)
	if err != nil {
		return nil, err
	}

	return WrapResponse(response), nil
}

func (c *RequestBuilder) url() string {
	s := c.prefix + c.path
	if strings.Contains(s, "?") {
		s += "&"
	} else {
		s += "?"
	}
	if len(c.query) > 0 {
		vs := url.Values{}
		for k, v := range c.query {
			vs.Set(k, v)
		}
		s += vs.Encode()
	}
	return s
}
