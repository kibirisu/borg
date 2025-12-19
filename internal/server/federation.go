package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kibirisu/borg/internal/util"
)

func (s *Server) federationRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/user/{username}", func(r chi.Router) {
			// TODO: add middleware providing username context instead of manually extracting it with every endpoint handler
			// also, currently i cannot remember if there would be any ActivityPub endpoints not starting with "/user",
			// so just mounting subrouter under "user/{username}" would be easier
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				user := chi.URLParam(r, "username")
				actor, err := s.service.Federation.GetLocalActor(r.Context(), user)
				if err != nil {
					util.WriteError(w, http.StatusBadRequest, err.Error())
				}
				util.WriteActivityJSON(w, http.StatusOK, &actor)
			})
			r.Post("/inbox", func(w http.ResponseWriter, r *http.Request) {
				panic("unimplemented")
			})
		})
	}
}
