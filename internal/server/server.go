package server

import (
	"io/fs"
	"net/http"

	"borg/internal/api"
	"borg/internal/config"
	"borg/internal/domain"
	"borg/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct {
	ds     domain.DataStore
	assets fs.FS
	conf   *config.Config
}

func NewServer(conf *config.Config, ds domain.DataStore) *http.Server {
	assets, err := web.GetAssets()
	if err != nil {
		panic(err)
	}
	server := &Server{
		ds,
		assets,
		conf,
	}
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(server.createAuthMiddleware())
	r.Route("/", func(r chi.Router) {
		r.Get("/*", server.handleRoot)
		r.Get("/static/*", server.handleAssets)
	})
	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:" + conf.ListenPort,
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

// GetApiUsersIdPosts implements api.ServerInterface.
func (s *Server) GetApiUsersIdPosts(w http.ResponseWriter, r *http.Request, id int) {
	getByUserId(s.ds.PostRepository(), id).ServeHTTP(w, r)
}

// PostApiAuthRegister implements api.ServerInterface.
func (s *Server) PostApiAuthRegister(w http.ResponseWriter, r *http.Request) {
	registerUser(s.ds.UserRepository()).ServeHTTP(w, r)
}

// PostApiAuthLogin implements api.ServerInterface.
func (s *Server) PostApiAuthLogin(w http.ResponseWriter, r *http.Request) {
	loginUser(s.ds.UserRepository()).ServeHTTP(w, r)
}
