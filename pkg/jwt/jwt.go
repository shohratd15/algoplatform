package jwt

import (
	"algoplatform/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var _ domain.TokenService = (*Service)(nil)

func New(secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) GenerateAccess(userID int64, email, role string) (string, error) {
	return s.generate(userID, email, role, "access", s.accessTTL)
}

func (s *Service) GenerateRefresh(userID int64, email, role string) (string, error) {
	return s.generate(userID, email, role, "refresh", s.refreshTTL)
}

func (s *Service) generate(userID int64, email, role, tokenType string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"type":    tokenType,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(s.secret)
}

func (s *Service) ParseAccess(token string) (domain.Claims, error) {
	return s.parseByType(token, "access")
}

func (s *Service) ParseRefresh(token string) (domain.Claims, error) {
	return s.parseByType(token, "refresh")
}

func (s *Service) parseByType(token, expectedType string) (domain.Claims, error) {
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
		tokenType, _ := claims["type"].(string)
		if tokenType != expectedType {
			return out, errors.New("invalid token type")
		}
		out.UserID = int64(claims["user_id"].(float64))
		out.Email = claims["email"].(string)
		out.Role = claims["role"].(string)

		return out, nil
	}

	return out, errors.New("invalid claims")
}
