package deer

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func WrapResponse(r *http.Response) *ResponseAdapter {
	return &ResponseAdapter{Raw: r}
}

type ResponseAdapter struct {
	Raw *http.Response
}

func (this *ResponseAdapter) Bytes() ([]byte, error) {
	defer this.Raw.Body.Close()
	content, err := ioutil.ReadAll(this.Raw.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (this *ResponseAdapter) Text() (string, error) {
	defer this.Raw.Body.Close()
	bs, err := this.Bytes()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (this *ResponseAdapter) BindWithJSON(value interface{}) error {
	defer this.Raw.Body.Close()
	return json.NewDecoder(this.Raw.Body).Decode(value)
}

func (this *ResponseAdapter) BindWithXML(value interface{}) error {
	defer this.Raw.Body.Close()
	return xml.NewDecoder(this.Raw.Body).Decode(value)
}
