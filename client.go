package deer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var DefaultClient = WrapClient(&http.Client{})

func Get(url string, query map[string]string, header map[string]string) (*ResponseAdapter, error) {
	return DefaultClient.Get(url, query, header)
}

func GetText(url string, query map[string]string, header map[string]string) (string, error) {
	return DefaultClient.GetText(url, query, header)
}

func GetJSON(url string, query map[string]string, header map[string]string, value interface{}) error {
	return DefaultClient.GetJSON(url, query, header, value)
}

func GetXML(url string, query map[string]string, header map[string]string, value interface{}) error {
	return DefaultClient.GetXML(url, query, header, value)
}

func Post(url string, body io.Reader, header map[string]string) (*ResponseAdapter, error) {
	return DefaultClient.Post(url, body, header)
}

func PostJSON(url string, value interface{}, header map[string]string) (*ResponseAdapter, error) {
	return DefaultClient.PostJSON(url, value, header)
}

func PostXML(url string, value interface{}, header map[string]string) (*ResponseAdapter, error) {
	return DefaultClient.PostXML(url, value, header)
}

func Do(method string, url string, body io.Reader, header map[string]string) (*ResponseAdapter, error) {
	return DefaultClient.Do(method, url, body, header)
}

func WrapClient(client *http.Client) *ClientAdapter {
	return &ClientAdapter{Raw: client}
}

type ClientAdapter struct {
	Raw *http.Client
}

func (client *ClientAdapter) Get(url string, query map[string]string, header map[string]string) (*ResponseAdapter, error) {
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += client.parseURLValues(query).Encode()
	return client.Do(http.MethodGet, url, nil, header)
}

func (client *ClientAdapter) GetText(url string, query map[string]string, header map[string]string) (string, error) {
	r, err := client.Get(url, query, header)
	if err != nil {
		return "", err
	}
	return r.Text()
}

func (client *ClientAdapter) GetJSON(url string, query map[string]string, header map[string]string, value interface{}) error {
	r, err := client.Get(url, query, header)
	if err != nil {
		return err
	}
	return r.BindWithJSON(value)
}

func (client *ClientAdapter) GetXML(url string, query map[string]string, header map[string]string, value interface{}) error {
	r, err := client.Get(url, query, header)
	if err != nil {
		return err
	}
	return r.BindWithXML(value)
}

func (client *ClientAdapter) Post(url string, body io.Reader, header map[string]string) (*ResponseAdapter, error) {
	return client.Do(http.MethodPost, url, body, header)
}

func (client *ClientAdapter) PostJSON(url string, value interface{}, header map[string]string) (*ResponseAdapter, error) {
	bs, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return client.Do(http.MethodPost, url, bytes.NewReader(bs), header)
}

func (client *ClientAdapter) PostXML(url string, value interface{}, header map[string]string) (*ResponseAdapter, error) {
	bs, err := xml.Marshal(value)
	if err != nil {
		return nil, err
	}
	return client.Do(http.MethodPost, url, bytes.NewReader(bs), header)
}

func (client *ClientAdapter) Do(method string, url string, body io.Reader, header map[string]string) (*ResponseAdapter, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Raw.Do(req)
	if err != nil {
		return nil, err
	}
	return WrapResponse(resp), nil
}

func (client *ClientAdapter) parseURLValues(m map[string]string) *url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, v)
	}
	return &values
}

func (client *ClientAdapter) parseHTTPHeader(m map[string]string) *http.Header {
	header := http.Header{}
	for k, v := range m {
		header.Set(k, v)
	}
	return &header
}
