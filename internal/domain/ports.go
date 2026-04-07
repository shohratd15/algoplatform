package domain

// TokenService — генерация и парсинг JWT
type TokenService interface {
	GenerateAccess(userID int64, email, role string) (string, error)
	GenerateRefresh(userID int64, email, role string) (string, error)
	ParseAccess(token string) (Claims, error)
	ParseRefresh(token string) (Claims, error)
}

type Claims struct {
	UserID int64
	Email  string
	Role   string
}

// Validator — обёртка над go-playground/validator
type Validator interface {
	Struct(v any) error
}
