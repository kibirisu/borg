package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/server/mapper"
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
	account, err := s.service.App.GetLocalAccount(r.Context(), user)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	util.WriteActivityJSON(w, http.StatusOK, mapper.AccountToActor(account))
}

func (s *Server) handleInbox(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	var activity domain.Activity
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "read error")
		return
	}

	if err = json.Unmarshal(bodyBytes, &activity); err != nil {
		return
	}

	switch activity.Type {
	case "Follow":
		var followReq domain.Follow
		if err = json.Unmarshal(bodyBytes, &followReq); err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		s.handleFollow(w, r, username, followReq)
	case "Create":
		var note domain.Note
		if err = json.Unmarshal(bodyBytes, &note); err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		// My goal is to build more generic / less concrete service functions
		// that will usually perform multiple repository operations.
		// At all cost, I would avoid just proxying repostiroy functions or explore more advanced sql queries.
		actor, err := s.service.App.GetAccount(r.Context(), db.GetAccountParams{
			Username: note.AttributedTo,
			Domain: sql.NullString{
				String: r.RemoteAddr,
				Valid:  true,
			}, // FIXME: RemoveAddr is shall be set to smth like <Host>:<Port>, we need to decide where servers accept TCP
		})
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if _, err = s.service.App.AddNote(r.Context(), db.CreateStatusParams{
			Uri:         note.ID,
			Url:         "TODO",
			Local:       sql.NullBool{Bool: false, Valid: true},
			Content:     note.Content,
			AccountID:   actor.ID,
			InReplyToID: sql.NullInt32{},
			ReblogOfID:  sql.NullInt32{},
		}); err != nil {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		w.WriteHeader(http.StatusCreated)
	case "Accept":
		log.Printf("Acknowledge activity: %s\n", activity.Type)
		w.WriteHeader(http.StatusAccepted)

	default:
		log.Printf("Received unsupported activity type: %s\n", activity.Type)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) fetchRemoteActor(ctx context.Context, uri string) (domain.Actor, error) {
	var actor domain.Actor

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

func (s *Server) resolveActor(ctx context.Context, raw json.RawMessage) domain.Actor {
	var actor domain.Actor
	if err := json.Unmarshal(raw, &actor); err == nil && actor.ID != "" {
		return actor
	}
	var uri string
	// not tested yet
	if err := json.Unmarshal(raw, &uri); err == nil && uri != "" {
		fetchedActor, err := s.fetchRemoteActor(ctx, uri)
		if err != nil {
			log.Printf("Warning: could not fetch remote actor %s: %v\n", uri, err)
			actor.ID = uri
			return actor
		}
		return fetchedActor
	}

	return actor
}

func (s *Server) handleFollow(
	w http.ResponseWriter,
	r *http.Request,
	localUsername string,
	msg domain.Follow,
) {
	followerActor := s.resolveActor(r.Context(), msg.Actor)
	followeeActor := s.resolveActor(r.Context(), msg.Object)

	if followeeActor.PreferredUsername != localUsername {
		util.WriteError(w, http.StatusNotFound, "conflicting information")
		return
	}
	localAccount, err := s.service.App.GetLocalAccount(r.Context(), localUsername)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "user not found")
		return
	}
	u, err := url.Parse(followerActor.ID)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, "domain could not be defined")
		return
	}
	originDomain := u.Host

	addRemoteAccount := mapper.ActorToAccountCreate(&followerActor, originDomain)
	followerAccount, err := s.service.App.AddRemoteAccount(r.Context(), addRemoteAccount)
	if err != nil {
		util.WriteError(
			w,
			http.StatusInternalServerError,
			"Error when adding remote actor: "+err.Error(),
		)
		return
	}
	createFollow := db.CreateFollowParams{
		Uri:             "",
		AccountID:       followerAccount.ID,
		TargetAccountID: localAccount.ID,
	}
	_, err = s.service.App.CreateFollow(r.Context(), &createFollow)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Error when adding follow: "+err.Error())
		return
	}

	log.Printf("User %s received follow from %s\n", localUsername, followerAccount.Username)

	accept := domain.Accept{
		Context: "https://www.w3.org/ns/activitystreams",
		ID:      fmt.Sprintf("%s/accept/%d", localAccount.Uri, time.Now().Unix()),
		Type:    "Accept",
		Actor:   localAccount.Uri,
		Object:  msg,
	}
	// TODO deliver
	_ = accept

	w.WriteHeader(http.StatusAccepted)
}
