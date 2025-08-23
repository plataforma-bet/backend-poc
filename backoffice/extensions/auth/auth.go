package auth

import (
	"backend-poc/backoffice/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type Subject struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}

type Claims struct {
	Subject
	jwt.RegisteredClaims
}

type JWTServices struct {
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	secret        string
}

func NewJWT(conf config.ConfigJWT) *JWTServices {
	return &JWTServices{
		accessExpiry:  time.Duration(conf.AccessExpiry) * time.Minute,
		refreshExpiry: time.Duration(conf.RefreshExpiry) * 24 * time.Hour,
		secret:        conf.Secret,
	}
}

func (j *JWTServices) GenerateAccessToken(user Subject) (string, error) {
	claims := &Claims{
		Subject: Subject{
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "backend-poc",
			Subject:   user.UserID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTServices) GenerateRefreshToken(user Subject) (string, error) {
	claims := &Claims{
		Subject: Subject{
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "backend-poc",
			Subject:   user.UserID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTServices) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTServices) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Gerar novos tokens
	user := Subject{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}

	newAccessToken, err := j.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := j.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (j *JWTServices) GetAccessTokenExpiry() time.Duration {
	return j.accessExpiry
}

func (j *JWTServices) GetRefreshTokenExpiry() time.Duration {
	return j.refreshExpiry
}
