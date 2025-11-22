package server

import (
	"io/fs"
	"net/http"

	"borg/internal/api"
	"borg/internal/domain"
	"borg/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct {
	ds     domain.DataStore
	assets fs.FS
}

func NewServer(listenPort string, ds domain.DataStore) *http.Server {
	assets, err := web.GetAssets()
	if err != nil {
		panic(err)
	}
	server := &Server{
		ds:     ds,
		assets: assets,
	}
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Get("/*", server.handleRoot)
		r.Get("/static/*", server.handleAssets)
	})
	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:" + listenPort,
	}
	return s
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, s.assets, "index.html")
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(s.assets).ServeHTTP(w, r)
}

// DeleteApiUsersId implements api.ServerInterface.
func (s *Server) DeleteApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	deleteByID(s.ds.UserRepository(), id).ServeHTTP(w, r)
}

// GetApiUsersId implements api.ServerInterface.
func (s *Server) GetApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	getByID(s.ds.UserRepository(), id).ServeHTTP(w, r)
}

// GetApiUsersByUsernameUsername implements api.ServerInterface.
func (s *Server) GetApiUsersByUsernameUsername(w http.ResponseWriter, r *http.Request, username string) {
	getByUsername(s.ds.UserRepository(), username).ServeHTTP(w, r)
}

// PostApiUsers implements api.ServerInterface.
func (s *Server) PostApiUsers(w http.ResponseWriter, r *http.Request) {
	create(s.ds.UserRepository()).ServeHTTP(w, r)
}

// PutApiUsersId implements api.ServerInterface.
func (s *Server) PutApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	update(s.ds.UserRepository()).ServeHTTP(w, r)
}

// DeleteApiPostsId implements api.ServerInterface.
func (s *Server) DeleteApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	deleteByID(s.ds.PostRepository(), id).ServeHTTP(w, r)
}

// GetApiPostsId implements api.ServerInterface.
func (s *Server) GetApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	getByID(s.ds.PostRepository(), id).ServeHTTP(w, r)
}

// PostApiPosts implements api.ServerInterface.
func (s *Server) PostApiPosts(w http.ResponseWriter, r *http.Request) {
	create(s.ds.PostRepository()).ServeHTTP(w, r)
}

// PutApiPostsId implements api.ServerInterface.
func (s *Server) PutApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	update(s.ds.PostRepository()).ServeHTTP(w, r)
}

// GetApiPosts implements api.ServerInterface.
func (s *Server) GetApiPosts(w http.ResponseWriter, r *http.Request) {
	getAll(s.ds.PostRepository()).ServeHTTP(w, r)
}

// GetApiUsersIdPosts implements api.ServerInterface.
func (s *Server) GetApiUsersIdPosts(w http.ResponseWriter, r *http.Request, id int) {
	getByUserId(s.ds.PostRepository(), id).ServeHTTP(w, r)
}
