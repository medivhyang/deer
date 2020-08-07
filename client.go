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

type ClientBuilder struct {
	prefix string
	method string
	path   string
	query  map[string]string
	header map[string]string
	body   io.Reader
	err    error
}

func Get(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodGet,
	}
}

func Post(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodPost,
	}
}

func Put(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodPut,
	}
}

func Delete(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodDelete,
	}
}

func Patch(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodPatch,
	}
}

func Options(path string) *ClientBuilder {
	return &ClientBuilder{
		path:   path,
		method: http.MethodOptions,
	}
}

func NewBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (c *ClientBuilder) Prefix(p string) *ClientBuilder {
	c.prefix = p
	return c
}

func (c *ClientBuilder) Get(path string) *ClientBuilder {
	c.method = http.MethodGet
	c.path = path
	return c
}

func (c *ClientBuilder) Post(path string) *ClientBuilder {
	c.method = http.MethodPost
	c.path = path
	return c
}

func (c *ClientBuilder) Put(path string) *ClientBuilder {
	c.method = http.MethodPut
	c.path = path
	return c
}

func (c *ClientBuilder) Delete(path string) *ClientBuilder {
	c.method = http.MethodDelete
	c.path = path
	return c
}

func (c *ClientBuilder) Patch(path string) *ClientBuilder {
	c.method = http.MethodPatch
	c.path = path
	return c
}

func (c *ClientBuilder) Options(path string) *ClientBuilder {
	c.method = http.MethodOptions
	c.path = path
	return c
}

func (c *ClientBuilder) Query(q map[string]string) *ClientBuilder {
	if c.query == nil {
		c.query = map[string]string{}
	}
	for k, v := range q {
		c.query[k] = v
	}
	return c
}

func (c *ClientBuilder) Header(h map[string]string) *ClientBuilder {
	if c.header == nil {
		c.header = map[string]string{}
	}
	for k, v := range h {
		c.header[k] = v
	}
	return c
}

func (c *ClientBuilder) WithTextBody(text string) *ClientBuilder {
	c.body = strings.NewReader(text)
	return c
}

func (c *ClientBuilder) WithJSONBody(v interface{}) *ClientBuilder {
	bs, err := json.Marshal(v)
	if err != nil {
		c.err = err
		return c
	}
	c.body = bytes.NewReader(bs)
	return c
}

func (c *ClientBuilder) WithXMLBody(v interface{}) *ClientBuilder {
	bs, err := xml.Marshal(v)
	if err != nil {
		c.err = err
		return c
	}
	c.body = bytes.NewReader(bs)
	return c
}

func (c *ClientBuilder) WithFile(filename string) *ClientBuilder {
	file, err := os.Open(filename)
	if err != nil {
		c.err = err
		return c
	}
	c.body = file
	return c
}

func (c *ClientBuilder) WithBodyReader(r io.Reader) *ClientBuilder {
	c.body = r
	return c
}

func (c *ClientBuilder) Do(client ...*http.Client) (*ResponseAdapter, error) {
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

func (c *ClientBuilder) url() string {
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
