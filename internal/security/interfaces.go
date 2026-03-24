package security

import "github.com/golang-jwt/jwt/v5"

type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error // nil = match, non-nil = mismatch or failure
}

type TokenService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint, email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
