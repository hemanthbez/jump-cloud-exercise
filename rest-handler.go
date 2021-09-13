package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type ctxKey struct{}

var routes = []route{
	newRoute("GET", "/", defaultHandler),
	newRoute("GET", "/hash/([0-9]+)", processPasswordGetRequest),
	newRoute("GET", "/stats", processStats),
	newRoute("GET", "/shutdown", processShutown),
	newRoute("POST", "/hash", processPasswordPostRequest),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {

	if shutdownEnabled {
		http.Error(w, "403 Forbidden - Cannot accept any new requests!", http.StatusForbidden)
		return
	}

	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ERROR: Missing Data \n")
	return
}
