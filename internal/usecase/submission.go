package usecase

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/repo/postgres"
	"context"
)

type SubmissionUsecase interface {
	Create(ctx context.Context, s *domain.Submission) (int64, error)
	Get(ctx context.Context, id int64) (*domain.Submission, error)
	ListPending(ctx context.Context, limit int) ([]domain.Submission, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type submissionUsecase struct {
	repo postgres.SubmissionRepository
}

func NewSubmissionUsecase(r postgres.SubmissionRepository) SubmissionUsecase {
	return &submissionUsecase{repo: r}
}

func (u *submissionUsecase) Create(ctx context.Context, s *domain.Submission) (int64, error) {
	return u.repo.Create(ctx, s)
}

func (u *submissionUsecase) Get(ctx context.Context, id int64) (*domain.Submission, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *submissionUsecase) ListPending(ctx context.Context, limit int) ([]domain.Submission, error) {
	return u.repo.GetPending(ctx, limit)
}

func (u *submissionUsecase) UpdateStatus(ctx context.Context, id int64, status string) error {
	return u.repo.UpdateStatus(ctx, id, status)
}
