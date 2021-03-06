package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type ResponseWriter interface {
	Raw() http.ResponseWriter
	StatusCode(statusCode int)
	Header(key string, value string) ResponseWriter
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

func (w *responseWriter) Header(key string, value string) ResponseWriter {
	w.raw.Header().Set(key, value)
	return w
}

func (w *responseWriter) Text(statusCode int, text string) {
	w.raw.Header().Set("Content-Type", "text/plain")
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, text); err != nil {
		panic(err)
	}
}

func (w *responseWriter) HTML(statusCode int, content string) {
	w.raw.Header().Set("Content-Type", "text/html")
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, content); err != nil {
		panic(err)
	}
}

func (w *responseWriter) JSON(statusCode int, value interface{}) {
	w.raw.Header().Set("Content-Type", "application/json")
	w.raw.WriteHeader(statusCode)
	if err := json.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}

func (w *responseWriter) XML(statusCode int, value interface{}) {
	w.raw.Header().Set("Content-Type", "application/xml")
	w.raw.WriteHeader(statusCode)
	if err := xml.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}
