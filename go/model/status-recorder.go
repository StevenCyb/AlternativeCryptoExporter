package model

import "net/http"

// StatusRecorder to record the status code from the ResponseWriter
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader writes the status code to header and record
func (rec *StatusRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}
