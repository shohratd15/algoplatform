package usecase

import (
	"algoplatform/internal/domain"
	repo "algoplatform/internal/repo/postgres"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string, role string) error
	Login(ctx context.Context, email, password string) (*domain.User, error)
}

type userUsecase struct {
	repo repo.UserRepository
}

func NewUserUsecase(r repo.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (uc *userUsecase) Register(ctx context.Context, username, email, password string, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		Role:         role,
		PasswordHash: string(hash),
	}

	return uc.repo.Create(ctx, user)
}

func (uc *userUsecase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, err
	}

	return user, nil
}
