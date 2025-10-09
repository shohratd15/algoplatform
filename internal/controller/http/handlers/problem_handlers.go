// internal/controller/http/handlers/problem.go
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

type ProblemHandler struct {
	usecase usecase.ProblemUsecase
	log     log.Logger
}

func NewProblemHandler(u usecase.ProblemUsecase, logger log.Logger) *ProblemHandler {
	return &ProblemHandler{
		usecase: u,
		log:     logger,
	}
}

func (h *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf(errors.ErrInvalidRequestBody, err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	p := &domain.Problem{
		Slug:       req.Slug,
		Difficulty: req.Difficulty,
	}

	statements := toDomainStatements(req.Statements)
	tests := toDomainTests(req.Tests)

	if err := h.usecase.Create(r.Context(), p, statements, tests); err != nil {
		h.log.Errorf(errors.ErrCreateProblem, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProblemHandler) List(w http.ResponseWriter, r *http.Request) {
	type respBody struct {
		ID         int64     `json:"id"`
		Slug       string    `json:"slug"`
		Difficulty string    `json:"difficulty"`
		CreatedAt  time.Time `json:"created_at"`
	}

	problems, err := h.usecase.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []respBody
	for _, p := range problems {
		r := respBody{
			ID:         p.ID,
			Slug:       p.Slug,
			Difficulty: p.Difficulty,
			CreatedAt:  p.CreatedAt,
		}

		resp = append(resp, r)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Errorf(errors.ErrEncodeJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProblemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf(errors.ErrParseIntID, err)
	}

	p, stmts, tests, err := h.usecase.GetById(r.Context(), id)
	if err != nil {
		h.log.Errorf(errors.ErrGetBySlug, err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"problem":    p,
		"statements": stmts,
		"tests":      tests,
	}); err != nil {
		h.log.Errorf(errors.ErrEncodeJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Errorf(errors.ErrParseIntID, err)
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		h.log.Errorf(errors.ErrDeleteProblem, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ProblemStatementDTO — простая структура для переноса данных о statement
type ProblemStatementDTO struct {
	Language  string `json:"language"`
	Title     string `json:"title"`
	Statement string `json:"statement"`
}

// ProblemTestDTO — для тестов
type ProblemTestDTO struct {
	ID             int64  `json:"id"` // Или int, если ID числовой
	InputData      string `json:"input_data"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       bool   `json:"is_sample"`
}

// CreateProblemRequest — DTO для запроса создания проблемы
type CreateProblemRequest struct {
	Slug       string                `json:"slug"`
	Difficulty string                `json:"difficulty"`
	Statements []ProblemStatementDTO `json:"statements"`
	Tests      []ProblemTestDTO      `json:"tests"`
}

// toDomainStatements — маппинг слайса DTO в domain
func toDomainStatements(dtos []ProblemStatementDTO) []domain.ProblemStatement {
	var statements []domain.ProblemStatement
	for _, dto := range dtos {
		statements = append(statements, domain.ProblemStatement{
			Language:  dto.Language,
			Title:     dto.Title,
			Statement: dto.Statement,
		})
	}
	return statements
}

// toDomainTests — аналогично для тестов
func toDomainTests(dtos []ProblemTestDTO) []domain.ProblemTest {
	var tests []domain.ProblemTest
	for _, dto := range dtos {
		tests = append(tests, domain.ProblemTest{
			ID:             dto.ID, // Если ID string в domain, иначе strconv
			InputData:      dto.InputData,
			ExpectedOutput: dto.ExpectedOutput,
			IsSample:       dto.IsSample,
		})
	}
	return tests
}
