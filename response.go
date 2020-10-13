package deer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

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
	Copy() Response
}

func WrapResponse(r *http.Response) Response {
	return &response{raw: r}
}

type response struct {
	raw  *http.Response
	read bool
}

func (r *response) Raw() (*http.Response, error) {
	if err := r.check(); err != nil {
		return nil, err
	}
	return r.raw, nil
}

func (r *response) Dump(body bool) ([]byte, error) {
	return httputil.DumpResponse(r.raw, body)
}

func (r *response) Copy() Response {
	if r.raw.Body == nil {
		return ErrorResponse(errors.New("deer: response copy: require body"))
	}
	buffer, err := ioutil.ReadAll(r.raw.Body)
	if err != nil {
		return ErrorResponse(err)
	}
	r.raw.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	return r
}

func (r *response) Write(writer io.Writer) error {
	if err := r.check(); err != nil {
		return err
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
	if err := r.check(); err != nil {
		return err
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
	if err := r.check(); err != nil {
		return nil, err
	}
	return r.raw.Body, nil
}

func (r *response) Bytes() ([]byte, error) {
	if err := r.check(); err != nil {
		return nil, err
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
	if err := r.check(); err != nil {
		return "", err
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
	if err := r.check(); err != nil {
		return err
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return json.NewDecoder(r.raw.Body).Decode(value)
}

func (r *response) XML(value interface{}) error {
	if err := r.check(); err != nil {
		return err
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return xml.NewDecoder(r.raw.Body).Decode(value)
}

func (r *response) check() error {
	if r.read {
		return errors.New("deer: http response: body has read")
	}
	return nil
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

func (r *errResponse) Copy() Response {
	return ErrorResponse(r.err)
}
