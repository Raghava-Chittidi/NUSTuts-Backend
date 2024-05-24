package studentAuth

import (
	"NUSTuts-Backend/internal/api"
	data "NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authInfo struct {
	ID           int
	Email        string
	Username     string
	AccessToken  string
	RefreshToken string
}

type StudentUser struct {
	Name     string
	Email    string
	Password string
}

type authStudentUser struct {
	name         string
	email        string
	passwordHash string
}

var authStudentUserDB = map[string]authStudentUser{}

var DefaultStudentUserService studentUserService

type studentUserService struct {
}

func (studentUserService) createStudentUser(newStudentUser StudentUser) error {
	_, ok := authStudentUserDB[newStudentUser.Email]
	if ok {
		fmt.Println("Student user already exists")
		return errors.New("student user already exists")
	}
	passwordHash, err := util.GetPasswordHash(newStudentUser.Password)
	if err != nil {
		return err
	}
	newAuthStudentUser := authStudentUser{
		name:         newStudentUser.Name,
		email:        newStudentUser.Email,
		passwordHash: passwordHash,
	}
	authStudentUserDB[newStudentUser.Email] = newAuthStudentUser
	return nil
}

func SignUpAsStudent(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string
		Email    string
		Password string
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Name == "" {
		util.ErrorJSON(w, errors.New("Invalid username!"))
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
		Modules:  []string{},
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
	log.Default().Println("Log in")
}

func getStudentUser(r *http.Request) StudentUser {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	return StudentUser{
		Name:     name,
		Email:    email,
		Password: password,
	}
}
