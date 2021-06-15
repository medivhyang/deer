package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

var ErrResponseBodyHasRead = newError("http response", "body has read")

func WrapResponse(r *http.Response) *Response {
	return &Response{Raw: r}
}

type Response struct {
	Raw  *http.Response
	read bool
}

func (r *Response) Dump(body bool) ([]byte, error) {
	return httputil.DumpResponse(r.Raw, body)
}

func (r *Response) Pipe(writer io.Writer) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	if _, err := io.Copy(writer, r.Raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *Response) SaveFile(filename string) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, r.Raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *Response) Stream() (io.ReadCloser, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	return r.Raw.Body, nil
}

func (r *Response) Bytes() ([]byte, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
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
		return "", ErrResponseBodyHasRead
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

func (r *Response) JSON(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	return json.NewDecoder(r.Raw.Body).Decode(value)
}

func (r *Response) XML(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.Raw.Body.Close()
		r.read = true
	}()
	return xml.NewDecoder(r.Raw.Body).Decode(value)
}
