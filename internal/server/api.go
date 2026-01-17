package server

import (
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/server/mapper"
	"github.com/kibirisu/borg/internal/util"
)

// PostAuthRegister implements api.ServerInterface.
func (s *Server) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var form api.AuthForm
	if err := util.ReadJSON(r, &form); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
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

// GetApiAccountsLookup implements api.ServerInterface.
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	acct := params.Acct
	handle, err := util.ParseHandle(acct, s.conf.ListenHost)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = handle
	panic("unimplemented")
	// if handle.Local {
	// 	log.Printf("lookup: local handle %s detected", acct)
	// 	account, err := s.service.App.GetLocalAccount(r.Context(), handle.Username)
	// 	if err != nil {
	// 		log.Println(err)
	// 		util.WriteError(w, http.StatusNotFound, err.Error())
	// 		return
	// 	}
	// 	log.Printf("lookup: found local account %s", account.Username)
	// 	util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
	// 	return
	// }
	//
	// if handle.Domain == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	//
	// log.Printf("lookup: remote handle %s detected, checking local cache", acct)
	// account, err := s.service.App.GetAccount(
	// 	r.Context(),
	// 	db.GetAccountParams{
	// 		Username: handle.Username,
	// 		Domain:   sql.NullString{String: handle.Domain, Valid: true},
	// 	},
	// )
	// if err == nil {
	// 	log.Printf("lookup: remote account %s@%s found locally", handle.Username, handle.Domain)
	// 	util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
	// 	return
	// }
	// if !errors.Is(err, sql.ErrNoRows) {
	// 	log.Println(err)
	// 	util.WriteError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	//
	// log.Printf(
	// 	"lookup: remote account %s@%s not cached, performing WebFinger lookup",
	// 	handle.Username,
	// 	handle.Domain,
	// )
	// util.WriteError(w, http.StatusInternalServerError, "unimplemented")
	// actor, err := s.service.Federation.processor.LookupActor(r.Context(), handle)
	// if err != nil {
	// 	log.Println(err)
	// 	util.WriteError(w, http.StatusBadGateway, err.Error())
	// 	return
	// }
	// if err != nil {
	// 	log.Println(err)
	// 	util.WriteError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// log.Printf("lookup: remote actor stored with username=%s domain=%s", row.Username, row.Domain.String)
	// util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(row))
}

// DeleteApiUsersId implements api.ServerInterface.
func (s *Server) DeleteApiUsersId(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PutApiUsersId implements api.ServerInterface.
func (s *Server) PutApiUsersId(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// DeleteApiPostsId implements api.ServerInterface.
func (s *Server) DeleteApiPostsId(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// GetApiPostsId implements api.ServerInterface.
func (s *Server) GetApiPostsId(w http.ResponseWriter, r *http.Request, id string) {
	info, err := s.service.App.GetPostByIDWithMetadata(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSON(w, http.StatusOK, *mapper.PostToAPIWithMetadata(&info.Status,
		&info.Account,
		int(info.LikeCount),
		int(info.ShareCount),
		int(info.CommentCount)))
}

// GetApiPostsIdComments implements api.ServerInterface.
func (s *Server) GetApiPostsIdComments(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PostApiPostsIdComments implements api.ServerInterface.
func (s *Server) PostApiPostsIdComments(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PostApiPostsIdLikes implements api.ServerInterface.
func (s *Server) PostApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PostApiPostsIdShares implements api.ServerInterface.
func (s *Server) PostApiPostsIdShares(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PutApiPostsId implements api.ServerInterface.
func (s *Server) PutApiPostsId(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented")
}

// PostApiAuthRegister implements api.ServerInterface.
func (s *Server) PostApiAuthRegister(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostApiAuthLogin implements api.ServerInterface.
func (s *Server) PostApiAuthLogin(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// GetApiUsersIdFavourites implements api.ServerInterface.
func (s *Server) GetApiUsersIdFavourites(w http.ResponseWriter, r *http.Request, id string) {
	posts, err := s.service.App.GetLikedPostsByAccountId(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	apiPosts := make([]api.Post, 0, len(posts))
	for _, info := range posts {
		converted := mapper.PostToAPIWithMetadata(
			&info.Status,
			&info.Account,
			int(info.LikeCount),
			int(info.ShareCount),
			int(info.CommentCount))
		apiPosts = append(apiPosts, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiPosts)
}

// GetApiUsersIdReblogged implements api.ServerInterface.
func (s *Server) GetApiUsersIdReblogged(w http.ResponseWriter, r *http.Request, id string) {
	posts, err := s.service.App.GetSharedPostsByAccountId(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	apiPosts := make([]api.Post, 0, len(posts))
	for _, info := range posts {
		converted := mapper.PostToAPIWithMetadata(
			&info.Status,
			&info.Account,
			int(info.LikeCount),
			int(info.ShareCount),
			int(info.CommentCount))
		apiPosts = append(apiPosts, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiPosts)
}

// GetApiUsersIdTimeline implements api.ServerInterface.
func (s *Server) GetApiUsersIdTimeline(w http.ResponseWriter, r *http.Request, id string) {
	posts, err := s.service.App.GetTimelinePostsByAccountId(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	apiPosts := make([]api.Post, 0, len(posts))
	for _, info := range posts {
		converted := mapper.PostToAPIWithMetadata(
			&info.Status,
			&info.Account,
			int(info.LikeCount),
			int(info.ShareCount),
			int(info.CommentCount))
		apiPosts = append(apiPosts, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiPosts)
}
