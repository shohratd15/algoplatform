// internal/controller/http/handlers/problem.go
package handlers

type CreateProblemRequest struct {
	Slug       string                `json:"slug"`
	Difficulty string                `json:"difficulty"`
	Statements []ProblemStatementDTO `json:"statements"`
}

type ProblemStatementDTO struct {
	Language  string `json:"language"` // "en", "ru", "tk"
	Title     string `json:"title"`
	Statement string `json:"statement"`
}

type ProblemResponse struct {
	ID         int64                 `json:"id"`
	Slug       string                `json:"slug"`
	Difficulty string                `json:"difficulty"`
	Statements []ProblemStatementDTO `json:"statements"`
}
