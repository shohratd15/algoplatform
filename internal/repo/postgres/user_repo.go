package postgres

import (
	"algoplatform/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

var _ UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	q := `INSERT INTO users (username, email, role_user, password_hash)
	      VALUES ($1,$2,$3,$4)`

	_, err := r.DB.Exec(ctx, q, u.Username, u.Email, u.Role, u.PasswordHash)

	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, username, email, role_user, password_hash, created_at
	       FROM users WHERE email=$1`

	var u domain.User
	if err := r.DB.QueryRow(ctx, q, email).Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.PasswordHash, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}
