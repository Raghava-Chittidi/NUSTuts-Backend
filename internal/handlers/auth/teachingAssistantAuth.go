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

func LoginAsTA(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	ta, err := data.GetTeachingAssistantByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("TA with this email does not exist!"))
		return
	}

	valid, err := util.VerifyPassword(payload.Password, ta.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("Incorrect Password!"))
		return
	}

	tutorial, err := data.GetTutorialById(int(ta.TutorialID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	authenticatedTA := api.TeachingAssistantAuthResponse{
		ID:          int(ta.ID),
		Name:        ta.Name,
		Email:       ta.Email,
		Role:        auth.GetRoleByEmail(ta.Email),
		Tutorial:    *tutorial,
	}
	
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedTA}, http.StatusOK)
}
