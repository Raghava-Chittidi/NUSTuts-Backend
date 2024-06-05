package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUpAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string
		Email    string
		Password string
		Modules  []string
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		util.ErrorJSON(w, errors.New("invalid name"), http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	_, err = dataaccess.GetStudentByEmail(payload.Email)
	if err != gorm.ErrRecordNotFound {
		util.ErrorJSON(w, errors.New("email is in use"), http.StatusBadRequest)
		return
	}

	if payload.Password == "" {
		util.ErrorJSON(w, errors.New("invalid password"), http.StatusBadRequest)
		return
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	student := models.Student{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPw),
		Modules:  payload.Modules,
	}

	result := database.DB.Table("students").Create(&student)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
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

	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Student created successfully", Data: authenticatedStudent}, http.StatusCreated)
}

func LoginAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	student, err := dataaccess.GetStudentByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("student with this email does not exist"), http.StatusNotFound)
		return
	}

	valid, err := util.VerifyPassword(payload.Password, student.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("incorrect password"), http.StatusUnauthorized)
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

	refreshCookie := auth.AuthObj.GenerateRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedStudent}, http.StatusOK)
}
