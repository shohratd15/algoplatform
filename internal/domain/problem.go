package domain

type Problem struct {
	ID         int64
	Slug       string
	Title      string
	Difficulty string // easy|medium|hard
}

type ProblemVersion struct {
	ID          int64
	ProblemID   int64
	Version     int
	TimeLimitMS int
	MemoryMB    int
	StatementMD string
}
