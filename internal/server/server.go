package server

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/web"
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Get("/*", server.serveFile("index.html"))
		r.Get("/static/*", server.handleAssets)
		r.Get("/api/docs", server.serveFile("docs.html"))
	})
	// API routes muszą być przed catch-all route
	h := api.HandlerWithOptions(
		server,
		api.ChiServerOptions{
			BaseRouter:  r,
			Middlewares: []api.MiddlewareFunc{server.createAuthMiddleware()},
		},
	)

	s := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              "0.0.0.0:" + conf.ListenPort,
	}
	return s
}

func (s *Server) serveFile(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, s.assets, file)
	}
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

// GetApiPostsIdComments implements api.ServerInterface.
func (s *Server) GetApiPostsIdComments(w http.ResponseWriter, r *http.Request, id int) {
	getByPostID(s.ds.CommentRepository(), id).ServeHTTP(w, r)
}

// PostApiPostsIdComments implements api.ServerInterface.
func (s *Server) PostApiPostsIdComments(w http.ResponseWriter, r *http.Request, id int) {
	create(s.ds.CommentRepository()).ServeHTTP(w, r)
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
	getByUserID(s.ds.PostRepository(), id).ServeHTTP(w, r)
}

// PostApiAuthRegister implements api.ServerInterface.
func (s *Server) PostApiAuthRegister(w http.ResponseWriter, r *http.Request) {
	registerUser(s.ds.UserRepository()).ServeHTTP(w, r)
}

// PostApiAuthLogin implements api.ServerInterface.
func (s *Server) PostApiAuthLogin(w http.ResponseWriter, r *http.Request) {
	loginUser(s.ds.UserRepository()).ServeHTTP(w, r)
}

// GetApiUsersIdFollowers implements api.ServerInterface.
func (s *Server) GetApiUsersIdFollowers(w http.ResponseWriter, r *http.Request, id int) {
	getFollowers(s.ds.UserRepository(), id).ServeHTTP(w, r)
}

// GetApiUsersIdFollowing implements api.ServerInterface.
func (s *Server) GetApiUsersIdFollowing(w http.ResponseWriter, r *http.Request, id int) {
	getFollowing(s.ds.UserRepository(), id).ServeHTTP(w, r)
}

// GetApiPosts implements api.ServerInterface.
func (s *Server) GetApiPosts(w http.ResponseWriter, r *http.Request) {
	getAll(s.ds.PostRepository()).ServeHTTP(w, r)
}
