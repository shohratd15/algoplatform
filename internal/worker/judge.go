// internal/worker/judge.go
package worker

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/usecase"
	"algoplatform/pkg/judge"
	"algoplatform/pkg/log"
	"context"
	"fmt"
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
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	w.log.Info("Judge Worker started and listening for submissions.")

	const numWorkers = 5
	jobs := make(chan domain.Submission, 50)

	for i := range numWorkers {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					return
				case s, ok := <-jobs:
					if !ok {
						return
					}
					w.process(ctx, s)
				}
			}
		}(i)
	}

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Judge Worker stopped.")
			close(jobs)
			return
		case <-ticker.C:
			// NOTE о race condition:
			// repo.GetPending использует "SELECT ... FOR UPDATE SKIP LOCKED"
			// и атомарно переводит статус queued → running прямо в запросе.
			// Поэтому повторная выборка тех же submissions на следующем тике
			// невозможна — они уже имеют статус running и не попадут в выборку.
			subs, err := w.subUsecase.ListPending(ctx, numWorkers)
			if err != nil {
				w.log.Errorf("failed to fetch pending submissions: %v", err)
				continue
			}

			if len(subs) == 0 {
				continue
			}

			w.log.Debugf("Found %d pending submissions to process.", len(subs))
			for _, s := range subs {
				jobs <- s
			}
		}
	}
}

func (w *JudgeWorker) process(ctx context.Context, s domain.Submission) {
	_, _, tests, err := w.probUsecase.GetById(ctx, s.ProblemID)
	if err != nil {
		w.log.Errorf("Submission %d: failed to get problem tests: %v", s.ID, err)
		_ = w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusError)
		return
	}

	finalStatus := domain.StatusAccepted

	for i, t := range tests {
		w.log.Debugf("Submission %d: Submitting test %d...", s.ID, i+1)

		token, err := w.client.Submit(ctx, judge.SubmissionRequest{
			LanguageID: s.LanguageID,
			SourceCode: s.SourceCode,
			Stdin:      t.InputData,
			Expected:   t.ExpectedOutput,
		})
		if err != nil {
			w.log.Errorf("Submission %d, Test %d: Judge0 Submit failed: %v", s.ID, i+1, err)
			finalStatus = domain.StatusError
			break
		}

		testStatus, err := w.pollForStatus(ctx, token, 500*time.Millisecond, 120)
		if err != nil {
			w.log.Errorf("Submission %d, Test %d: Polling failed: %v", s.ID, i+1, err)
			finalStatus = domain.StatusError
			break
		}

		if testStatus != judge.StatusAccepted {
			w.log.Debugf("Submission %d, Test %d failed with status ID: %d", s.ID, i+1, testStatus)
			finalStatus = w.mapJudgeStatusToDomain(testStatus)
			break
		}
	}

	if err := w.subUsecase.UpdateStatus(ctx, s.ID, finalStatus); err != nil {
		w.log.Errorf("Submission %d: failed to set final status %s: %v", s.ID, finalStatus, err)
	}
}

func (w *JudgeWorker) pollForStatus(ctx context.Context, token string, interval time.Duration, maxChecks int) (int, error) {
	for range maxChecks {
		time.Sleep(interval)

		res, err := w.client.GetResult(ctx, token)
		if err != nil {
			return 0, fmt.Errorf("get result failed: %w", err)
		}

		statusID := res.Status.ID
		if statusID == judge.StatusInQueue || statusID == judge.StatusProcessing {
			continue
		}

		return statusID, nil
	}

	return 0, fmt.Errorf("polling timeout reached for token %s", token)
}

func (w *JudgeWorker) mapJudgeStatusToDomain(judgeID int) string {
	switch judgeID {
	case judge.StatusAccepted:
		return domain.StatusAccepted
	case judge.StatusWrongAnswer:
		return domain.StatusWrong
	case 5:
		return domain.StatusTimeLimit
	case 6:
		return domain.StatusError
	default:
		return domain.StatusError
	}
}
