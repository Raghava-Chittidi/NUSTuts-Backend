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
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Role can be used to identify whether a user is a Student or a Teaching Assistant for access control
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
