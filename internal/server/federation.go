package server

import (
	"encoding/json"
	"net/http"

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
	resp.Subject = params.Resource
	json.NewEncoder(w).Encode(&resp)
}

