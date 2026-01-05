package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/util"
)

func (s *Server) federationRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/user/{username}", func(r chi.Router) {
			r.Get("/", s.handleGetActor)
			r.Post("/inbox", s.handleInbox)
		})
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

func (s *Server) handleInbox(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	_ = username

	var object domain.ObjectOrLink
	if err := util.ReadJSON(r, &object); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := s.service.Federation.ProcessInbox(r.Context(), &object); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) fetchRemoteActor(ctx context.Context, uri string) (domain.ActorOld, error) {
	var actor domain.ActorOld

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return actor, err
	}
	req.Header.Set("Accept", "application/activity+json")

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return actor, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return actor, fmt.Errorf("remote server returned status %d", resp.StatusCode)
	}

	if err = util.ReadJSON(req, &actor); err != nil {
		return actor, fmt.Errorf("failed to decode actor: %w", err)
	}

	return actor, nil
}
