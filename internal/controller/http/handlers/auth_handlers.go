// internal/controller/http/handlers/user.go
package handlers

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/errors"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/log"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	uc     usecase.UserUsecase
	tokens domain.TokenService
	log    log.Logger
}

func NewUserHandler(uc usecase.UserUsecase, tokens domain.TokenService, logger log.Logger) *UserHandler {
	return &UserHandler{
		uc:     uc,
		tokens: tokens,
		log:    logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err := h.uc.Register(r.Context(), req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		h.log.Errorf(errors.ErrUserRegister, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Errorf(errors.ErrUserLogin, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)

		return
	}

	token, err := h.tokens.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		h.log.Errorf(errors.ErrGenerateToken)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		h.log.Errorf(errors.ErrEncodeJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
