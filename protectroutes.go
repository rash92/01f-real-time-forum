package main

import (
	"forum/controller"
	"net/http"
)

func protectGetRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			controller.PageErrors(w, r, tmpl, 400, "Bad Request")
			return
		}
		h(w, r)

	}
}

func protectPostRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			controller.PageErrors(w, r, tmpl, 400, "Bad Request")
			return
		}
		h(w, r)
	}
}

func protectPostGetRequests(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "GET" {
			controller.PageErrors(w, r, tmpl, 400, "Bad requests")
			return
		}
		h(w, r)
	}
}
