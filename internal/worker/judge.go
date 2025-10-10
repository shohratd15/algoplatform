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

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Judge Worker stopped.")
			return
		case <-ticker.C:
			subs, err := w.subUsecase.ListPending(ctx, 5)
			if err != nil {
				w.log.Errorf("failed to fetch pending submissions: %v", err)
				continue
			}

			// Если нет новых решений, просто ждем следующий тик.
			if len(subs) == 0 {
				continue
			}

			w.log.Debugf("Found %d pending submissions to process.", len(subs))
			for _, s := range subs {
				// Запускаем обработку каждого решения в отдельной горутине.
				go w.process(ctx, s)
			}
		}
	}
}

func (w *JudgeWorker) process(ctx context.Context, s domain.Submission) {
	if err := w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusRunning); err != nil {
		w.log.Errorf("Submission %d: failed to set status RUNNING: %v", s.ID, err)
		return
	}

	_, _, tests, err := w.probUsecase.GetById(ctx, s.ProblemID)
	if err != nil {
		w.log.Errorf("Submission %d: failed to get problem tests: %v", s.ID, err)
		// Обновляем статус на ERROR, если не удалось получить тесты
		_ = w.subUsecase.UpdateStatus(ctx, s.ID, domain.StatusError)
		return
	}

	finalStatus := domain.StatusAccepted

	for i, t := range tests {
		w.log.Debugf("Submission %d: Submitting test %d...", s.ID, i+1)

		// Отправляем код в Judge0
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

		// 4. Polling Loop: Ожидание результата
		pollInterval := 500 * time.Millisecond
		maxChecks := 120 // Максимум 20 секунд ожидания (40 * 500мс)

		testStatus, err := w.pollForStatus(ctx, token, pollInterval, maxChecks)

		if err != nil {
			w.log.Errorf("Submission %d, Test %d: Polling failed: %v", s.ID, i+1, err)
			finalStatus = domain.StatusError
			break
		}

		// 5. Проверка финального статуса теста
		if testStatus != judge.StatusAccepted {
			w.log.Debugf("Submission %d, Test %d failed with status ID: %d", s.ID, i+1, testStatus)
			// Конвертируем статус Judge0 (например, 4) в статус предметной области (WrongAnswer)
			finalStatus = w.mapJudgeStatusToDomain(testStatus)
			break
		}
	}

	// 6. Обновляем финальный статус решения
	if err := w.subUsecase.UpdateStatus(ctx, s.ID, finalStatus); err != nil {
		w.log.Errorf("Submission %d: failed to set final status %s: %v", s.ID, finalStatus, err)
	}
}

// pollForStatus опрашивает Judge0, пока не будет получен финальный статус или не истечет таймаут.
func (w *JudgeWorker) pollForStatus(ctx context.Context, token string, interval time.Duration, maxChecks int) (int, error) {
	for i := 0; i < maxChecks; i++ {
		time.Sleep(interval) // Ждем перед каждой проверкой

		res, err := w.client.GetResult(ctx, token)
		if err != nil {
			return 0, fmt.Errorf("get result failed: %w", err)
		}

		statusID := res.Status.ID

		// Если статус = Queued (1) или Processing (2), продолжаем опрос
		if statusID == judge.StatusInQueue || statusID == judge.StatusProcessing {
			continue
		}

		// Получен финальный статус
		return statusID, nil
	}

	// Если цикл завершился по таймауту
	return 0, fmt.Errorf("polling timeout reached for token %s", token)
}

// mapJudgeStatusToDomain конвертирует ID статуса Judge0 в константу домена.
func (w *JudgeWorker) mapJudgeStatusToDomain(judgeID int) string {
	switch judgeID {
	case judge.StatusAccepted:
		return domain.StatusAccepted
	case judge.StatusWrongAnswer:
		return domain.StatusWrong
	case 5: // Time Limit Exceeded
		return domain.StatusTimeLimit
	case 6: // Compilation Error
		return domain.StatusError // Или domain.StatusCompileError, если он у вас есть
	// Добавьте другие важные статусы (Runtime Error, Memory Limit)
	default:
		// Для всех остальных (включая Internal Error Judge0)
		return domain.StatusError
	}
}
