// internal/transport/http/router.go
package httpi

import (
	"net/http"
)

type Router struct{ mux *http.ServeMux }

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.mux.HandleFunc(pattern, h)
}

func (r *Router) ListenAndServe(path string) {
	http.ListenAndServe(path, r.mux)
}
