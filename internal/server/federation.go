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
		r.Route("/users/{id}", func(r chi.Router) {
			r.Get("/", s.handleGetActor)
			r.Get("/followers", s.handleActorFollowers)
			r.Get("/following", s.handleActorFollowing)
			r.Get("/follow/{followId}", s.handleGetFollow)
			r.Get("/statuses/{statusId}", s.handleGetStatus)
			r.Get("/likes/{likeId}", s.handleGetLike)
			r.Get("/reblogs/{reblogId}", func(http.ResponseWriter, *http.Request) {})
			r.Get("/requests/{requestId}", func(http.ResponseWriter, *http.Request) {})
			r.Post("/inbox", s.handleInbox)
		})
	}
}

func (s *Server) handleGetActor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actor, err := s.service.Federation.GetActor(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, actor)
}

func (s *Server) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "statusId")
	status, err := s.service.Federation.GetStatus(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, status)
}

func (s *Server) handleGetLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "likeId")
	like, err := s.service.Federation.GetLike(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, like)
}

func (s *Server) handleGetFollow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "followId")
	follow, err := s.service.Federation.GetFollow(r.Context(), id)
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
	id := chi.URLParam(r, "id")

	var object domain.ObjectOrLink
	if err := util.ReadJSON(r, &object); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	job, err := s.service.Federation.ProcessIncoming(r.Context(), &object, id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.worker.Enqueue(job)
	w.WriteHeader(http.StatusAccepted)
}
