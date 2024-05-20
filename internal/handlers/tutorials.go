package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

var tutorials []string = []string{"CS1101S", "CS2040S", "CS2030S"}

func CreateTutorial(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Create a tutorial")
}

func GetTutorials(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert []string to JSON []byte
	tutorialsBytes, err := json.Marshal(tutorials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON []byte to the response
	w.Write(tutorialsBytes)
}

func GetTutorialsByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tutorials[0]))
}

func UpdateTutorialByID(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("First TA changed to CS2100")
	tutorials[0] = "CS2100"
	w.Write([]byte(tutorials[0]))
}
