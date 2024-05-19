package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type TeachingAssistant struct{}

var teachingAssistants []string = []string{"Alex", "Zap", "Bobby"}

func (s *TeachingAssistant) CreateTA(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Create a TA")
}

func (s *TeachingAssistant) GetTAs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert []string to JSON []byte
	teachingAssistantsBytes, err := json.Marshal(teachingAssistants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON []byte to the response
	w.Write(teachingAssistantsBytes)
}

func (s *TeachingAssistant) GetTAByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(teachingAssistants[0]))
}

func (s *TeachingAssistant) UpdateTAByID(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("First TA changed to Lex")
	teachingAssistants[0] = "Lex"
	w.Write([]byte(teachingAssistants[0]))
}

func (s *TeachingAssistant) DeleteTAByID(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("First TA deleted")
	teachingAssistants = teachingAssistants[1:]
}
