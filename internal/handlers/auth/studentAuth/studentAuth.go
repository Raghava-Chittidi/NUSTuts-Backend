package studentAuth

import (
	"log"
	"net/http"
)

func SignUpAsStudent(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Sign up")
}

func LoginAsStudent(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log in")
}
