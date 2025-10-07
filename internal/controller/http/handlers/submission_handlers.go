// internal/controller/http/handlers/submission.go
package handlers

type CreateSubmissionRequest struct {
	ProblemID  int64  `json:"problem_id"`
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
}

type SubmissionResponse struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"` // queued, running, accepted, wrong_answer ...
	ProblemID  int64  `json:"problem_id"`
	LanguageID int    `json:"language_id"`
}
