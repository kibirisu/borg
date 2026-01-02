package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/domain"
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

	if err := s.service.App.Register(r.Context(), form); err != nil {
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
	// // we must check if account is local or from other instance
	// // if from other instance we do webfinger lookup
	acct := params.Acct
	arr := strings.Split(acct, "@")
	username := arr[0]
	addr := arr[1]
	//
	if addr == s.conf.ListenHost {
		account, err := s.service.App.GetLocalAccount(r.Context(), username)
		if err != nil {
			log.Println(err)
			util.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
	} else if addr != "" {
		// actor, err := s.ds.Raw().GetAccount(r.Context(), db.GetAccountParams{username, sql.NullString{domain, true}})
		account, err := s.service.App.GetLocalAccount(r.Context(), username)
		if err != nil {
			// we should do webfinger lookup at this point
			// code bellow will be move to worker

			client := http.Client{Timeout: 2 * time.Second}
			req, err := http.NewRequest("GET", "http://"+addr+"/.well-known/webfinger", nil)
			q := req.URL.Query()
			q.Set("resource", acct)
			req.URL.RawQuery = q.Encode()
			if err != nil {
				log.Println(err)
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			var webfinger api.WebFingerResponse
			if err = util.ReadJSON(r, &webfinger); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			_ = resp.Body.Close()

			// at this point we successfully looked up a account
			// and we should ask the other server for actor associated with the account

			req, err = http.NewRequest("GET", webfinger.Links[0].Href, nil)
			if err != nil {
				log.Println(err)
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			resp, err = client.Do(req)
			if err != nil {
				log.Println(err)
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			var actor domain.Actor
			if err = util.ReadJSON(r, &actor); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			log.Println(actor)
			// we fetched remote actor
			// we must store it in database and return account in response
			_ = resp.Body.Close()
			row, err := s.service.Federation.CreateActor(r.Context(), *mapper.ActorToDB(&actor, addr))
			if err != nil {
				log.Println(err)
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			log.Println(row)
			util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(row))
		}
		util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
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
    var comment api.NewComment
    if err := util.ReadJSON(r, &comment); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return 
    }
	currentUserID := comment.UserID

    commenter, err := s.service.App.GetAccountById(r.Context(), currentUserID)
    parentPost, err := s.service.App.GetPostById(r.Context(), id)
    if err != nil {
        http.Error(w, "Parent post not found", http.StatusNotFound)
        return
    }

    dbComment := mapper.NewCommentToDB(&comment) 
    status, err := s.service.App.AddNote(r.Context(), *dbComment)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    parentAuthor, _ := s.service.App.GetAccountById(r.Context(), int(parentPost.AccountID))
    
    s.service.App.DeliverToFollowers(w, r, currentUserID, func(recipientURI string) any {
        return mapper.PostToCreateNote(&status, &commenter, parentAuthor.FollowersUri)
    })

    if commenter.Domain != parentAuthor.Domain {
        APComment := mapper.PostToCreateNote(&status, &commenter, parentAuthor.Uri)
        util.DeliverToEndpoint(parentAuthor.InboxUri, APComment)
    }

    util.WriteJSON(w, http.StatusCreated, nil)
}

// GetApiPostsIdLikes implements api.ServerInterface.
func (s *Server) GetApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPostsIdLikes implements api.ServerInterface.
func (s *Server) PostApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id int) {
    var newLike api.NewLike
    if err := util.ReadJSON(r, &newLike); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return 
    }
	currentUserID := newLike.UserID

    liker, err := s.service.App.GetAccountById(r.Context(), currentUserID)
    post, err := s.service.App.GetPostById(r.Context(), id)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    like, err := s.service.App.AddFavourite(r.Context(), currentUserID, id)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    APLike := mapper.DBToFavourite(&like, &liker, post)
    
    author, err := s.service.App.GetAccountById(r.Context(), int(post.AccountID))
    if err == nil && liker.Domain != author.Domain {
        util.DeliverToEndpoint(author.InboxUri, APLike)
    }

    util.WriteJSON(w, http.StatusCreated, nil)
}

// GetApiPostsIdShares implements api.ServerInterface.
func (s *Server) GetApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// PostApiPostsIdShares implements api.ServerInterface.
func (s *Server) PostApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
    var newShare api.NewShare // Assuming this exists in your generated API code
    if err := util.ReadJSON(r, &newShare); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return 
    }
    currentUserID := newShare.UserID

    sharer, err := s.service.App.GetAccountById(r.Context(), currentUserID)
    post, err := s.service.App.GetPostById(r.Context(), id)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
    status := mapper.NewShareToDB(&newShare)

    share, err := s.service.App.AddNote(r.Context(), *status)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }


    author, err := s.service.App.GetAccountById(r.Context(), int(post.AccountID))
    if err == nil && sharer.Domain != author.Domain {
	APAnnounce := mapper.PostToCreateNote(&share, &sharer, author.Uri)
        util.DeliverToEndpoint(author.InboxUri, APAnnounce)
    }

    s.service.App.DeliverToFollowers(w, r, currentUserID, func(recipientURI string) any {
	APAnnounce := mapper.PostToCreateNote(&share, &sharer, author.FollowersUri)
        return APAnnounce
    })

    util.WriteJSON(w, http.StatusCreated, nil)
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
        return mapper.PostToCreateNote(&status, &poster, poster.FollowersUri)
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
