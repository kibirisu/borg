package server

import (
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/api"
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

	token, err := s.service.API.Login(r.Context(), form)
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

	if err := s.service.API.Register(r.Context(), form); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
}

// GetApiAccountsLookup implements api.ServerInterface.
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	// we must check if account is local or from other instance
	// if from other instance we do webfinger lookup
	// acct := params.Acct
	// log.Println(acct)
	// arr := strings.Split(acct, "@")
	// username := arr[0]
	// domain := arr[1]
	//
	// if domain == s.conf.ListenHost {
	// 	actor, err := s.ds.Raw().GetActor(r.Context(), username)
	// 	if err != nil {
	// 		log.Println(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	account := api.Account{
	// 		Acct:        acct,
	// 		DisplayName: actor.DisplayName.String,
	// 		Id:          int(actor.ID),
	// 		Url:         actor.Url,
	// 		Username:    actor.Username,
	// 	}
	// 	json.NewEncoder(w).Encode(&account)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// } else if domain != "" {
	// 	actor, err := s.ds.Raw().GetAccount(r.Context(), db.GetAccountParams{username, sql.NullString{domain, true}})
	// 	if err != nil {
	// 		// we should do webfinger lookup at this point
	//
	// 		client := http.Client{Timeout: 2 * time.Second}
	// 		req, err := http.NewRequest("GET", "http://"+domain+"/.well-known/webfinger", nil)
	// 		q := req.URL.Query()
	// 		q.Set("resource", acct)
	// 		req.URL.RawQuery = q.Encode()
	// 		if err != nil {
	// 			log.Println(err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		resp, err := client.Do(req)
	// 		if err != nil {
	// 			log.Println(err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		var webfinger api.WebFingerResponse
	// 		if err = json.NewDecoder(req.Body).Decode(&webfinger); err != nil {
	// 			log.Println(err)
	// 			_ = resp.Body.Close()
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		_ = resp.Body.Close()
	//
	// 		// at this point we successfully looked up a account
	// 		// and we should ask the other server for actor associated with the account
	//
	// 		req, err = http.NewRequest("GET", webfinger.Links[0].Href, nil)
	// 		if err != nil {
	// 			log.Println(err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		resp, err = client.Do(req)
	// 		if err != nil {
	// 			log.Println(err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		var actor Actor
	// 		if err = json.NewDecoder(resp.Body).Decode(&actor); err != nil {
	// 			log.Println(err)
	// 			_ = resp.Body.Close()
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		log.Println(actor)
	// 		// we fetched remote actor
	// 		// we must store it in database and return account in response
	// 		_ = resp.Body.Close()
	// 		row, err := s.ds.Raw().CreateActor(r.Context(), db.CreateActorParams{
	// 			Username:    username, // probably...
	// 			Uri:         actor.ID,
	// 			DisplayName: sql.NullString{actor.PreferredUsername, true}, // probably not...
	// 			Domain:      sql.NullString{domain, true},
	// 			InboxUri:    actor.Inbox,
	// 			OutboxUri:   actor.Outbox,
	// 			Url:         "", // TODO: we should send web profile addr in webfinger
	// 		})
	// 		if err != nil {
	// 			log.Println(err)
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		log.Println(row)
	// 		account := api.Account{
	// 			Acct:        acct,
	// 			DisplayName: row.DisplayName.String,
	// 			Id:          int(row.ID),
	// 			Url:         row.Url,
	// 			Username:    row.Username,
	// 		}
	// 		_ = json.NewEncoder(w).Encode(&account)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 	}
	// 	account := api.Account{
	// 		Acct:        acct,
	// 		DisplayName: actor.DisplayName.String,
	// 		Id:          int(actor.ID),
	// 		Url:         actor.Url,
	// 		Username:    actor.Username,
	// 	}
	// 	_ = json.NewEncoder(w).Encode(&account)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }
}

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id int) {
	// we were requested to create new follow relation
	// we should extract user data from auth token (auth middleware not fully operational yet)
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
	panic("unimplemented")
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
