package security

import (
	"fmt"
	"time"

	"go-microservice/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(cfg config.JWTConfig) *JWTService {
	return &JWTService{
		secretKey:     cfg.SecretKey,
		accessExpiry:  cfg.AccessTokenExpiry,
		refreshExpiry: cfg.RefreshTokenExpiry,
	}
}

func (s *JWTService) GenerateAccessToken(userID uint, email string) (string, error) {
	return s.generateToken(userID, email, "access", s.accessExpiry)
}

func (s *JWTService) GenerateRefreshToken(userID uint, email string) (string, error) {
	return s.generateToken(userID, email, "refresh", s.refreshExpiry)
}

func (s *JWTService) generateToken(userID uint, email, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
