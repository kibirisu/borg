package server

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/web"
)

var _ api.ServerInterface = (*Server)(nil)

type Server struct {
	ds     domain.DataStore
	assets fs.FS
	conf   *config.Config
}

// GetApiAccountsLookup implements api.ServerInterface.
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	// we must check if account is local or from other instance
	// if from other instance we do webfinger lookup
	acct := params.Acct
	log.Println(acct)
	arr := strings.Split(acct, "@")
	username := arr[0]
	domain := arr[1]

	if domain == s.conf.ListenHost {
		actor, err := s.ds.Raw().GetActor(r.Context(), username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		account := api.Account{
			Acct:        acct,
			DisplayName: actor.DisplayName.String,
			Id:          int(actor.ID),
			Url:         actor.Url,
			Username:    actor.Username,
		}
		json.NewEncoder(w).Encode(&account)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else if domain != "" {
		actor, err := s.ds.Raw().GetAccount(r.Context(), db.GetAccountParams{username, sql.NullString{domain, true}})
		if err != nil {
			// we should do webfinger lookup at this point

			client := http.Client{Timeout: 2 * time.Second}
			req, err := http.NewRequest("GET", "http://"+domain+"/.well-known/webfinger", nil)
			q := req.URL.Query()
			q.Set("resource", acct)
			req.URL.RawQuery = q.Encode()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var webfinger api.WebFingerResponse
			if err = json.NewDecoder(req.Body).Decode(&webfinger); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_ = resp.Body.Close()

			// at this point we successfully looked up a account
			// and we should ask the other server for actor associated with the account

			req, err = http.NewRequest("GET", webfinger.Links[0].Href, nil)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp, err = client.Do(req)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var actor Actor
			if err = json.NewDecoder(resp.Body).Decode(&actor); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(actor)
			// we fetched remote actor
			// we must store it in database and return account in response
			_ = resp.Body.Close()
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		account := api.Account{
			Acct:        acct,
			DisplayName: actor.DisplayName.String,
			Id:          int(actor.ID),
			Url:         actor.Url,
			Username:    actor.Username,
		}
		json.NewEncoder(w).Encode(&account)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

type Actor struct {
	Context           any    `json:"@context"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
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
		r.Get("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
			// needed for actor identification before we can even follow one

			username := chi.URLParam(r, "username")
			log.Println(username)
			actor, err := server.ds.Raw().GetActor(r.Context(), username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(actor)
			// we need to build AP response here
			object := Actor{
				Context:           "https://www.w3.org/ns/activitystreams",
				ID:                actor.Uri,
				Type:              "Person",
				PreferredUsername: actor.Username,
				Inbox:             actor.InboxUri,
				Outbox:            actor.OutboxUri,
				Following:         actor.FollowingUri,
				Followers:         actor.FollowersUri,
			}
			w.Header().Set("Content-Type", "application/activity+json")
			json.NewEncoder(w).Encode(&object)
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/user/{username}/inbox", func(w http.ResponseWriter, r *http.Request) {
			username := chi.URLParam(r, "username")
			log.Println(username)
		})
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
