package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/server/mapper"
	"github.com/kibirisu/borg/internal/util"
)

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

// GetApiAccountsLookup implements api.ServerInterface.
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	// // we must check if account is local or from other instance
	// // if from other instance we do webfinger lookup
	acct := params.Acct
	handle, err := util.ParseHandle(acct, s.conf.ListenHost)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if handle.Local {
		log.Printf("lookup: local handle %s detected", acct)
		account, err := s.service.App.GetLocalAccount(r.Context(), handle.Username)
		if err != nil {
			log.Println(err)
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("lookup: found local account %s", account.Username)
		util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
		return
	}

	if handle.Domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("lookup: remote handle %s detected, checking local cache", acct)
	account, err := s.service.App.GetAccount(
		r.Context(),
		db.GetAccountParams{
			Username: handle.Username,
			Domain:   sql.NullString{String: handle.Domain, Valid: true},
		},
	)
	if err == nil {
		log.Printf("lookup: remote account %s@%s found locally", handle.Username, handle.Domain)
		util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("lookup: remote account %s@%s not cached, performing WebFinger lookup", handle.Username, handle.Domain)
	actor, err := s.service.Federation.LookupRemoteActor(r.Context(), handle)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	log.Printf("lookup: persisting remote actor %s", actor.ID)
	row, err := s.service.Federation.CreateActor(
		r.Context(),
		*mapper.ActorToDB(actor, handle.Domain),
	)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("lookup: remote actor stored with username=%s domain=%s", row.Username, row.Domain.String)
	util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(row))
}

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id int) {
	container, ok := r.Context().Value("token").(*tokenContainer)
    
    if !ok || container == nil || container.id == nil {
        util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
        return
    }
    currentUserID := *container.id
	if currentUserID == id {
		http.Error(w, "Tried to follow oneself", http.StatusBadRequest)
		return
	}
	follower, err := s.service.App.GetAccountById(r.Context(), currentUserID)
	followee, err := s.service.App.GetAccountById(r.Context(), id)
	follow, err := s.service.App.FollowAccount(r.Context(), currentUserID, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	APfollow := mapper.DBToFollow(follow, &follower, &followee)
	log.Println(followee.InboxUri)
	if follower.Domain != followee.Domain {
		util.DeliverToEndpoint(followee.InboxUri, APfollow)
	}
	util.WriteJSON(w, http.StatusCreated, nil);
}

// DeleteApiUsersId implements api.ServerInterface.
func (s *Server) DeleteApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiUsersId implements api.ServerInterface.
func (s *Server) GetApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiUsers implements api.ServerInterface.
func (s *Server) PostApiUsers(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PutApiUsersId implements api.ServerInterface.
func (s *Server) PutApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// DeleteApiPostsId implements api.ServerInterface.
func (s *Server) DeleteApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiPostsId implements api.ServerInterface.
func (s *Server) GetApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiPostsIdComments implements api.ServerInterface.
func (s *Server) GetApiPostsIdComments(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPostsIdComments implements api.ServerInterface.
func (s *Server) PostApiPostsIdComments(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiPostsIdLikes implements api.ServerInterface.
func (s *Server) GetApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPostsIdLikes implements api.ServerInterface.
func (s *Server) PostApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiPostsIdShares implements api.ServerInterface.
func (s *Server) GetApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPostsIdShares implements api.ServerInterface.
func (s *Server) PostApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPosts implements api.ServerInterface.
func (s *Server) PostApiPosts(w http.ResponseWriter, r *http.Request) {
    container, ok := r.Context().Value("token").(*tokenContainer)
    if !ok || container == nil || container.id == nil {
        util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
        return
    }
    currentUserID := *container.id
    poster, err := s.service.App.GetAccountById(r.Context(), currentUserID)
    var newPost api.NewPost
    if err := util.ReadJSON(r, &newPost); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return 
    }
    var newDBPost = mapper.NewPostToDB(&newPost, true)
    status, err := s.service.App.AddNote(r.Context(), *newDBPost)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    s.service.App.DeliverToFollowers(w, r, newPost.UserID, func(recipientURI string) any {
        return mapper.PostToCreateNote(&status, &poster, []string{recipientURI})
    })
    util.WriteJSON(w, http.StatusCreated, nil);
}

// PutApiPostsId implements api.ServerInterface.
func (s *Server) PutApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiUsersIdPosts implements api.ServerInterface.
func (s *Server) GetApiUsersIdPosts(w http.ResponseWriter, r *http.Request, id int) {
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

// GetApiUsersIdFollowers implements api.ServerInterface.
func (s *Server) GetApiUsersIdFollowers(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiUsersIdFollowing implements api.ServerInterface.
func (s *Server) GetApiUsersIdFollowing(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiPosts implements api.ServerInterface.
func (s *Server) GetApiPosts(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
