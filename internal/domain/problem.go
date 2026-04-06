package domain

import "time"

type Problem struct {
	ID         int64     `json:"id"`
	Slug       string    `json:"slug"`
	Difficulty string    `json:"difficulty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProblemStatement struct {
	ProblemID int64  `json:"problem_id"`
	Language  string `json:"language"`
	Title     string `json:"title"`
	Statement string `json:"statement"`
}

type ProblemTest struct {
	ID             int64  `json:"id"`
	ProblemID      int64  `json:"problem_id"`
	InputData      string `json:"input_data"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       bool   `json:"is_sample"`
}
