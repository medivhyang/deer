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
	client  *http.Client
	prefix  string
	method  string
	path    string
	queries map[string]string
	headers map[string]string
	body    io.Reader
	err     error
}

type RequestTemplate struct {
	Client  *http.Client
	Prefix  string
	Queries map[string]string
	Headers map[string]string
}

func (t *RequestTemplate) NewBuilder() *RequestBuilder {
	return NewBuilder().Client(t.Client).Prefix(t.Prefix).Headers(t.Headers).Queries(t.Queries)
}

func NewBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

func Get(path string) *RequestBuilder {
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

func Post(path string) *RequestBuilder {
	return NewBuilder().Post(path)
}

func Put(path string) *RequestBuilder {
	return NewBuilder().Put(path)
}

func Delete(path string) *RequestBuilder {
	return NewBuilder().Delete(path)
}

func Patch(path string) *RequestBuilder {
	return NewBuilder().Patch(path)
}

func Options(path string) *RequestBuilder {
	return NewBuilder().Options(path)
}

func (b *RequestBuilder) Client(client *http.Client) *RequestBuilder {
	b.client = client
	return b
}

func (b *RequestBuilder) Prefix(p string) *RequestBuilder {
	b.prefix = p
	return b
}

func (b *RequestBuilder) Get(path string) *RequestBuilder {
	b.method = http.MethodGet
	b.path = path
	return b
}

func (b *RequestBuilder) GetJSON(path string, result interface{}) error {
	return b.Get(path).JSON(result)
}

func (b *RequestBuilder) GetText(path string) (string, error) {
	return b.Get(path).Text()
}

func (b *RequestBuilder) GetFile(path string, filename string) error {
	return b.Get(path).WriteFile(filename)
}

func (b *RequestBuilder) GetStream(path string) (io.ReadCloser, error) {
	return b.Get(path).Stream()
}

func (b *RequestBuilder) Post(path string) *RequestBuilder {
	b.method = http.MethodPost
	b.path = path
	return b
}

func (b *RequestBuilder) PostJSON(path string, body interface{}) *RequestBuilder {
	return b.Post(path).WithJSONBody(body)
}

func (b *RequestBuilder) Put(path string) *RequestBuilder {
	b.method = http.MethodPut
	b.path = path
	return b
}

func (b *RequestBuilder) PutJSON(path string, body interface{}) *RequestBuilder {
	return b.Put(path).WithJSONBody(body)
}

func (b *RequestBuilder) Delete(path string) *RequestBuilder {
	b.method = http.MethodDelete
	b.path = path
	return b
}

func (b *RequestBuilder) Patch(path string) *RequestBuilder {
	b.method = http.MethodPatch
	b.path = path
	return b
}

func (b *RequestBuilder) PatchJSON(path string, body interface{}) *RequestBuilder {
	return b.Patch(path).WithJSONBody(body)
}

func (b *RequestBuilder) Options(path string) *RequestBuilder {
	b.method = http.MethodOptions
	b.path = path
	return b
}

func (b *RequestBuilder) Query(k, v string) *RequestBuilder {
	b.queries[k] = v
	return b
}

func (b *RequestBuilder) Queries(m map[string]string) *RequestBuilder {
	if b.queries == nil {
		b.queries = map[string]string{}
	}
	for k, v := range m {
		b.queries[k] = v
	}
	return b
}

func (b *RequestBuilder) Header(k, v string) *RequestBuilder {
	b.headers[k] = v
	return b
}

func (b *RequestBuilder) Headers(m map[string]string) *RequestBuilder {
	if b.headers == nil {
		b.headers = map[string]string{}
	}
	for k, v := range m {
		b.headers[k] = v
	}
	return b
}

func (b *RequestBuilder) WithTextBody(text string) *RequestBuilder {
	b.body = strings.NewReader(text)
	return b
}

func (b *RequestBuilder) WithJSONBody(v interface{}) *RequestBuilder {
	bs, err := json.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *RequestBuilder) WithXMLBody(v interface{}) *RequestBuilder {
	bs, err := xml.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *RequestBuilder) WithFile(filename string) *RequestBuilder {
	file, err := os.Open(filename)
	if err != nil {
		b.err = err
		return b
	}
	b.body = file
	return b
}

func (b *RequestBuilder) WithReaderBody(r io.Reader) *RequestBuilder {
	b.body = r
	return b
}

func (b *RequestBuilder) Do(client ...*http.Client) Response {
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

	for k, v := range b.headers {
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

func (b *RequestBuilder) Text() (string, error) {
	return b.Do().Text()
}

func (b *RequestBuilder) JSON(value interface{}) error {
	return b.Do().JSON(value)
}

func (b *RequestBuilder) XML(value interface{}) error {
	return b.Do().XML(value)
}

func (b *RequestBuilder) Stream() (io.ReadCloser, error) {
	return b.Do().Stream()
}

func (b *RequestBuilder) Write(writer io.Writer) error {
	return b.Do().Write(writer)
}

func (b *RequestBuilder) WriteFile(filename string) error {
	return b.Do().WriteFile(filename)
}

func (b *RequestBuilder) url() string {
	s := b.prefix + b.path
	if strings.Contains(s, "?") {
		s += "&"
	} else {
		s += "?"
	}
	if len(b.queries) > 0 {
		vs := url.Values{}
		for k, v := range b.queries {
			vs.Set(k, v)
		}
		s += vs.Encode()
	}
	return s
}
