package deer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func WrapResponse(r *http.Response) *Response {
	return &Response{Raw: r}
}

type Response struct {
	read bool
	Raw  *http.Response
}

func (r *Response) ReadCloser() io.ReadCloser {
	if r.read {
		return nil
	}
	return r.Raw.Body
}

func (r *Response) Bytes() ([]byte, error) {
	if r.read {
		return nil, errors.New("deer: http response: body has read")
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	content, err := ioutil.ReadAll(r.Raw.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *Response) Text() (string, error) {
	if r.read {
		return "", errors.New("deer: http response: body has read")
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	bs, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (r *Response) BindWithJSON(value interface{}) error {
	if r.read {
		return errors.New("deer: http response: body has read")
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	return json.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *Response) BindWithXML(value interface{}) error {
	if r.read {
		return errors.New("deer: http response: body has read")
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	return xml.NewDecoder(r.Raw.Body).Decode(value)
}
