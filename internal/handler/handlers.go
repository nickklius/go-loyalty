package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-loyalty/config"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/middleware"
	"github.com/nickklius/go-loyalty/internal/usecase"
	"github.com/nickklius/go-loyalty/internal/utils"
)

type Handler struct {
	u usecase.UseCase
	l zap.Logger
	c config.Config
}

func newRoutes(handler *chi.Mux, l *zap.Logger, u *usecase.UseCase, c *config.Config) {
	r := &Handler{*u, *l, *c}

	handler.Post("/api/user/register", r.Register)
	handler.Post("/api/user/login", r.Login)
	handler.Post("/api/user/orders", r.CreateOrder)

	handler.Get("/api/user/orders", r.GetOrders)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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

	err = h.u.CreateUser(r.Context(), user)

	if err != nil {
		statusCode, msg := parseError(err)
		http.Error(w, msg, statusCode)
		return
	}

	u, err := h.u.CheckPassword(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := utils.CreateToken(u.ID, h.c.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "nil body", http.StatusBadRequest)
		return
	}

	var user entity.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := h.u.CheckPassword(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := utils.CreateToken(u.ID, h.c.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token.AccessToken)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "text/plain")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "nil body", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if !utils.ValidLuhnNumber(number) {
		http.Error(w, "wrong order number", http.StatusUnprocessableEntity)
		return
	}

	userID := middleware.GetClaims(r.Context())

	order := entity.Order{
		UserID: userID,
		Number: strconv.Itoa(number),
		Status: "NEW",
	}

	err = h.u.CreateOrder(r.Context(), order)
	if err != nil {
		statusCode, msg := parseError(err)
		http.Error(w, msg, statusCode)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")

	userID := middleware.GetClaims(r.Context())

	orders, err := h.u.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
