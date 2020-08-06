package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

func WrapResponseWriter(w http.ResponseWriter) *ResponseWriterAdapter {
	return &ResponseWriterAdapter{Raw: w}
}

type ResponseWriterAdapter struct {
	Raw http.ResponseWriter
}

func (this *ResponseWriterAdapter) Text(statusCode int, text string) error {
	this.Raw.Header().Set("Content-Type", "text/plain")
	this.Raw.WriteHeader(statusCode)
	_, err := io.WriteString(this.Raw, text)
	return err
}

func (this *ResponseWriterAdapter) HTML(statusCode int, content string) error {
	this.Raw.Header().Set("Content-Type", "text/html")
	this.Raw.WriteHeader(statusCode)
	_, err := io.WriteString(this.Raw, content)
	return err
}

func (this *ResponseWriterAdapter) JSON(statusCode int, value interface{}) error {
	this.Raw.Header().Set("Content-Type", "application/json")
	this.Raw.WriteHeader(statusCode)
	return json.NewEncoder(this.Raw).Encode(value)
}

func (this *ResponseWriterAdapter) XML(statusCode int, value interface{}) error {
	this.Raw.Header().Set("Content-Type", "application/xml")
	this.Raw.WriteHeader(statusCode)
	return xml.NewEncoder(this.Raw).Encode(value)
}
