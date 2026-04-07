// internal/controller/http/handlers/problem_handlers.go
package handlers

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/errors"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/log"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProblemHandler struct {
	usecase usecase.ProblemUsecase
	val     domain.Validator
	log     log.Logger
}

func NewProblemHandler(u usecase.ProblemUsecase, v domain.Validator, logger log.Logger) *ProblemHandler {
	return &ProblemHandler{
		usecase: u,
		val:     v,
		log:     logger,
	}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UpsertProblemRequest
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

	p := &domain.Problem{
		Slug:       req.Slug,
		Difficulty: req.Difficulty,
	}

	if err := h.usecase.Create(r.Context(), p, toDomainStatements(req.Statements), toDomainTests(req.Tests)); err != nil {
		h.log.Errorf(errors.ErrCreateProblem, err)
		http.Error(w, "failed to create problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf(errors.ErrParseIntID, err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpsertProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err = h.val.Struct(&req); err != nil {
		h.log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := &domain.Problem{
		Slug:       req.Slug,
		Difficulty: req.Difficulty,
	}

	if err = h.usecase.Update(r.Context(), id, p, toDomainStatements(req.Statements), toDomainTests(req.Tests)); err != nil {
		h.log.Errorf("failed to update problem: %v", err)
		http.Error(w, "failed to update problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request) {
	type respItem struct {
		ID         int64     `json:"id"`
		Slug       string    `json:"slug"`
		Difficulty string    `json:"difficulty"`
		CreatedAt  time.Time `json:"created_at"`
	}

	problems, err := h.usecase.List(r.Context())
	if err != nil {
		h.log.Errorf("failed to list problems: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := make([]respItem, 0, len(problems))
	for _, p := range problems {
		resp = append(resp, respItem{
			ID:         p.ID,
			Slug:       p.Slug,
			Difficulty: p.Difficulty,
			CreatedAt:  p.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProblemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf(errors.ErrParseIntID, err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, stmts, tests, err := h.usecase.GetById(r.Context(), id)
	if err != nil {
		h.log.Errorf(errors.ErrGetBySlug, err)
		http.Error(w, "problem not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"problem":    p,
		"statements": stmts,
		"tests":      tests,
	})
}

func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf(errors.ErrParseIntID, err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		h.log.Errorf(errors.ErrDeleteProblem, err)
		http.Error(w, "failed to delete problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---- DTOs ----

type ProblemStatementDTO struct {
	Language  string `json:"language"  validate:"required"`
	Title     string `json:"title"     validate:"required"`
	Statement string `json:"statement" validate:"required"`
}

type ProblemTestDTO struct {
	ID             int64  `json:"id"`
	InputData      string `json:"input_data"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       bool   `json:"is_sample"`
}

type UpsertProblemRequest struct {
	Slug       string                `json:"slug"       validate:"required"`
	Difficulty string                `json:"difficulty" validate:"required"`
	Statements []ProblemStatementDTO `json:"statements" validate:"required,min=1,dive"`
	Tests      []ProblemTestDTO      `json:"tests"      validate:"required,min=1,dive"`
}

func toDomainStatements(dtos []ProblemStatementDTO) []domain.ProblemStatement {
	out := make([]domain.ProblemStatement, 0, len(dtos))
	for _, dto := range dtos {
		out = append(out, domain.ProblemStatement{
			Language:  normalizeStatementLanguage(dto.Language),
			Title:     dto.Title,
			Statement: dto.Statement,
		})
	}
	return out
}

func toDomainTests(dtos []ProblemTestDTO) []domain.ProblemTest {
	out := make([]domain.ProblemTest, 0, len(dtos))
	for _, dto := range dtos {
		out = append(out, domain.ProblemTest{
			ID:             dto.ID,
			InputData:      dto.InputData,
			ExpectedOutput: dto.ExpectedOutput,
			IsSample:       dto.IsSample,
		})
	}
	return out
}

func normalizeStatementLanguage(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "en", "eng":
		return "eng"
	case "ru", "rus":
		return "rus"
	case "tm", "tkm":
		return "tkm"
	default:
		return strings.ToLower(strings.TrimSpace(lang))
	}
}
