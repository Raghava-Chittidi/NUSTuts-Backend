package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var AuthObj Auth

func InitialiseAuthObj() error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}
	
	AuthObj = Auth{
		AccessTokenExpiry: time.Minute * 15,
		RefreshTokenExpiry: time.Hour * 24,
		Issuer: "example.com",
		Audience: "example.com",
		Secret: os.Getenv("SECRET"),
		CookieDomain: "",
		CookieName: "session-cookie",
		CookiePath: "/",
	}

	return nil
}

func (a *Auth) GenerateTokens(user *AuthenticatedUser) (Tokens, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)

	claims["role"] = user.Role
	claims["email"] = user.Email
	claims["typ"] = "JWT"
	claims["iss"] = a.Issuer
	claims["aud"] = a.Audience
	claims["sub"] = fmt.Sprint(user.ID)
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().UTC().Add(a.AccessTokenExpiry).Unix()

	signedAccessToken, err := accessToken.SignedString([]byte (a.Secret))
	if err != nil {
		return Tokens{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["role"] = user.Role
	refreshTokenClaims["email"] = user.Email
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["exp"] = time.Now().UTC().Add(a.RefreshTokenExpiry).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte (a.Secret))
	if err != nil {
		return Tokens{}, err
	}

	tokens := Tokens{
		AccessToken: signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokens, nil
}

/* SameSite - If true can help to prevent cookie from being transmitted cross-site to protect 
   against CSRF attacks
   Secure - Cookie is only sent over HTTPS
   HttpOnly - Cookie cannot be accessed by javascript */
func (a *Auth) GenerateRefreshCookie(refreshTokenString string) *http.Cookie {
	return &http.Cookie{
		Name: a.CookieName,
		Path: a.CookiePath,
		Value: refreshTokenString,
		Domain: a.CookieDomain,
		MaxAge: int(a.RefreshTokenExpiry.Seconds()),
		Expires: time.Now().Add(a.RefreshTokenExpiry),
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		HttpOnly: true,
	}
}

func (a *Auth) DeleteRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name: a.CookieName,
		Path: a.CookiePath,
		Value: "",
		Domain: a.CookieDomain,
		MaxAge: -1,
		Expires: time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		HttpOnly: true,
	}
}

func (a *Auth) VerifyToken(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errors.New("no authorization header")
	}

	arr := strings.Split(authHeader, " ")
	if len(arr) != 2 || arr[0] != "Bearer" {
		return "", nil, errors.New("invalid authorization header")
	}

	token := arr[1]
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(a.Secret), nil
	})

	if err != nil {
		return "", nil, err
	}

	if claims.Issuer != a.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}