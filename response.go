package deer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

var ErrResponseBodyHasRead = errors.New("deer: http response: body has read")

type Response interface {
	Raw() (*http.Response, error)
	Stream() (io.ReadCloser, error)
	Bytes() ([]byte, error)
	Text() (string, error)
	JSON(value interface{}) error
	XML(value interface{}) error
	Write(writer io.Writer) error
	WriteFile(filename string) error
	Dump(body bool) ([]byte, error)
}

func WrapResponse(r *http.Response) Response {
	return &response{raw: r}
}

type response struct {
	raw  *http.Response
	read bool
}

func (r *response) Raw() (*http.Response, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	return r.raw, nil
}

func (r *response) Dump(body bool) ([]byte, error) {
	return httputil.DumpResponse(r.raw, body)
}

func (r *response) Write(writer io.Writer) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	if _, err := io.Copy(writer, r.raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *response) WriteFile(filename string) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, r.raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *response) Stream() (io.ReadCloser, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	return r.raw.Body, nil
}

func (r *response) Bytes() ([]byte, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	content, err := ioutil.ReadAll(r.raw.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *response) Text() (string, error) {
	if r.read {
		return "", ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	bs, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (r *response) JSON(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return json.NewDecoder(r.raw.Body).Decode(value)
}

func (r *response) XML(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return xml.NewDecoder(r.raw.Body).Decode(value)
}

type errResponse struct {
	err error
}

func ErrorResponse(err error) Response {
	if err == nil {
		err = errors.New("deer: err response: unspecified error")
	}
	return &errResponse{err: err}
}

func (r *errResponse) Raw() (*http.Response, error) {
	return nil, r.err
}

func (r *errResponse) Stream() (io.ReadCloser, error) {
	return nil, r.err
}

func (r *errResponse) Bytes() ([]byte, error) {
	return nil, r.err
}

func (r *errResponse) Text() (string, error) {
	return "", r.err
}

func (r *errResponse) BindWithJSON(value interface{}) error {
	return r.err
}

func (r *errResponse) BindWithXML(value interface{}) error {
	return r.err
}

func (r *errResponse) JSON(value interface{}) error {
	return r.err
}

func (r *errResponse) XML(value interface{}) error {
	return r.err
}

func (r *errResponse) Write(writer io.Writer) error {
	return r.err
}

func (r *errResponse) WriteFile(filename string) error {
	return r.err
}

func (r *errResponse) Dump(body bool) ([]byte, error) {
	return nil, r.err
}
