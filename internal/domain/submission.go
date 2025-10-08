package domain

import "time"

type Submission struct {
	ID         int64
	UserID     int64
	ProblemID  int64
	LanguageID int
	SourceCode string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Возможные статусы
const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusAccepted  = "accepted"
	StatusWrong     = "wrong_answer"
	StatusError     = "error"
	StatusTimeLimit = "time_limit"
	StatusMemory    = "memory_limit"
)
