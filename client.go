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

type requestBuilder struct {
	client *http.Client
	prefix string
	method string
	path   string
	query  map[string]string
	header map[string]string
	body   io.Reader
	err    error
}

func NewBuilder() *requestBuilder {
	return &requestBuilder{}
}

func Get(path string) *requestBuilder {
	return NewBuilder().Get(path)
}

func GetJSON(path string, result interface{}) error {
	return NewBuilder().GetJSON(path, result)
}

func GetText(path string) (string, error) {
	return NewBuilder().GetText(path)
}

func GetFile(path string, filename string) error {
	return NewBuilder().GetFile(path, filename)
}

func GetStream(path string) (io.ReadCloser, error) {
	return NewBuilder().GetStream(path)
}

func Post(path string) *requestBuilder {
	return NewBuilder().Post(path)
}

func Put(path string) *requestBuilder {
	return NewBuilder().Put(path)
}

func Delete(path string) *requestBuilder {
	return NewBuilder().Delete(path)
}

func Patch(path string) *requestBuilder {
	return NewBuilder().Patch(path)
}

func Options(path string) *requestBuilder {
	return NewBuilder().Options(path)
}

func (b *requestBuilder) FromTemplate() {

}

func (b *requestBuilder) Client(client *http.Client) *requestBuilder {
	b.client = client
	return b
}

func (b *requestBuilder) Prefix(p string) *requestBuilder {
	b.prefix = p
	return b
}

func (b *requestBuilder) Get(path string) *requestBuilder {
	b.method = http.MethodGet
	b.path = path
	return b
}

func (b *requestBuilder) GetJSON(path string, result interface{}) error {
	return b.Get(path).JSON(result)
}

func (b *requestBuilder) GetText(path string) (string, error) {
	return b.Get(path).Text()
}

func (b *requestBuilder) GetFile(path string, filename string) error {
	return b.Get(path).WriteFile(filename)
}

func (b *requestBuilder) GetStream(path string) (io.ReadCloser, error) {
	return b.Get(path).Stream()
}

func (b *requestBuilder) Post(path string) *requestBuilder {
	b.method = http.MethodPost
	b.path = path
	return b
}

func (b *requestBuilder) PostJSON(path string, body interface{}) *requestBuilder {
	return b.Post(path).WithJSONBody(body)
}

func (b *requestBuilder) Put(path string) *requestBuilder {
	b.method = http.MethodPut
	b.path = path
	return b
}

func (b *requestBuilder) PutJSON(path string, body interface{}) *requestBuilder {
	return b.Put(path).WithJSONBody(body)
}

func (b *requestBuilder) Delete(path string) *requestBuilder {
	b.method = http.MethodDelete
	b.path = path
	return b
}

func (b *requestBuilder) Patch(path string) *requestBuilder {
	b.method = http.MethodPatch
	b.path = path
	return b
}

func (b *requestBuilder) PatchJSON(path string, body interface{}) *requestBuilder {
	return b.Patch(path).WithJSONBody(body)
}

func (b *requestBuilder) Options(path string) *requestBuilder {
	b.method = http.MethodOptions
	b.path = path
	return b
}

func (b *requestBuilder) Query(k, v string) *requestBuilder {
	b.query[k] = v
	return b
}

func (b *requestBuilder) Queries(q map[string]string) *requestBuilder {
	if b.query == nil {
		b.query = map[string]string{}
	}
	for k, v := range q {
		b.query[k] = v
	}
	return b
}

func (b *requestBuilder) Header(k, v string) *requestBuilder {
	b.header[k] = v
	return b
}

func (b *requestBuilder) Headers(h map[string]string) *requestBuilder {
	if b.header == nil {
		b.header = map[string]string{}
	}
	for k, v := range h {
		b.header[k] = v
	}
	return b
}

func (b *requestBuilder) WithTextBody(text string) *requestBuilder {
	b.body = strings.NewReader(text)
	return b
}

func (b *requestBuilder) WithJSONBody(v interface{}) *requestBuilder {
	bs, err := json.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *requestBuilder) WithXMLBody(v interface{}) *requestBuilder {
	bs, err := xml.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *requestBuilder) WithFile(filename string) *requestBuilder {
	file, err := os.Open(filename)
	if err != nil {
		b.err = err
		return b
	}
	b.body = file
	return b
}

func (b *requestBuilder) WithBodyReader(r io.Reader) *requestBuilder {
	b.body = r
	return b
}

func (b *requestBuilder) Do(client ...*http.Client) Response {
	if b.err != nil {
		return ErrorResponse(b.err)
	}

	aURL := b.url()
	if aURL == "" {
		return ErrorResponse(fmt.Errorf("deer: http client: require url"))
	}
	if b.method == "" {
		return ErrorResponse(errors.New("deer: http client: require method"))
	}

	request, err := http.NewRequest(b.method, aURL, b.body)
	if err != nil {
		return ErrorResponse(err)
	}

	for k, v := range b.header {
		request.Header.Set(k, v)
	}

	var finalClient *http.Client
	if len(client) > 0 {
		finalClient = client[0]
	}
	if finalClient == nil {
		finalClient = b.client
	}
	if finalClient == nil {
		finalClient = http.DefaultClient
	}

	response, err := finalClient.Do(request)
	if err != nil {
		return ErrorResponse(err)
	}

	return WrapResponse(response)
}

func (b *requestBuilder) Text() (string, error) {
	return b.Do().Text()
}

func (b *requestBuilder) JSON(value interface{}) error {
	return b.Do().JSON(value)
}

func (b *requestBuilder) XML(value interface{}) error {
	return b.Do().XML(value)
}

func (b *requestBuilder) Stream() (io.ReadCloser, error) {
	return b.Do().Stream()
}

func (b *requestBuilder) Write(writer io.Writer) error {
	return b.Do().Write(writer)
}

func (b *requestBuilder) WriteFile(filename string) error {
	return b.Do().WriteFile(filename)
}

func (b *requestBuilder) url() string {
	s := b.prefix + b.path
	if strings.Contains(s, "?") {
		s += "&"
	} else {
		s += "?"
	}
	if len(b.query) > 0 {
		vs := url.Values{}
		for k, v := range b.query {
			vs.Set(k, v)
		}
		s += vs.Encode()
	}
	return s
}
