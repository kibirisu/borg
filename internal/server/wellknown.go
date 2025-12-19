package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/util"
)

// GetWellKnownWebfinger implements api.ServerInterface.
func (s *Server) GetWellKnownWebfinger(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetWellKnownWebfingerParams,
) {
	// Ok, so someone asks us about user's identifier
	// We should provide the identifier here
	// In ActivityStreams all objects "should have unique global identifiers"
	var resp api.WebFingerResponse
	resource := strings.TrimPrefix(params.Resource, "acct:")
	log.Println(resource, resp)

	// Here we should query database for a user contained in resource
	// And return minimal response as defined in WebFinger spec

	// we need helper function for parsing "acct"
	// TODO: extract username from user handle
	// FIME: be aware that below function unnecessarily map account object to activity
	actor, err := s.service.Federation.GetLocalActor(r.Context(), resource)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
	}

	// WebFinger pointing to queried actor
	resp.Subject = resource
	resp.Links = append(
		resp.Links,
		api.WebFingerLink{
			Rel:  "self",
			Type: "application/activity+json",
			Href: actor.ID,
		},
	)
	util.WriteWebFingerJSON(w, http.StatusOK, &resp)
}
