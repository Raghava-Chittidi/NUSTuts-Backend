package handlers

import (
	"log"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Sign up")
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log in")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log out")
}
