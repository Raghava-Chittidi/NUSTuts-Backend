package auth

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	data "NUSTuts-Backend/internal/dataaccess"
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
		util.ErrorJSON(w, err)
		return
	}

	if payload.Name == "" {
		util.ErrorJSON(w, errors.New("Invalid name!"))
		return
	}

	_, err = mail.ParseAddress(payload.Email)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, errors.New("Invalid email!"))
		return
	}

	_, err = data.GetStudentByEmail(payload.Email)
	if err != gorm.ErrRecordNotFound {
		util.ErrorJSON(w, errors.New("Email is in use!"))
		return
	}

	if payload.Password == "" {
		util.ErrorJSON(w, errors.New("Invalid password!"))
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
	log.Println(student)
	result := database.DB.Table("students").Create(&student)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Student created successfully"}, http.StatusCreated)
}

func LoginAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	student, err := data.GetStudentByEmail(payload.Email)
	if err != nil {
		util.ErrorJSON(w, errors.New("Student with this email does not exist!"))
		return
	}

	valid, err := util.VerifyPassword(payload.Password, student.Password)
	if err != nil || !valid {
		util.ErrorJSON(w, errors.New("Incorrect Password!"))
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
	}

	util.WriteJSON(w, api.Response{Message: "Login successful", Data: authenticatedStudent}, http.StatusOK)
}
