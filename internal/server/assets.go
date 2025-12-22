package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) staticRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/*", s.serveFile("index.html"))
		r.Get("/static/*", s.handleAssets)
		r.Get("/api/doc", s.serveFile("docs.html"))
	}
}

func (s *Server) serveFile(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, s.assets, file)
	}
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(s.assets).ServeHTTP(w, r)
}
