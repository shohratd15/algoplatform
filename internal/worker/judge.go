// internal/worker/judge.go
package worker

// import (
// 	"context"
// 	// "os/exec"
// 	// "bytes"
// 	// "time"
// 	"algoplatform/internal/domain"
// )

// type InMemoryJudge struct {
// 	queue chan int64
// 	repo  domain.SubmissionRepo
// 	// problemRepo, testRepo, logger и т.д.
// }

// func NewInMemoryJudge(repo domain.SubmissionRepo) *InMemoryJudge {
// 	j := &InMemoryJudge{queue: make(chan int64, 1024), repo: repo}
// 	go j.loop()
// 	return j
// }

// func (j *InMemoryJudge) Enqueue(ctx context.Context, id int64) error {
// 	select {
// 	case j.queue <- id:
// 		return nil
// 	default:
// 		return context.DeadlineExceeded
// 	}
// }

// func (j *InMemoryJudge) loop() {
// 	for id := range j.queue {
// 		_ = j.repo.UpdateStatus(context.Background(), id, "running")
// 		// TODO: получить исходник/тесты, собрать двоичный файл, прогнать тесты с таймаутом.
// 		// MVP: выставить "accepted" как заглушку.
// 		_ = j.repo.UpdateStatus(context.Background(), id, "accepted")
// 	}
// }
