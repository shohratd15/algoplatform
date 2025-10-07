package jwt

import (
	"algoplatform/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	ttl    time.Duration
}

var _ domain.TokenService = (*Service)(nil)

func New(secret string, ttl time.Duration) *Service {
	return &Service{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *Service) Generate(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(s.ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(s.secret)
}

func (s *Service) Parse(token string) (domain.Claims, error) {
	var out domain.Claims

	tok, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected method")
		}

		return s.secret, nil
	})
	if err != nil || !tok.Valid {
		return out, errors.New("invalid token")
	}

	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		out.UserID = int64(claims["user_id"].(float64))
		out.Email = claims["email"].(string)
		out.Role = claims["role"].(string)

		return out, nil
	}

	return out, errors.New("invalid claims")
}
