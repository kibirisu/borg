package server

import (
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/util"
)

// GetWellKnownWebfinger implements api.ServerInterface.
func (s *Server) GetWellKnownWebfinger(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetWellKnownWebfingerParams,
) {
	webfinger, err := s.service.App.WebfingerAccount(r.Context(), params)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.WriteWebFingerJSON(w, http.StatusOK, webfinger)
}

// PostAuthRegister implements api.ServerInterface.
func (s *Server) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var form api.AuthForm
	if err := util.ReadJSON(r, &form); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("auth register: incoming request username=%s", form.Username)

	if err := s.service.App.Register(r.Context(), form); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("auth register: user %s created successfully", form.Username)
	w.WriteHeader(http.StatusCreated)
}

// PostAuthLogin implements api.ServerInterface.
func (s *Server) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var form api.AuthForm
	if err := util.ReadJSON(r, &form); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.service.App.Login(r.Context(), form)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Authorization", "Bearer: "+token)
	w.WriteHeader(http.StatusOK)
}

// GetApiAccountsId implements api.ServerInterface.
func (s *Server) GetApiAccountsId(w http.ResponseWriter, r *http.Request, id string) {
	account, err := s.service.App.GetAccount(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, account)
}

// GetApiAccountsIdStatuses implements api.ServerInterface.
func (s *Server) GetApiAccountsIdStatuses(w http.ResponseWriter, r *http.Request, id string) {
	statuses, err := s.service.App.GetAccountStatuses(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, statuses)
}

// GetApiAccountsIdFollowers implements api.ServerInterface.
func (s *Server) GetApiAccountsIdFollowers(w http.ResponseWriter, r *http.Request, id string) {
	followers, err := s.service.App.GetAccountFollowers(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, followers)
}

// GetApiAccountsIdFollowing implements api.ServerInterface.
func (s *Server) GetApiAccountsIdFollowing(w http.ResponseWriter, r *http.Request, id string) {
	following, err := s.service.App.GetAccountFollowing(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, following)
}

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.FollowAccount(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// PostApiAccountsIdUnfollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdUnfollow(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.UnfollowAccount(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// GetApiAccountsLookup implements api.ServerInterface.
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	account, err := s.service.App.LookupAccount(r.Context(), params.Acct)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, account)
}

// PostApiStatuses implements api.ServerInterface.
func (s *Server) PostApiStatuses(w http.ResponseWriter, r *http.Request) {
	var status api.PostApiStatusesJSONBody
	if err := util.ReadJSON(r, &status); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	job, err := s.service.App.CreateStatus(r.Context(), status)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// GetApiStatusesId implements api.ServerInterface.
func (s *Server) GetApiStatusesId(w http.ResponseWriter, r *http.Request, id string) {
	status, err := s.service.App.ViewStatus(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
	}
	util.WriteJSON(w, http.StatusOK, status)
}

// PostApiStatusesIdFavourite implements api.ServerInterface.
func (s *Server) PostApiStatusesIdFavourite(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.FavouriteStatus(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// PostApiStatusesIdUnfavourite implements api.ServerInterface.
func (s *Server) PostApiStatusesIdUnfavourite(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.UnfavouriteStatus(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// PostApiStatusesIdReblog implements api.ServerInterface.
func (s *Server) PostApiStatusesIdReblog(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.ReblogStatus(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// PostApiStatusesIdUnreblog implements api.ServerInterface.
func (s *Server) PostApiStatusesIdUnreblog(w http.ResponseWriter, r *http.Request, id string) {
	job, err := s.service.App.UnreblogStatus(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.worker.Enqueue(job)
}

// GetApiStatusesIdReplies implements api.ServerInterface.
func (s *Server) GetApiStatusesIdReplies(w http.ResponseWriter, r *http.Request, id string) {
	replies, err := s.service.App.GetStatusReplies(r.Context(), id)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, replies)
}

// GetApiTimelinesHome implements api.ServerInterface.
func (s *Server) GetApiTimelinesHome(w http.ResponseWriter, r *http.Request) {
	statuses, err := s.service.App.ViewHomeTimeline(r.Context())
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, statuses)
}
// GetApiTimelinesFavourite implements api.ServerInterface.
func (s *Server) GetApiTimelinesFavourite(w http.ResponseWriter, r *http.Request) {
	statuses, err := s.service.App.ViewFavouriteTimeline(r.Context())
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, statuses)
}
