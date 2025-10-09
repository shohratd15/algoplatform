// internal/controller/http/handlers/submission.go
package handlers

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/errors"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SubmissionHandler struct {
	usecase usecase.SubmissionUsecase
	log     log.Logger
}

func NewSubmissionHandler(u usecase.SubmissionUsecase, logger log.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		usecase: u,
		log:     logger,
	}
}

// POST /submissions
func (h *SubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     int64  `json:"user_id"`
		ProblemID  int64  `json:"problem_id"`
		LanguageID int    `json:"language_id"`
		SourceCode string `json:"source_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidRequestBody, http.StatusBadRequest)

		return
	}

	sub := &domain.Submission{
		UserID:     req.UserID,
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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{"id": id}); err != nil {
		h.log.Errorf(errors.ErrEncodeJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /submissions?id=123
func (h *SubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	sub, err := h.usecase.Get(r.Context(), parseID(idStr))
	if err != nil {
		h.log.Errorf(errors.ErrSubmissionNotFound, err)
		http.Error(w, errors.ErrSubmissionNotFound, http.StatusNotFound)
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

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Errorf(errors.ErrEncodeJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseID(s string) int64 {
	var id int64
	_, _ = fmt.Sscan(s, &id)

	return id
}
