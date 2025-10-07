package domain

// JWT-сервис (генерация и парсинг)
type TokenService interface {
	Generate(userID int64, email, role string) (string, error)
	Parse(token string) (Claims, error)
}

type Claims struct {
	UserID int64
	Email  string
	Role   string
}

// Валидация входных DTO (обертка над go-playground/validator)
type Validator interface {
	Struct(v any) error
}
