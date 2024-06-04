package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	data "NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
)

// Commented as sign up for TA is not needed
// func SignUpAsTA(w http.ResponseWriter, r *http.Request) {
// 	log.Default().Println("Sign up")
// }

func LoginAsTeachingAssistant(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	teachingAssistant, err := data.GetTeachingAssistantByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("Teaching Assistant with this email does not exist!"), http.StatusNotFound)
		return
	}

	valid, err := util.VerifyPassword(payload.Password, teachingAssistant.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("Incorrect Password!"), http.StatusUnauthorized)
		return
	}

	tutorial, err := data.GetTutorialById(int(teachingAssistant.TutorialID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	authenticatedTeachingAssistant := api.TeachingAssistantAuthResponse{
		ID:          int(teachingAssistant.ID),
		Name:        teachingAssistant.Name,
		Email:       teachingAssistant.Email,
		Role:        auth.RoleTeachingAssistant,
		Tutorial:    *tutorial,
	}
	
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedTeachingAssistant}, http.StatusOK)
}
