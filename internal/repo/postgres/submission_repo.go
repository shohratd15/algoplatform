package postgres

import (
	"algoplatform/internal/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionRepository interface {
	Create(ctx context.Context, s *domain.Submission) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Submission, error)
	GetPending(ctx context.Context, limit int) ([]domain.Submission, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type SubmissionRepo struct {
	DB *pgxpool.Pool
}

func NewSubmissionRepo(db *pgxpool.Pool) *SubmissionRepo {
	return &SubmissionRepo{DB: db}
}

func (r *SubmissionRepo) Create(ctx context.Context, s *domain.Submission) (int64, error) {
	q := `
		INSERT INTO submissions (user_id, problem_id, language_id, source_code, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRow(ctx, q, s.UserID, s.ProblemID, s.LanguageID, s.SourceCode, s.Status, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SubmissionRepo) GetByID(ctx context.Context, id int64) (*domain.Submission, error) {
	q := `
		SELECT id, user_id, problem_id, language_id, source_code, status, created_at, updated_at
		FROM submissions WHERE id=$1
	`
	var s domain.Submission
	err := r.DB.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID, &s.SourceCode, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubmissionRepo) GetPending(ctx context.Context, limit int) ([]domain.Submission, error) {
	q := `
		SELECT id, user_id, problem_id, language_id, source_code, status, created_at, updated_at
		FROM submissions WHERE status=$1
		ORDER BY created_at ASC
		LIMIT $2
	`
	rows, err := r.DB.Query(ctx, q, domain.StatusQueued, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []domain.Submission
	for rows.Next() {
		var s domain.Submission
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID, &s.SourceCode, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	return subs, nil
}

func (r *SubmissionRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	q := `UPDATE submissions SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := r.DB.Exec(ctx, q, status, time.Now(), id)

	return err
}
