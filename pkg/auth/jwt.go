package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zulfikarmuzakir/e_procurement/internal/domain"
)

type JWTAuth struct {
	accessSecret  []byte
	refreshSecret []byte
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

func NewJWTAuth(accessSecret, refreshSecret string) *JWTAuth {
	return &JWTAuth{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

type Claims struct {
	UserID   int64
	Name     string
	Username string
	Email    string
	Role     string
	Status   string
	jwt.RegisteredClaims
}

func (j *JWTAuth) GenerateAccessToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *JWTAuth) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt: Set the expiration time to 7 days from now
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			// IssuedAt: Set the issued time to the current time
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

func (j *JWTAuth) ValidateToken(tokenString string, isAccessToken bool) (*Claims, error) {
	// Choose the appropriate secret key based on token type
	secretKey := j.accessSecret
	if !isAccessToken {
		secretKey = j.refreshSecret
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
