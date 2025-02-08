package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		fmt.Printf("\n\nURL: %s\nReq: %+v\n\n", r.RequestURI, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func ParseForm(r *http.Request) map[string]string {
	r.ParseForm()
	formContents := map[string]string{}

	for key, value := range r.Form {
		formContents[key] = value[0]
	}
	return formContents
}

func ServeWebpage(path ...string) *template.Template {
	finalPath := filepath.Join(path...)
	return template.Must(template.ParseFiles(finalPath))
}
