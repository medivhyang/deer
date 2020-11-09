package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type ResponseWriter interface {
	Raw() http.ResponseWriter
	SetHeader(key string, value string) *responseWriter
	SetStatusCode(statusCode int)
	Text(statusCode int, text string)
	HTML(statusCode int, content string)
	JSON(statusCode int, value interface{})
	XML(statusCode int, value interface{})
	Error(statusCode int, errorMessage ...string)
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

func (w *responseWriter) SetHeader(key string, value string) *responseWriter {
	w.raw.Header().Set(key, value)
	return w
}

func (w *responseWriter) SetStatusCode(statusCode int) {
	w.raw.WriteHeader(statusCode)
}

func (w *responseWriter) Text(statusCode int, text string) {
	w.raw.Header().Set("Content-Type", "text/plain")
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, text); err != nil {
		http.Error(w.raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *responseWriter) HTML(statusCode int, content string) {
	w.raw.Header().Set("Content-Type", "text/html")
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, content); err != nil {
		http.Error(w.raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *responseWriter) JSON(statusCode int, value interface{}) {
	w.raw.Header().Set("Content-Type", "application/json")
	w.raw.WriteHeader(statusCode)
	if err := json.NewEncoder(w.raw).Encode(value); err != nil {
		http.Error(w.raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *responseWriter) XML(statusCode int, value interface{}) {
	w.raw.Header().Set("Content-Type", "application/xml")
	w.raw.WriteHeader(statusCode)
	if err := xml.NewEncoder(w.raw).Encode(value); err != nil {
		http.Error(w.raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *responseWriter) Error(statusCode int, errorMessage ...string) {
	if len(errorMessage) > 0 {
		http.Error(w.raw, errorMessage[0], statusCode)
		return
	}
	http.Error(w.raw, http.StatusText(statusCode), statusCode)
}