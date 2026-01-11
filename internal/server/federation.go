package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/util"
)

func (s *Server) federationRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/user/{username}", func(r chi.Router) {
			r.Get("/", s.handleGetActor)
			r.Get("/followers", s.handleActorFollowers)
			r.Get("/following", s.handleActorFollowing)
			r.Post("/inbox", s.handleInbox)
		})
		r.Get("/statuses/{id}", s.handleGetStatus)
		r.Get("/likes/{id}", s.handleGetLike)
		r.Get("/follow/{id}", s.handleGetFollow)
	}
}

func (s *Server) handleGetActor(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")
	actor, err := s.service.Federation.GetLocalActor(r.Context(), user)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, actor)
}

func (s *Server) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	uri := chi.URLParam(r, "id")
	status, err := s.service.Federation.GetStatus(r.Context(), uri)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, status)
}

func (s *Server) handleGetLike(w http.ResponseWriter, r *http.Request) {
	uri := chi.URLParam(r, "id")
	like, err := s.service.Federation.GetLike(r.Context(), uri)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, like)
}

func (s *Server) handleGetFollow(w http.ResponseWriter, r *http.Request) {
	uri := chi.URLParam(r, "id")
	follow, err := s.service.Federation.GetFollow(r.Context(), uri)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, follow)
}

func (s *Server) handleActorFollowers(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")

	var pagePtr *int = nil
	pageParam, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil {
		pagePtr = &pageParam
	}
	collection, err := s.service.Federation.GetActorFollowers(r.Context(), user, pagePtr)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, collection)
}

func (s *Server) handleActorFollowing(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")

	var pagePtr *int = nil
	pageParam, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil {
		pagePtr = &pageParam
	}
	collection, err := s.service.Federation.GetActorFollowing(r.Context(), user, pagePtr)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, collection)
}

func (s *Server) handleInbox(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	_ = username

	var object domain.ObjectOrLink
	if err := util.ReadJSON(r, &object); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	job, err := s.service.Federation.ProcessIncoming(r.Context(), &object)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
	}
	s.worker.Enqueue(job)
	w.WriteHeader(http.StatusAccepted)
}
