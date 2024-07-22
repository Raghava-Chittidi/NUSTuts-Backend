package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	data "NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
)

func LoginAsTeachingAssistant(w http.ResponseWriter, r *http.Request) {
	var payload api.LoginPayload

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if email provided exists
	teachingAssistant, err := data.GetTeachingAssistantByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("invalid credentials"), http.StatusNotFound)
		return
	}

	// Check if password provided is valid
	valid, err := util.VerifyPassword(payload.Password, teachingAssistant.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	authUser := auth.AuthenticatedUser{
		ID:          int(teachingAssistant.ID),
		Name:        teachingAssistant.Name,
		Email:       teachingAssistant.Email,
		Role:        auth.RoleTeachingAssistant,
	}

	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
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
		Tokens: 	 tokens,
	}
	
	// Generate and set new refresh cookie
	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedTeachingAssistant}, http.StatusOK)
}
