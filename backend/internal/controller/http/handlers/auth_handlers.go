// internal/controller/http/handlers/auth_handlers.go
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
	val    domain.Validator
	log    log.Logger
}

func NewUserHandler(uc usecase.UserUsecase, tokens domain.TokenService, v domain.Validator, logger log.Logger) *UserHandler {
	return &UserHandler{
		uc:     uc,
		tokens: tokens,
		val:    v,
		log:    logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.val.Struct(&req); err != nil {
		h.log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.uc.Register(r.Context(), req.Username, req.Email, req.Password); err != nil {
		h.log.Errorf(errors.ErrUserRegister, err)
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.val.Struct(&req); err != nil {
		h.log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Errorf(errors.ErrUserLogin, err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokens.GenerateAccess(user.ID, user.Email, user.Role)
	if err != nil {
		h.log.Errorf(errors.ErrGenerateToken)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.tokens.GenerateRefresh(user.ID, user.Email, user.Role)
	if err != nil {
		h.log.Errorf(errors.ErrGenerateToken)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"role":          user.Role,
	})
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.val.Struct(&req); err != nil {
		h.log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := h.tokens.ParseRefresh(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.tokens.GenerateAccess(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		h.log.Errorf(errors.ErrGenerateToken)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.tokens.GenerateRefresh(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		h.log.Errorf(errors.ErrGenerateToken)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"role":          claims.Role,
	})
}
