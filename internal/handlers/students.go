package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

var students []string = []string{"Alice", "Bob", "Charlie"}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Create a student")
}

func GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert []string to JSON []byte
	studentBytes, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON []byte to the response
	w.Write(studentBytes)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(students[0]))
}

func UpdateStudentByID(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("First student changed to Zara")
	students[0] = "Zara"
	w.Write([]byte(students[0]))
}

func DeleteByID(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("First student deleted")
	students = students[1:]
}
