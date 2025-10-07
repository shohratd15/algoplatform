package domain

import (
	"time"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	Role         string
	PasswordHash string
	CreateAt     time.Time
}
