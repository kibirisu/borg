package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/util"
)

// https://docs.joinmastodon.org/spec/webfinger/#intro
// not fully implemented, only things needed for basic actor resolving

// GetWellKnownWebfinger implements api.ServerInterface.
func (s *Server) GetWellKnownWebfinger(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetWellKnownWebfingerParams,
) {
	// Ok, so someone asks us about user's identifier
	// We should provide the identifier here
	var resp api.WebFingerResponse
	username, domain, err := util.ParseAccount(params.Resource)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if domain != s.conf.ListenHost {
		msg := fmt.Sprintf("domain mismatch: got %s, expected %s", domain, s.conf.ListenHost)
		util.WriteError(w, http.StatusNotFound, msg)
		return
	}

	// Here we should query database for a user contained in resource
	// And return minimal response as defined in WebFinger spec

	// FIME: be aware that below function unnecessarily map account object to activity
	account, err := s.service.App.GetLocalAccount(r.Context(), username)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
	}

	// WebFinger pointing to queried actor
	resp.Subject = params.Resource
	resp.Links = append(
		resp.Links,
		api.WebFingerLink{
			Rel:  "self",
			Type: "application/activity+json",
			Href: account.Uri,
		},
	)
	util.WriteWebFingerJSON(w, http.StatusOK, &resp)
}
