package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticatedUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
	// Hackish way to implement union type
	// StudentUser *models.Student           `json:"studentUser"`
	// TAUser      *models.TeachingAssistant `json:"taUser"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type Auth struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Audience           string
	Secret             string
	CookieDomain       string
	CookieName         string
	CookiePath         string
}
