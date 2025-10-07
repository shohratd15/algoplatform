package domain

type Submission struct {
	ID             int64
	UserID         int64
	ProblemID      int64
	ProblemVersion int64
	LanguageID     int16
	SourceCode     string
	Status         string // queued|running|accepted|...
}
