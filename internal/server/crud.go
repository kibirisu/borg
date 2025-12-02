package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kibirisu/borg/internal/domain"
)

func create[R domain.Repository[T, Create, Update], T, Create, Update any](
	repo R,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Create
		if err := json.NewDecoder(r.Body).Decode(item); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := repo.Create(r.Context(), item); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getByID[R domain.Repository[T, Create, Update], T, Create, Update any](
	repo R,
	id int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, err := repo.GetByID(r.Context(), int32(id))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(item)
	}
}

func getByUserID[R domain.UserScopedRepository[T, Create, Update], T, Create, Update any](
	repo R,
	id int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.GetByUserID(r.Context(), int32(id))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&items)
	}
}

func getByPostID[R domain.PostScopedRepository[T, Create, Update], T, Create, Update any](
	repo R,
	id int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.GetByPostID(r.Context(), int32(id))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&items)
	}
}

func deleteByID[R domain.Repository[T, Create, Update], T, Create, Update any](
	repo R,
	id int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(id)
		if err := repo.Delete(r.Context(), int32(id)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func update[R domain.Repository[T, Create, Update], T, Create, Update any](
	repo R,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Update
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := repo.Update(r.Context(), item); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getFollowers(repo domain.UserRepository, id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := repo.GetFollowers(r.Context(), int32(id))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&users)
	}
}

func getFollowing(repo domain.UserRepository, id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := repo.GetFollowed(r.Context(), int32(id))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&users)
	}
}

func getAll[R interface {
	GetAll(context.Context) ([]T, error)
}, T any](
	repo R,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.GetAll(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&items)
	}
}
