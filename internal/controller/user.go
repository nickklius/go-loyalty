package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/auth"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type userRoutes struct {
	u usecase.User
	l zap.Logger
	c config.Config
}

func newUserRoutes(handler *chi.Mux, l *zap.Logger, u usecase.User, c *config.Config) {
	r := &userRoutes{u, *l, *c}

	handler.Post("/api/user/register", r.Register)
	handler.Post("/api/user/login", r.Login)
}

func (ur *userRoutes) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var user entity.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ur.u.CreateUser(r.Context(), user)

	if err != nil {
		statusCode, msg := parseError(err)
		http.Error(w, msg, statusCode)
		return
	}

	u, err := ur.u.CheckPassword(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.CreateToken(u.ID, ur.c.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)

	w.WriteHeader(http.StatusOK)
}

func (ur *userRoutes) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "the body is missing", http.StatusBadRequest)
		return
	}

	var user entity.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := ur.u.CheckPassword(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(u.ID, ur.c.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)

	w.WriteHeader(http.StatusOK)
}
