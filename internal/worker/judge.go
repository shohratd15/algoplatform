// internal/worker/judge.go
package worker

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/log"
	"context"
	"time"
)

type JudgeWorker struct {
	submissionUsecase usecase.SubmissionUsecase
	log               log.Logger
}

func NewJudgeWorker(s usecase.SubmissionUsecase, l log.Logger) *JudgeWorker {
	return &JudgeWorker{submissionUsecase: s, log: l}
}

func (w *JudgeWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			subs, err := w.submissionUsecase.ListPending(ctx, 5)
			if err != nil {
				w.log.Errorf("fetch pending", err)
				continue
			}

			for _, s := range subs {
				w.process(ctx, s)
			}
		}
	}
}

func (w *JudgeWorker) process(ctx context.Context, s domain.Submission) {
	w.log.Infof("Processing submission %d", s.ID)
	_ = w.submissionUsecase.UpdateStatus(ctx, s.ID, domain.StatusRunning)

	// Пока просто имитация выполнения
	time.Sleep(2 * time.Second)
	_ = w.submissionUsecase.UpdateStatus(ctx, s.ID, domain.StatusAccepted)
}
