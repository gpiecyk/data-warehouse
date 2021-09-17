package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gpiecyk/data-warehouse/internal/users"
)

func (handler *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(users.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdUser, err := handler.api.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("%s/%d", r.URL.Path, createdUser.ID)
	userBytes, err := json.Marshal(createdUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Location", path)
	w.WriteHeader(http.StatusCreated)
	w.Write(userBytes)
}

func (handler *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(users.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	value := vars["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := handler.api.UpdateUser(r.Context(), user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userBytes, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(userBytes)
}

func (handler *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.api.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	id, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := handler.api.GetUserById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(userBytes)
}

func (handler *Handlers) FindUsers(w http.ResponseWriter, r *http.Request) {
	limitString := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitString)

	users, err := handler.api.FindUsersWithLimit(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersBytes)
}
