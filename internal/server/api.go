package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/ap"
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

	log.Printf(
		"lookup: remote account %s@%s not cached, performing WebFinger lookup",
		handle.Username,
		handle.Domain,
	)
	util.WriteError(w, http.StatusInternalServerError, "unimplemented")
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

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id int) {
	container, ok := r.Context().Value(TokenContextKey).(*tokenContainer)

	if !ok || container == nil || container.id == nil {
		util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	currentUserID := *container.id
	if currentUserID == id {
		http.Error(w, "Tried to follow oneself", http.StatusBadRequest)
		return
	}
	follower, err := s.service.App.GetAccountByID(r.Context(), currentUserID)
	followee, err := s.service.App.GetAccountByID(r.Context(), id)
	follow, err := s.service.App.FollowAccount(r.Context(), currentUserID, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// APfollow := mapper.DBToFollow(follow, &follower, &followee)
	followActivity := ap.NewActivity(nil)
	actor := ap.NewActor(nil)
	actor.SetLink(follower.Uri)
	object := ap.NewActor(nil)
	object.SetLink(followee.Uri)
	followActivity.SetObject(ap.Activity[any]{
		ID:     follow.Uri,
		Type:   "Follow",
		Actor:  actor,
		Object: object.(ap.Objecter[any]),
	})
	log.Println(followee.InboxUri)
	if follower.Domain != followee.Domain {
		util.DeliverToEndpoint(followee.InboxUri, followActivity.GetRaw())
	}
	util.WriteJSON(w, http.StatusCreated, nil)
}

// DeleteApiUsersId implements api.ServerInterface.
func (s *Server) DeleteApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiUsersId implements api.ServerInterface.
func (s *Server) GetApiUsersId(w http.ResponseWriter, r *http.Request, id int) {
	user, err := s.service.App.GetAccountByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSON(w, http.StatusOK, *mapper.AccountToAPI(&user))
}

// PostApiUsers implements api.ServerInterface.
func (s *Server) PostApiUsers(w http.ResponseWriter, r *http.Request) {
	var newUser api.NewUser
	if err := util.ReadJSON(r, &newUser); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("api users: creating user username=%s", newUser.Username)

	// Use existing Register method which creates both account and user
	authForm := api.AuthForm{
		Username: newUser.Username,
		Password: newUser.Password,
	}
	if err := s.service.App.Register(r.Context(), authForm); err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the created account to return user data
	account, err := s.service.App.GetLocalAccount(r.Context(), newUser.Username)
	if err != nil {
		log.Println(err)
		util.WriteError(w, http.StatusInternalServerError, "Failed to retrieve created user")
		return
	}

	// Get followers and following counts
	followers, err := s.service.App.GetAccountFollowers(r.Context(), int(account.ID))
	if err != nil {
		log.Println(err)
		// Continue even if counts fail, use 0 as default
		followers = []db.Account{}
	}

	following, err := s.service.App.GetAccountFollowing(r.Context(), int(account.ID))
	if err != nil {
		log.Println(err)
		// Continue even if counts fail, use 0 as default
		following = []db.Account{}
	}

	user := mapper.AccountToUserAPI(account, len(followers), len(following))

	log.Printf("api users: user %s created successfully with id=%d", newUser.Username, user.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	util.WriteJSON(w, http.StatusCreated, *user)
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

	commenter, err := s.service.App.GetAccountByID(r.Context(), currentUserID)
	parentPost, err := s.service.App.GetPostByID(r.Context(), id)
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

	parentAuthor, _ := s.service.App.GetAccountByID(r.Context(), int(parentPost.AccountID))

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
	likes, err := s.service.App.GetPostLikes(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	apiLikes := make([]api.Like, 0, len(likes))

	for _, like := range likes {
		converted := mapper.LikeToAPI(&like)
		apiLikes = append(apiLikes, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiLikes)
}

// PostApiPostsIdLikes implements api.ServerInterface.
func (s *Server) PostApiPostsIdLikes(w http.ResponseWriter, r *http.Request, id int) {
	var newLike api.NewLike
	if err := util.ReadJSON(r, &newLike); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	currentUserID := newLike.UserID

	liker, err := s.service.App.GetAccountByID(r.Context(), currentUserID)
	post, err := s.service.App.GetPostByID(r.Context(), id)
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

	author, err := s.service.App.GetAccountByID(r.Context(), int(post.AccountID))
	if err == nil && liker.Domain != author.Domain {
		util.DeliverToEndpoint(author.InboxUri, APLike)
	}

	util.WriteJSON(w, http.StatusCreated, nil)
}

// GetApiPostsIdShares implements api.ServerInterface.
func (s *Server) GetApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
	shares, err := s.service.App.GetPostShares(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	apiShares := make([]api.Post, 0, len(shares))

	for _, share := range shares {
		converted := mapper.PostToAPI(&share)
		apiShares = append(apiShares, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiShares)
}

// PostApiPostsIdShares implements api.ServerInterface.
func (s *Server) PostApiPostsIdShares(w http.ResponseWriter, r *http.Request, id int) {
	var newShare api.NewShare
	if err := util.ReadJSON(r, &newShare); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	currentUserID := newShare.UserID

	sharer, err := s.service.App.GetAccountByID(r.Context(), currentUserID)
	post, err := s.service.App.GetPostByID(r.Context(), id)
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

	author, err := s.service.App.GetAccountByID(r.Context(), int(post.AccountID))
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
	container, ok := r.Context().Value(TokenContextKey).(*tokenContainer)
	if !ok || container == nil || container.id == nil {
		util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	currentUserID := *container.id
	poster, err := s.service.App.GetAccountByID(r.Context(), currentUserID)
	var newPost api.NewPost
	if err := util.ReadJSON(r, &newPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	newDBPost := mapper.NewPostToDB(&newPost, true)
	status, err := s.service.App.AddNote(r.Context(), *newDBPost)
	if err != nil {
		http.Error(w, "Cannot add note"+err.Error(), http.StatusInternalServerError)
		return
	}

	s.service.App.DeliverToFollowers(w, r, newPost.UserID, func(recipientURI string) any {
		create := ap.NewCreateActivity(nil)
		actor := ap.NewActor(nil)
		actor.SetLink(poster.Uri)
		note := ap.NewNote(nil)
		replies := ap.NewNoteCollection(nil)
		replies.SetLink("YO MAMA TODO")
		note.SetObject(ap.Note{
			ID:           status.Uri,
			Type:         "Note",
			Content:      status.Content,
			InReplyTo:    ap.NewNote(nil),
			Published:    status.CreatedAt,
			AttributedTo: ap.NewActor(nil),
			To:           []string{recipientURI},
			CC:           []string{recipientURI},
			Replies:      replies,
		})
		create.SetObject(ap.Activity[ap.Note]{
			ID:     "TODO",
			Type:   "Create",
			Actor:  actor,
			Object: note,
		})
		return create.GetRaw()
	})
	util.WriteJSON(w, http.StatusCreated, nil)
}

// PutApiPostsId implements api.ServerInterface.
func (s *Server) PutApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	panic("unimplemented")
}

// GetApiUsersIdPosts implements api.ServerInterface.
func (s *Server) GetApiUsersIdPosts(w http.ResponseWriter, r *http.Request, id int) {
	posts, err := s.service.App.GetPostByAccountID(r.Context(), id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	apiLikes := make([]api.Post, 0, len(posts))

	for _, info := range posts {
		converted := mapper.PostToAPIWithMetadata(&info.Status,
			&info.Account,
			int(info.LikeCount),
			int(info.ShareCount),
			int(info.CommentCount))
		apiLikes = append(apiLikes, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiLikes)
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
	followers, err := s.service.App.GetAccountFollowers(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		return
	}

	apiFollowers := make([]api.Account, 0, len(followers))
	for _, follower := range followers {
		apiFollowers = append(apiFollowers, *mapper.AccountToAPI(&follower))
	}

	util.WriteJSON(w, http.StatusOK, apiFollowers)
}

// GetApiUsersIdFollowing implements api.ServerInterface.
func (s *Server) GetApiUsersIdFollowing(w http.ResponseWriter, r *http.Request, id int) {
	following, err := s.service.App.GetAccountFollowing(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch following", http.StatusInternalServerError)
		return
	}

	apiFollowers := make([]api.Account, 0, len(following))
	for _, follower := range following {
		apiFollowers = append(apiFollowers, *mapper.AccountToAPI(&follower))
	}

	util.WriteJSON(w, http.StatusOK, apiFollowers)
}

// GetApiPosts implements api.ServerInterface.
func (s *Server) GetApiPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.service.App.GetLocalPosts(r.Context())
	if err != nil {
		http.Error(w, "Database error "+err.Error(), http.StatusNotFound)
		return
	}

	apiLikes := make([]api.Post, 0, len(posts)) // i see smth bad right here...

	for _, info := range posts {
		converted := mapper.PostToAPIWithMetadata(&info.Status,
			&info.Account,
			int(info.LikeCount),
			int(info.ShareCount),
			int(info.CommentCount))
		apiLikes = append(apiLikes, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiLikes)
}

// GetApiUsersIdFavourites implements api.ServerInterface.
func (s *Server) GetApiUsersIdFavourites(w http.ResponseWriter, r *http.Request, id int) {
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
func (s *Server) GetApiUsersIdReblogged(w http.ResponseWriter, r *http.Request, id int) {
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
func (s *Server) GetApiUsersIdTimeline(w http.ResponseWriter, r *http.Request, id int) {
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
