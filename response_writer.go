package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

const (
	headerContentType = "Content-Type"

	mimeText = "text/plain"
	mimeHTML = "text/html"
	mimeJSON = "application/json"
	mimeXML  = "application/xml"
)

type ResponseWriter interface {
	Raw() http.ResponseWriter
	StatusCode(statusCode int)
	Header(key string, value string) *responseWriter
	Text(statusCode int, text string)
	HTML(statusCode int, content string)
	JSON(statusCode int, value interface{})
	XML(statusCode int, value interface{})
}

type responseWriter struct {
	raw http.ResponseWriter
}

func WrapResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{raw: w}
}

func (w *responseWriter) Raw() http.ResponseWriter {
	return w.raw
}

func (w *responseWriter) StatusCode(statusCode int) {
	w.raw.WriteHeader(statusCode)
}

func (w *responseWriter) Header(key string, value string) *responseWriter {
	w.raw.Header().Set(key, value)
	return w
}

func (w *responseWriter) Text(statusCode int, text string) {
	w.raw.Header().Set(headerContentType, mimeText)
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, text); err != nil {
		panic(err)
	}
}

func (w *responseWriter) HTML(statusCode int, content string) {
	w.raw.Header().Set(headerContentType, mimeHTML)
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, content); err != nil {
		panic(err)
	}
}

func (w *responseWriter) JSON(statusCode int, value interface{}) {
	w.raw.Header().Set(headerContentType, mimeJSON)
	w.raw.WriteHeader(statusCode)
	if err := json.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}

func (w *responseWriter) XML(statusCode int, value interface{}) {
	w.raw.Header().Set(headerContentType, mimeXML)
	w.raw.WriteHeader(statusCode)
	if err := xml.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}
