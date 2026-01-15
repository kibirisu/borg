package server

import (
	"database/sql"
	"encoding/json/v2"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
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
	// we must check if account is local or from other instance
	// if from other instance we do webfinger lookup
	acct := params.Acct
	handle, err := util.ParseHandle(acct, s.conf.ListenHost, s.conf.ListenPort)
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
	} else {
		log.Printf("lookup: remote handle %s detected, checking local cache", acct)
		account, err := s.service.App.GetAccount(r.Context(), db.GetAccountParams{
			Username: handle.Username,
			Domain: sql.NullString{
				String: handle.Domain,
				Valid:  true,
			},
		})
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

		client := http.Client{Timeout: 5 * time.Second}
		actorURL := "http://" + handle.Domain + "/user/" + handle.Username

		reqActor, err := http.NewRequestWithContext(r.Context(), http.MethodGet, actorURL, nil)
		if err != nil {
			log.Println(err)
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("lookup: fetching remote actor %s", actorURL)
		actorResp, err := client.Do(reqActor)
		if err != nil {
			log.Println(err)
			util.WriteError(w, http.StatusBadGateway, err.Error())
			return
		}
		defer func() {
			err = actorResp.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		if actorResp.StatusCode != http.StatusOK {
			util.WriteError(w, http.StatusBadGateway, "remote actor fetch failed")
			return
		}

		var object domain.ObjectOrLink
		if err := json.UnmarshalRead(actorResp.Body, &object); err != nil {
			log.Println(err)
			util.WriteError(w, http.StatusBadGateway, err.Error())
			return
		}
		actor := ap.NewActor(&object)
		actorData := actor.GetObject()

		params := db.CreateActorParams{
			Username: actorData.PreferredUsername,
			Uri:      actorData.ID,
			Domain: sql.NullString{
				String: handle.Domain,
				Valid:  true,
			},
			DisplayName:  sql.NullString{},
			InboxUri:     actorData.Inbox,
			OutboxUri:    actorData.Outbox,
			FollowersUri: actorData.Followers,
			FollowingUri: actorData.Following,
			Url:          actorData.ID,
		}

		log.Printf(
			"lookup: creating remote account %s@%s from actor %s",
			actorData.PreferredUsername, handle.Domain, actorData.ID,
		)
		account, err = s.service.App.AddRemoteAccount(r.Context(), &params)
		if err != nil {
			log.Println(err)
			// If already exists, fetch the stored record.
			existing, getErr := s.service.App.GetAccount(r.Context(), db.GetAccountParams{
				Username: actorData.PreferredUsername,
				Domain: sql.NullString{
					String: handle.Domain,
					Valid:  true,
				},
			})
			if getErr != nil {
				util.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(existing))
			return
		}

		util.WriteJSON(w, http.StatusOK, mapper.AccountToAPI(account))
	}
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
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	followee, err := s.service.App.GetAccountByID(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	follow, err := s.service.App.FollowAccount(r.Context(), currentUserID, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	followActivity := ap.NewFollowActivity(nil)
	actor := ap.NewActor(nil)
	actor.SetLink(follower.Uri)
	actoree := ap.NewActor(nil)
	actoree.SetLink(followee.Uri)
	followActivity.SetObject(ap.Activity[ap.Actor]{
		ID:     follow.Uri,
		Type:   "Follow",
		Actor:  actor,
		Object: actoree,
	})
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
	account, err := s.service.App.GetAccountByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get followers and following counts
	followers, err := s.service.App.GetAccountFollowers(r.Context(), id)
	if err != nil {
		log.Println(err)
		followers = []db.Account{}
	}

	following, err := s.service.App.GetAccountFollowing(r.Context(), id)
	if err != nil {
		log.Println(err)
		following = []db.Account{}
	}

	user := mapper.AccountToUserAPI(&account, len(followers), len(following))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSON(w, http.StatusOK, *user)
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
	authForm := api.AuthForm(newUser)
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
	container, ok := r.Context().Value(TokenContextKey).(*tokenContainer)
	if !ok || container == nil || container.id == nil {
		util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	currentUserID := *container.id

	if id != currentUserID {
		util.WriteError(w, http.StatusForbidden, "Forbidden: You can only update your own account")
		return
	}

	_, err := s.service.App.GetAccountByID(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	var update api.UpdateUser
	if err := util.ReadJSON(r, &update); err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var bio *string
	if update.Bio != nil {
		bio = update.Bio
	}

	updatedAccount, err := s.service.App.UpdateAccount(r.Context(), id, bio)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	followers, err := s.service.App.GetAccountFollowers(r.Context(), id)
	if err != nil {
		log.Println(err)
		followers = []db.Account{}
	}

	following, err := s.service.App.GetAccountFollowing(r.Context(), id)
	if err != nil {
		log.Println(err)
		following = []db.Account{}
	}

	user := mapper.AccountToUserAPI(&updatedAccount, len(followers), len(following))
	util.WriteJSON(w, http.StatusOK, *user)
}

// DeleteApiPostsId implements api.ServerInterface.
func (s *Server) DeleteApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	container, ok := r.Context().Value(TokenContextKey).(*tokenContainer)
	if !ok || container == nil || container.id == nil {
		util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	currentUserID := *container.id

	post, err := s.service.App.GetPostByID(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if int(post.AccountID) != currentUserID {
		util.WriteError(w, http.StatusForbidden, "Forbidden: You can only delete your own posts")
		return
	}

	err = s.service.App.DeletePost(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetApiPostsId implements api.ServerInterface.
func (s *Server) GetApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	info, err := s.service.App.GetPostByIDWithMetadata(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	util.WriteJSON(w, http.StatusOK, *mapper.PostToAPIWithMetadata(&info.Status,
		&info.Account,
		int(info.LikeCount),
		int(info.ShareCount),
		int(info.CommentCount)))
}

// GetApiPostsIdComments implements api.ServerInterface.
func (s *Server) GetApiPostsIdComments(w http.ResponseWriter, r *http.Request, id int) {
	comments, err := s.service.App.GetCommentsByPostID(r.Context(), id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	apiComments := make([]api.Comment, 0, len(comments))
	for _, comment := range comments {
		converted := mapper.StatusToComment(&comment)
		apiComments = append(apiComments, *converted)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSON(w, http.StatusOK, apiComments)
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
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
	}
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
	activity := mapper.StatusToCreateActivity(status, commenter, parentPost)

	s.service.App.DeliverToFollowers(w, r, currentUserID, func(recipientURI string) any {
		return activity.GetRaw()
	})

	// deliver to original post author
	if commenter.Domain != parentAuthor.Domain {
		util.DeliverToEndpoint(parentAuthor.InboxUri, activity.GetRaw())
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
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	post, err := s.service.App.GetPostByID(r.Context(), id)
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	like, err := s.service.App.AddFavourite(r.Context(), currentUserID, id)
	if err != nil {
		if isUniqueViolation(err) {
			util.WriteError(w, http.StatusConflict, "Post already liked")
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	actor := ap.NewActor(nil)
	actor.SetLink(liker.Uri)
	note := ap.NewNote(nil)
	note.SetLink(post.Uri)
	activity := ap.NewLikeActivity(nil)
	activity.SetObject(ap.Activity[ap.Note]{
		ID:     like.Uri,
		Type:   "Like",
		Actor:  actor,
		Object: note,
	})

	author, err := s.service.App.GetAccountByID(r.Context(), int(post.AccountID))
	if err == nil && liker.Domain != author.Domain {
		util.DeliverToEndpoint(author.InboxUri, activity)
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
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
	}
	post, err := s.service.App.GetPostByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	status := mapper.NewShareToDB(&newShare)

	share, err := s.service.App.AddNote(r.Context(), *status)
	if err != nil {
		if isUniqueViolation(err) {
			util.WriteError(w, http.StatusConflict, "Post already shared")
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	activity := mapper.StatusToCreateActivity(share, sharer, nil)

	// send activity to post author
	author, err := s.service.App.GetAccountByID(r.Context(), int(post.AccountID))
	if err == nil && sharer.Domain != author.Domain {
		util.DeliverToEndpoint(author.InboxUri, activity.GetRaw())
	}

	// send activity to my followers
	s.service.App.DeliverToFollowers(w, r, currentUserID, func(recipientURI string) any {
		return activity.GetRaw()
	})

	util.WriteJSON(w, http.StatusCreated, nil)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func applyReshareInfo(
	post *api.Post,
	resharedBy sql.NullString,
	resharedByID sql.NullInt32,
) {
	if resharedBy.Valid {
		post.ResharedBy = &resharedBy.String
	}
	if resharedByID.Valid {
		value := int(resharedByID.Int32)
		post.ResharedById = &value
	}
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
	if err != nil {
		util.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
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

	activity := mapper.StatusToCreateActivity(status, poster, nil)

	s.service.App.DeliverToFollowers(w, r, newPost.UserID, func(recipientURI string) any {
		return activity.GetRaw()
	})
	util.WriteJSON(w, http.StatusCreated, nil)
}

// PutApiPostsId implements api.ServerInterface.
func (s *Server) PutApiPostsId(w http.ResponseWriter, r *http.Request, id int) {
	// 1. Authorization - check if user is authenticated
	container, ok := r.Context().Value(TokenContextKey).(*tokenContainer)
	if !ok || container == nil || container.id == nil {
		util.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	currentUserID := *container.id

	// 2. Check if post exists and get current post data
	post, err := s.service.App.GetPostByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// 3. Check ownership - only owner can update their post
	if int(post.AccountID) != currentUserID {
		util.WriteError(w, http.StatusForbidden, "Forbidden: You can only update your own posts")
		return
	}

	// 4. Read request body
	var update api.UpdatePost
	if err := util.ReadJSON(r, &update); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 5. Validate content - check if content is provided
	if update.Content == nil || *update.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	// 6. Update the post
	_, err = s.service.App.UpdatePost(r.Context(), id, *update.Content)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 7. Get updated post with metadata for response
	info, err := s.service.App.GetPostByIDWithMetadata(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 8. Return updated post with metadata
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	util.WriteJSON(w, http.StatusOK, *mapper.PostToAPIWithMetadata(
		&info.Status,
		&info.Account,
		int(info.LikeCount),
		int(info.ShareCount),
		int(info.CommentCount)))
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
		applyReshareInfo(converted, info.ResharedBy, info.ResharedByID)
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
		applyReshareInfo(converted, info.ResharedBy, info.ResharedByID)
		apiLikes = append(apiLikes, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiLikes)
}

// GetApiUsersIdFavourites implements api.ServerInterface.
func (s *Server) GetApiUsersIdFavourites(w http.ResponseWriter, r *http.Request, id int) {
	posts, err := s.service.App.GetLikedPostsByAccountID(r.Context(), id)
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
	posts, err := s.service.App.GetSharedPostsByAccountID(r.Context(), id)
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
		applyReshareInfo(converted, info.ResharedBy, info.ResharedByID)
		apiPosts = append(apiPosts, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiPosts)
}

// GetApiUsersIdTimeline implements api.ServerInterface.
func (s *Server) GetApiUsersIdTimeline(w http.ResponseWriter, r *http.Request, id int) {
	posts, err := s.service.App.GetTimelinePostsByAccountID(r.Context(), id)
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
		applyReshareInfo(converted, info.ResharedBy, info.ResharedByID)
		apiPosts = append(apiPosts, *converted)
	}

	util.WriteJSON(w, http.StatusOK, apiPosts)
}
