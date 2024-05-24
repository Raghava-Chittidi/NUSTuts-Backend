package handlers

import (
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Log out")
}
