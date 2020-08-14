package deer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
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

func (w *ResponseWriter) Text(statusCode int, text string) {
	w.Raw.Header().Set("Content-Type", "text/plain")
	w.Raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.Raw, text); err != nil {
		http.Error(w.Raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *ResponseWriter) HTML(statusCode int, content string) {
	w.Raw.Header().Set("Content-Type", "text/html")
	w.Raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.Raw, content); err != nil {
		http.Error(w.Raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *ResponseWriter) JSON(statusCode int, value interface{}) {
	w.Raw.Header().Set("Content-Type", "application/json")
	w.Raw.WriteHeader(statusCode)
	if value != nil {
		if  err := json.NewEncoder(w.Raw).Encode(value); err != nil {
			http.Error(w.Raw, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	reflectType := reflect.TypeOf(value)
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	var nilContent string
	switch reflectType.Kind() {
	case reflect.Slice:
		nilContent = "[]"
	case reflect.Struct, reflect.Map:
		nilContent = "{}"
	default:
		nilContent = "null"
	}
	if _, err := io.WriteString(w.Raw, nilContent); err != nil {
		http.Error(w.Raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *ResponseWriter) XML(statusCode int, value interface{}) {
	w.Raw.Header().Set("Content-Type", "application/xml")
	w.Raw.WriteHeader(statusCode)
	if err := xml.NewEncoder(w.Raw).Encode(value); err != nil {
		http.Error(w.Raw, err.Error(), http.StatusInternalServerError)
	}
}

func (w *ResponseWriter) Error(statusCode int, error string) {
	http.Error(w.Raw, error, statusCode)
}