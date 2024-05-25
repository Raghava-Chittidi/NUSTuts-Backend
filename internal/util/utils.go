package util

import (
	"NUSTuts-Backend/internal/api"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Reads request body into data
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// Max request size
	maxBytes := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	return err
}

// Writes json data onto the response and sets status code
func WriteJSON(w http.ResponseWriter, resData api.Response, status int) error {
	res, err := json.Marshal(resData)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(res)
	return err
}

// Writes json error onto the response and sets status code
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	// Default is 400 Bad Request
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	resData := api.Response{Message: err.Error(), Error: err}
	WriteJSON(w, resData, statusCode)
}

func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}
