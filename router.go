package main

import (
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]http.HandlerFunc
}

func NewRouter() (*Router, error) {
	return &Router{
		routes: make(map[string]http.HandlerFunc),
	}, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for pattern, handler := range r.routes {
		if isMatch(pattern, req.URL.Path) {
			handler(w, req)
			return
		}
	}

	// No matching route found
	http.NotFound(w, req)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.routes[pattern] = handler
}

func isMatch(pattern, path string) bool {
	parts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(parts) != len(pathParts) {
		return false
	}

	for i, part := range parts {
		if !strings.HasPrefix(part, "{") && part != pathParts[i] {
			return false
		}
	}

	return true
}
