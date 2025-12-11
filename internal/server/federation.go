package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kibirisu/borg/internal/api"
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

	// Here we should query database for a user contained in resource
	// And return minimal response as defined in WebFinger spec

	// WebFinger pointing to queried actor
	resp.Subject = resource
	resp.Links = append(
		resp.Links,
		api.WebFingerLink{
			Rel:  "self",
			Type: "application/activity+json",
			Href: "http://127.0.0.1:8080/users/user", // ActivityPub user URI
		},
	)

	json.NewEncoder(w).Encode(&resp)
	w.WriteHeader(http.StatusOK)
}

