package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func RefreshAuthStatus(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == auth.AuthObj.CookieName {
			claims := &auth.Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(auth.AuthObj.Secret), nil
			})
			if err != nil {
				util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
			
			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				util.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user := auth.AuthenticatedUser{
				ID: userId,
				Name: "",
				Email: "",
				Role: claims.Role,
			}

			if claims.Role.UserType == "student" {
				student, err := dataaccess.GetStudentById(userId)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}
				
				user.Name = student.Name
				user.Email = student.Email

				tokens, err := auth.AuthObj.GenerateTokens(&user)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}

				tutorials, err := dataaccess.GetTutorialsByStudentId(userId)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}

				authenticatedStudent := api.StudentAuthResponse{
					ID:          int(student.ID),
					Name:        student.Name,
					Email:       student.Email,
					Role:        auth.RoleStudent,
					Modules:	 student.Modules,
					Tutorials:	 *tutorials,
					Tokens: 	 tokens,
				}

				refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
				http.SetCookie(w, refreshCookie)
				util.WriteJSON(w, api.Response{Message: "Refreshed auth status successfully!", Data: authenticatedStudent}, http.StatusOK)
				return
			} else if claims.Role.UserType == "teachingAssistant" {
				teachingAssistant, err := dataaccess.GetTeachingAssistantById(userId)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}

				user.Name = teachingAssistant.Name
				user.Email = teachingAssistant.Email

				tokens, err := auth.AuthObj.GenerateTokens(&user)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}

				tutorial, err := dataaccess.GetTutorialById(teachingAssistant.TutorialID)
				if err != nil {
					util.ErrorJSON(w, err, http.StatusInternalServerError)
					return 
				}

				authenticatedTA := api.TeachingAssistantAuthResponse{
					ID:          int(teachingAssistant.ID),
					Name:        teachingAssistant.Name,
					Email:       teachingAssistant.Email,
					Role:        auth.RoleTeachingAssistant,
					Tutorial:	 *tutorial,
					Tokens: 	 tokens,
				}

				refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
				http.SetCookie(w, refreshCookie)
				util.WriteJSON(w, api.Response{Message: "Refreshed auth status successfully!", Data: authenticatedTA}, http.StatusOK)
				return
			} else {
				util.ErrorJSON(w, errors.New("invalid user"), http.StatusUnauthorized)
				return 
			}
		}
	}

	util.WriteJSON(w, api.Response{Message: "No refresh cookie!"}, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, auth.AuthObj.DeleteRefreshCookie())
	util.WriteJSON(w, api.Response{Message: "Logged out successfully!"}, http.StatusOK)
}
