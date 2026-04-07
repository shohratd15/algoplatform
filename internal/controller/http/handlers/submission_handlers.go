// internal/controller/http/handlers/submission_handlers.go
package handlers

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/errors"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/log"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SubmissionHandler struct {
	usecase usecase.SubmissionUsecase
	val     domain.Validator
	log     log.Logger
}

func NewSubmissionHandler(u usecase.SubmissionUsecase, v domain.Validator, logger log.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		usecase: u,
		val:     v,
		log:     logger,
	}
}

// POST /api/submissions
func (h *SubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProblemID  int64  `json:"problem_id"  validate:"required,min=1"`
		LanguageID int    `json:"language_id" validate:"required,min=1"`
		SourceCode string `json:"source_code" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	if err := h.val.Struct(&req); err != nil {
		h.log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(domain.ClaimsKey).(domain.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sub := &domain.Submission{
		UserID:     claims.UserID,
		ProblemID:  req.ProblemID,
		LanguageID: req.LanguageID,
		SourceCode: req.SourceCode,
		Status:     domain.StatusQueued,
	}

	id, err := h.usecase.Create(r.Context(), sub)
	if err != nil {
		h.log.Errorf(errors.ErrSubmissionCreate, err)
		http.Error(w, errors.ErrSubmissionCreate, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

// GET /api/submissions?id=123
func (h *SubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf("failed to parse id: %v", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		h.log.Errorf(errors.ErrSubmissionNotFound, err)
		http.Error(w, errors.ErrSubmissionNotFound, http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(domain.ClaimsKey).(domain.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "admin" && sub.UserID != claims.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	resp := struct {
		ID         int64     `json:"id"`
		UserID     int64     `json:"user_id"`
		ProblemID  int64     `json:"problem_id"`
		LanguageID int       `json:"language_id"`
		SourceCode string    `json:"source_code"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}{
		ID:         sub.ID,
		UserID:     sub.UserID,
		ProblemID:  sub.ProblemID,
		LanguageID: sub.LanguageID,
		SourceCode: sub.SourceCode,
		Status:     sub.Status,
		CreatedAt:  sub.CreatedAt,
		UpdatedAt:  sub.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}
