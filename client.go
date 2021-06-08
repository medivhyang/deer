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

func (t *RequestTemplate) New() *RequestBuilder {
	return NewBuilder().Client(t.Client).Prefix(t.Prefix).Headers(t.Headers).Queries(t.Queries)
}

func NewBuilder() *RequestBuilder {
	return &RequestBuilder{}
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

func (b *RequestBuilder) SaveFile(path string, filename string) error {
	return b.Get(path).File(filename)
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
	return b.Post(path).WriteJSON(body)
}

func (b *RequestBuilder) Put(path string) *RequestBuilder {
	b.method = http.MethodPut
	b.path = path
	return b
}

func (b *RequestBuilder) PutJSON(path string, body interface{}) *RequestBuilder {
	return b.Put(path).WriteJSON(body)
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
	return b.Patch(path).WriteJSON(body)
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

func (b *RequestBuilder) WriteText(text string) *RequestBuilder {
	b.body = strings.NewReader(text)
	return b
}

func (b *RequestBuilder) WriteJSON(v interface{}) *RequestBuilder {
	bs, err := json.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *RequestBuilder) WriteXML(v interface{}) *RequestBuilder {
	bs, err := xml.Marshal(v)
	if err != nil {
		b.err = err
		return b
	}
	b.body = bytes.NewReader(bs)
	return b
}

func (b *RequestBuilder) WriteFile(filename string) *RequestBuilder {
	file, err := os.Open(filename)
	if err != nil {
		b.err = err
		return b
	}
	b.body = file
	return b
}

func (b *RequestBuilder) ReaderBody(r io.Reader) *RequestBuilder {
	b.body = r
	return b
}

func (b *RequestBuilder) Do(client ...*http.Client) Response {
	if b.err != nil {
		return ErrorResponse(b.err)
	}

	aURL := b.generateURL()
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

func (b *RequestBuilder) Pipe(writer io.Writer) error {
	return b.Do().Pipe(writer)
}

func (b *RequestBuilder) File(filename string) error {
	return b.Do().SaveFile(filename)
}

func (b *RequestBuilder) generateURL() string {
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
