package middlewares

import (
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// godotenv.Load("../../.env")
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
		if privilege != auth.RoleStudent.Privilege && strings.Contains(urlPath, "student") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if privilege != auth.RoleTeachingAssistant.Privilege && strings.Contains(urlPath, "teachingAssistant") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Ensure students cannot access any files and attendance routes without "student" inside the route 
		if privilege == auth.RoleStudent.Privilege && (strings.Contains(urlPath, "files") || strings.Contains(urlPath, "attendance")) && 
				!strings.Contains(urlPath, "student") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func ValidateTutorialID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := auth.AuthObj.VerifyToken(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userType := claims.Role.UserType
		url := r.URL.String()
		re := regexp.MustCompile(`(\d)+`)
		matches := re.FindAllString(url, -1)
		tutorialId, err := strconv.Atoi(matches[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Ensure Student/TA who is making this request for a tutorial, is in that tutorial
		if userType == "student" {
			valid, err := dataaccess.CheckIfStudentInTutorialById(id, tutorialId)
			if err != nil || !valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			valid, err := dataaccess.CheckIfTeachingAssistantInTutorialById(id, tutorialId)
			if err != nil || !valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}