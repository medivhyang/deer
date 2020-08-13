package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

func WrapResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{Raw: w}
}

type ResponseWriter struct {
	Raw http.ResponseWriter
}

func (w *ResponseWriter) SetHeader(key string, value string) *ResponseWriter {
	w.Raw.Header().Set(key, value)
	return w
}

func (w *ResponseWriter) SetStatusCode(statusCode int) {
	w.Raw.WriteHeader(statusCode)
}

func (w *ResponseWriter) Text(statusCode int, text string) error {
	w.Raw.Header().Set("Content-Type", "text/plain")
	w.Raw.WriteHeader(statusCode)
	_, err := io.WriteString(w.Raw, text)
	return err
}

func (w *ResponseWriter) HTML(statusCode int, content string) error {
	w.Raw.Header().Set("Content-Type", "text/html")
	w.Raw.WriteHeader(statusCode)
	_, err := io.WriteString(w.Raw, content)
	return err
}

func (w *ResponseWriter) JSON(statusCode int, value interface{}) error {
	w.Raw.Header().Set("Content-Type", "application/json")
	w.Raw.WriteHeader(statusCode)
	return json.NewEncoder(w.Raw).Encode(value)
}

func (w *ResponseWriter) XML(statusCode int, value interface{}) error {
	w.Raw.Header().Set("Content-Type", "application/xml")
	w.Raw.WriteHeader(statusCode)
	return xml.NewEncoder(w.Raw).Encode(value)
}
