package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	data "NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
)

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
		util.ErrorJSON(w, errors.New("ta with this email does not exist"))
		return
	}

	valid, err := util.VerifyPassword(payload.Password, ta.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("incorrect password"))
		return
	}

	authUser := auth.AuthenticatedUser{
		ID:          int(ta.ID),
		Name:        ta.Name,
		Email:       ta.Email,
		Role:        auth.RoleTeachingAssistant,
	}

	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
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
		Role:        auth.RoleTeachingAssistant,
		Tutorial:    *tutorial,
		Tokens: 	 tokens,
	}
	
	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedTA}, http.StatusOK)
}
