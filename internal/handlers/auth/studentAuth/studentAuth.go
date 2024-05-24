package studentAuth

import (
	"NUSTuts-Backend/internal/util"
	"errors"
	"fmt"
	"log"
	"net/http"
)

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
	newStudentUser := getStudentUser(r)
	err := DefaultStudentUserService.createStudentUser(newStudentUser)
	if err != nil {
		log.Println(err, "Student Sign Up Failed")
	}
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
