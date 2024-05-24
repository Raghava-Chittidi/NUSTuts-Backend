package middlewares

import (
	"NUSTuts-Backend/internal/auth"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		godotenv.Load("../../.env")
		clientUrl := os.Getenv("CLIENT_URL")
		w.Header().Set("Access-Control-Allow-Origin", clientUrl)
		w.Header().Set("Access-Control-Allow-Credentials", "true")


		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func AuthoriseUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if valid access token is present
		_, claims, err := auth.AuthObj.VerifyToken(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if role is allowed to access route
		privilege := claims.Role.Privilege
		urlPath := r.URL.Path

		// If the role privileges do not match the route they are trying to access
		if (privilege != auth.RoleStudent.Privilege && strings.Contains(urlPath, "students")) || 
			(privilege != auth.RoleTeachingAssistant.Privilege && strings.Contains(urlPath, "teaching-assistant")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}