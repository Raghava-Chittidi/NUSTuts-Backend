package handlers

import (
	"log"
	"net/http"
)

type Auth struct{}

func (s *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Sign up")
}

func (s *Auth) Login(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log in")
}

func (s *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log out")
}
