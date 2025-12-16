package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

// Since there gonna be a lot of API requests triggering our server to perform HTTP request to another server,
// we build worker that will perform those request and return to the calling handler, ideally.
// The worker shall work separately from the server thread performing tasks enqueued by the handlers

// PostAuthLogin implements api.ServerInterface.
func (s *Server) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var form api.AuthForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	auth, err := s.ds.Raw().AuthData(r.Context(), form.Username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(form.Password)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  auth.ID,
		"iss":  "http://" + s.conf.ListenHost,
		"name": form.Username,
	})
	token, err := jwt.SignedString([]byte("changeme"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer: "+token)
	w.WriteHeader(http.StatusOK)
}

// PostAuthRegister implements api.ServerInterface.
func (s *Server) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var form api.AuthForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// in future we need validate that provided username is sanitized and without "@" character
	uri := "http://" + s.conf.ListenHost + "/users/" + form.Username
	actor, err := s.ds.Raw().CreateActor(r.Context(), db.CreateActorParams{
		Username:    form.Username,
		Uri:         uri,
		DisplayName: sql.NullString{}, // hassle to maintain that, gonna abandon display name
		Domain:      sql.NullString{},
		InboxUri:    uri + "/inbox",
		OutboxUri:   uri + "/outbox",
		Url:         "http://" + s.conf.ListenHost + "/profiles/" + form.Username, // maybe we make profile page on /@<username> ???
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = s.ds.Raw().CreateUser(r.Context(), db.CreateUserParams{
		AccountID:    actor.ID,
		PasswordHash: "",
	}); err != nil {
		// bad, very BAD thing happened, probably good idea is to delete actor from db
		// or use one query but i don't remember postgres this much if that would even work
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetApiAccountsLookup implements api.ServerInterface.
// DEMO
func (s *Server) GetApiAccountsLookup(
	w http.ResponseWriter,
	r *http.Request,
	params api.GetApiAccountsLookupParams,
) {
	// we must check if account is local or from other instance
	// if from other instance we do webfinger lookup
	acct := params.Acct
	log.Println(acct)
	arr := strings.Split(acct, "@")
	username := arr[0]
	domain := arr[1]

	if domain == s.conf.ListenHost {
		actor, err := s.ds.Raw().GetActor(r.Context(), username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		account := api.Account{
			Acct:        acct,
			DisplayName: actor.DisplayName.String,
			Id:          int(actor.ID),
			Url:         actor.Url,
			Username:    actor.Username,
		}
		json.NewEncoder(w).Encode(&account)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else if domain != "" {
		actor, err := s.ds.Raw().GetAccount(r.Context(), db.GetAccountParams{username, sql.NullString{domain, true}})
		if err != nil {
			// we should do webfinger lookup at this point

			client := http.Client{Timeout: 2 * time.Second}
			req, err := http.NewRequest("GET", "http://"+domain+"/.well-known/webfinger", nil)
			q := req.URL.Query()
			q.Set("resource", acct)
			req.URL.RawQuery = q.Encode()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var webfinger api.WebFingerResponse
			if err = json.NewDecoder(req.Body).Decode(&webfinger); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_ = resp.Body.Close()

			// at this point we successfully looked up a account
			// and we should ask the other server for actor associated with the account

			req, err = http.NewRequest("GET", webfinger.Links[0].Href, nil)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp, err = client.Do(req)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var actor Actor
			if err = json.NewDecoder(resp.Body).Decode(&actor); err != nil {
				log.Println(err)
				_ = resp.Body.Close()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(actor)
			// we fetched remote actor
			// we must store it in database and return account in response
			_ = resp.Body.Close()
			row, err := s.ds.Raw().CreateActor(r.Context(), db.CreateActorParams{
				Username:    username, // probably...
				Uri:         actor.ID,
				DisplayName: sql.NullString{actor.PreferredUsername, true}, // probably not...
				Domain:      sql.NullString{domain, true},
				InboxUri:    actor.Inbox,
				OutboxUri:   actor.Outbox,
				Url:         "", // TODO: we should send web profile addr in webfinger
			})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(row)
			account := api.Account{
				Acct:        acct,
				DisplayName: row.DisplayName.String,
				Id:          int(row.ID),
				Url:         row.Url,
				Username:    row.Username,
			}
			_ = json.NewEncoder(w).Encode(&account)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
		account := api.Account{
			Acct:        acct,
			DisplayName: actor.DisplayName.String,
			Id:          int(actor.ID),
			Url:         actor.Url,
			Username:    actor.Username,
		}
		_ = json.NewEncoder(w).Encode(&account)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// PostApiAccountsIdFollow implements api.ServerInterface.
func (s *Server) PostApiAccountsIdFollow(w http.ResponseWriter, r *http.Request, id int) {
	// we were requested to create new follow relation
	// we should extract user data from auth token (auth middleware not fully operational yet)
}
