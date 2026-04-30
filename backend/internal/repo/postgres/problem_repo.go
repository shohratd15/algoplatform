package postgres

import (
	"algoplatform/internal/domain"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProblemRepository interface {
	CreateProblem(ctx context.Context, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error
	UpdateProblem(ctx context.Context, id int64, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error
	GetAllProblems(ctx context.Context) ([]domain.Problem, error)
	GetProblemById(ctx context.Context, id int64) (*domain.Problem, []domain.ProblemStatement, []domain.ProblemTest, error)
	DeleteProblem(ctx context.Context, id int64) error
}

type ProblemRepo struct {
	db *pgxpool.Pool
}

func (r *ProblemRepo) CreateProblem(ctx context.Context, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error {
	return r.upsertProblem(ctx, 0, p, stmts, tests, false)
}

func (r *ProblemRepo) UpdateProblem(ctx context.Context, id int64, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest) error {
	return r.upsertProblem(ctx, id, p, stmts, tests, true)
}

func (r *ProblemRepo) upsertProblem(ctx context.Context, id int64, p *domain.Problem, stmts []domain.ProblemStatement, tests []domain.ProblemTest, isUpdate bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	problemID := id
	if isUpdate {
		tag, execErr := tx.Exec(ctx, `
			UPDATE problems
			SET slug = $1, difficulty = $2
			WHERE id = $3
		`, p.Slug, p.Difficulty, id)
		if execErr != nil {
			return execErr
		}
		if tag.RowsAffected() == 0 {
			return tx.Commit(ctx)
		}

		if _, execErr = tx.Exec(ctx, `DELETE FROM problem_statements WHERE problem_id = $1`, id); execErr != nil {
			return execErr
		}
		if _, execErr = tx.Exec(ctx, `DELETE FROM problem_tests WHERE problem_id = $1`, id); execErr != nil {
			return execErr
		}
	} else {
		err = tx.QueryRow(ctx, `
			INSERT INTO problems (slug, difficulty) VALUES ($1, $2) RETURNING id
			`, p.Slug, p.Difficulty).Scan(&problemID)
		if err != nil {
			return err
		}
	}

	for _, s := range stmts {
		_, err = tx.Exec(ctx, `
				INSERT INTO problem_statements (problem_id, language, title, statement)
				VALUES($1, $2, $3, $4)`,
			problemID, s.Language, s.Title, s.Statement)
		if err != nil {
			return err
		}
	}

	for _, t := range tests {
		_, err = tx.Exec(ctx, `
			INSERT INTO problem_tests (problem_id, input_data, expected_output, is_sample)
			VALUES($1, $2, $3, $4)`,
			problemID, t.InputData, t.ExpectedOutput, t.IsSample)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *ProblemRepo) DeleteProblem(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM problems WHERE id=$1`, id)
	return err
}

func (r *ProblemRepo) GetAllProblems(ctx context.Context) ([]domain.Problem, error) {
	rows, err := r.db.Query(ctx, `SELECT id, slug, difficulty, created_at FROM problems ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var problems []domain.Problem
	for rows.Next() {
		var p domain.Problem
		if err := rows.Scan(&p.ID, &p.Slug, &p.Difficulty, &p.CreatedAt); err != nil {
			return nil, err
		}

		problems = append(problems, p)
	}

	return problems, nil
}

func (r *ProblemRepo) GetProblemById(ctx context.Context, id int64) (*domain.Problem, []domain.ProblemStatement, []domain.ProblemTest, error) {
	var p domain.Problem
	err := r.db.QueryRow(ctx,
		`SELECT id, slug, difficulty, created_at FROM problems WHERE id=$1`, id).
		Scan(&p.ID, &p.Slug, &p.Difficulty, &p.CreatedAt)
	if err != nil {
		return nil, nil, nil, err
	}

	stmtRows, err := r.db.Query(ctx, `SELECT language, title, statement FROM problem_statements WHERE problem_id=$1`, p.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	defer stmtRows.Close()

	var stmts []domain.ProblemStatement
	for stmtRows.Next() {
		var s domain.ProblemStatement
		if err := stmtRows.Scan(&s.Language, &s.Title, &s.Statement); err != nil {
			return nil, nil, nil, err
		}

		s.ProblemID = p.ID
		stmts = append(stmts, s)
	}

	testsRows, err := r.db.Query(ctx, `SELECT id, input_data, expected_output, is_sample FROM problem_tests WHERE problem_id=$1`, p.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	defer testsRows.Close()

	var tests []domain.ProblemTest
	for testsRows.Next() {
		var t domain.ProblemTest
		if err := testsRows.Scan(&t.ID, &t.InputData, &t.ExpectedOutput, &t.IsSample); err != nil {
			return nil, nil, nil, err
		}

		t.ProblemID = p.ID
		tests = append(tests, t)
	}

	return &p, stmts, tests, nil
}

func NewProblemRepo(db *pgxpool.Pool) *ProblemRepo {
	return &ProblemRepo{db: db}
}

var _ ProblemRepository = (*ProblemRepo)(nil)
