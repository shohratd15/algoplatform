package domain

import "time"

type Problem struct {
	ID         int64
	Slug       string
	Difficulty string
	CreatedAt  time.Time
}

type ProblemStatement struct {
	ProblemID int64
	Language  string
	Title     string
	Statement string
}

type ProblemTest struct {
	ID             int64
	ProblemID      int64
	InputData      string
	ExpectedOutput string
	IsSample       bool
}
