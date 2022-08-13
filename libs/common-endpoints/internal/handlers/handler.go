package handlers

import "net/http"

// Router defines router behavior.
type Router interface {
	HandleFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request))
}

// Result standard result for the service.
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}
