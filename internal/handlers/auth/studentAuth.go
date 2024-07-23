package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUpAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload api.StudentSignupPayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	// Check if name provided is non-empty
	if payload.Name == "" {
		util.ErrorJSON(w, errors.New("invalid name"))
		return
	}

	// Check if email provided is valid
	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("invalid email"))
		return
	}

	// Check if email is already in use
	_, err = dataaccess.GetStudentByEmail(payload.Email)
	if err != gorm.ErrRecordNotFound {
		util.ErrorJSON(w, errors.New("email is in use"))
		return
	}

	// Check if password provided is valid
	if payload.Password == "" || len(payload.Password) < 6 {
		util.ErrorJSON(w, errors.New("invalid password"))
		return
	}

	// Hash the password provided
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	student, err := dataaccess.CreateStudent(payload.Name, payload.Email, string(hashedPw), payload.Modules)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	authUser := auth.AuthenticatedUser{
		ID:          int(student.ID),
		Name:        student.Name,
		Email:       student.Email,
		Role:        auth.RoleStudent,
	}

	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tutorials, err := dataaccess.GetTutorialsByStudentId(int(student.ID))
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

	// Generate and set new refresh cookie
	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Student created successfully", Data: authenticatedStudent}, http.StatusCreated)
}

func LoginAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload api.LoginPayload

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	// Check if email provided exists
	student, err := dataaccess.GetStudentByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("invalid credentials"), http.StatusNotFound)
		return
	}

	// Check if password provided is valid
	valid, err := util.VerifyPassword(payload.Password, student.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	authUser := auth.AuthenticatedUser{
		ID:          int(student.ID),
		Name:        student.Name,
		Email:       student.Email,
		Role:        auth.RoleStudent,
	}

	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tutorials, err := dataaccess.GetTutorialsByStudentId(int(student.ID))
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

	// Generate and set new refresh cookie
	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedStudent}, http.StatusOK)
}
