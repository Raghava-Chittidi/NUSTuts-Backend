package teachingAssistantAuth

import (
	"log"
	"net/http"
)

func SignUpAsTA(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Sign up")
}

func LoginAsTA(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log in")
}
