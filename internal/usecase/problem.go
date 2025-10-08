package usecase

import (
	"algoplatform/internal/domain"
	"algoplatform/internal/repo/postgres"
	"context"
)

type ProblemUsecase interface {
	Create(ctx context.Context, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error
	List(ctx context.Context) ([]domain.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Problem, []domain.ProblemStatement, []domain.ProblemTest, error)
	Delete(ctx context.Context, id int64) error
}

type problemUsecase struct {
	repo postgres.ProblemRepository
}

func (u *problemUsecase) Create(ctx context.Context, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error {
	return u.repo.CreateProblem(ctx, p, stmts, tests)
}

func (u *problemUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.DeleteProblem(ctx, id)
}

func (u *problemUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Problem, []domain.ProblemStatement, []domain.ProblemTest, error) {
	return u.repo.GetProblemBySlug(ctx, slug)
}

func (u *problemUsecase) List(ctx context.Context) ([]domain.Problem, error) {
	return u.repo.GetAllProblems(ctx)
}

func NewProblemUsecase(repo postgres.ProblemRepository) *problemUsecase {
	return &problemUsecase{repo: repo}
}
