// internal/worker/judge.go
package worker

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/judge"
	"algoplatform/pkg/log"
	"context"
	"time"
)

type JudgeWorker struct {
	subUsecase  usecase.SubmissionUsecase
	probUsecase usecase.ProblemUsecase
	client      *judge.Client
	log         log.Logger
}

func NewJudgeWorker(
	sub usecase.SubmissionUsecase,
	prob usecase.ProblemUsecase,
	client *judge.Client,
	log log.Logger,
) *JudgeWorker {
	return &JudgeWorker{
		subUsecase:  sub,
		probUsecase: prob,
		client:      client,
		log:         log,
	}
}

func (w *JudgeWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			subs, err := w.subUsecase.ListPending(ctx, 5)
			if err != nil {
				w.log.Errorf("fetch pending", err)
				continue
			}

			for _, s := range subs {
				go w.process(ctx, s)
			}
		}
	}
}

func (w *JudgeWorker) process(ctx context.Context, s domain.Submission) {
	_ = w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusRunning)

	_, _, tests, err := w.probUsecase.GetById(ctx, s.ProblemID)
	if err != nil {
		w.log.Errorf("get problem tests", err)

		return
	}

	allPassed := true
	for _, t := range tests {
		token, err := w.client.Submit(ctx, judge.SubmissionRequest{
			LanguageID: s.LanguageID,
			SourceCode: s.SourceCode,
			Stdin:      t.InputData,
			Expected:   t.ExpectedOutput,
		})
		if err != nil {
			allPassed = false

			break
		}

		time.Sleep(2 * time.Second)
		res, _ := w.client.GetResult(ctx, token)
		if res.Status.ID != 3 {
			allPassed = false

			break
		}
	}

	if allPassed {
		_ = w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusAccepted)
	} else {
		_ = w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusWrong)
	}
}
