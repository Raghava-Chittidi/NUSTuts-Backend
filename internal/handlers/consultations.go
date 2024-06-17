package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func GetConsultationsForTutorialForDate(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	consultDate := r.URL.Query().Get("date")

	// Check if date is in the correct format
	_, err = time.Parse("02-01-2006", consultDate)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultations, err := dataaccess.GetAllConsultationsForTutorialForDate(tutorialId, consultDate)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if len(*consultations) < 2 {
		util.GenerateConsultationsForDate(tutorialId, consultDate)
		consultations, err = dataaccess.GetAllConsultationsForTutorialForDate(tutorialId, consultDate)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	res := api.ConsultationsResponse{Consultations: *consultations}
	util.WriteJSON(w, api.Response{Message: "Consultations for tutorial for date fetched successfully!", 
		Data: res}, http.StatusOK)
}

func GetConsultationsForTutorialForStudent(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	studentId, err := strconv.Atoi(chi.URLParam(r, "studentId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	consultations, err := dataaccess.GetAllConsultationsForTutorialForStudent(tutorialId, studentId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.ConsultationsResponse{Consultations: *consultations}
	util.WriteJSON(w, api.Response{Message: "Consultations for student fetched successfully!", Data: res}, http.StatusOK)
}

// func UpdateConsultationById(w http.ResponseWriter, r *http.Request) {
// 	// If book is true, book the consultation, else unbook it
// 	book, err := strconv.ParseBool(r.URL.Query().Get("book"))
// 	if err != nil {
// 		util.ErrorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
// 	if err != nil {
// 		util.ErrorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	_, claims, err := auth.AuthObj.VerifyToken(w, r)
// 	userID := claims.Subject // The "sub" claim is typically used for the user ID
// 	userIDInt, err := strconv.Atoi(userID)
// 	if err != nil {
// 		util.ErrorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	if book {
// 		dataaccess.BookConsultationById(consultationId, userIDInt)
// 	} else {
// 		dataaccess.UnbookConsultationById(consultationId, userIDInt)
// 	}
// }

func BookConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, claims, err := auth.AuthObj.VerifyToken(w, r)
	userID := claims.Subject // The "sub" claim is typically used for the user ID
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	dataaccess.BookConsultationById(consultationId, userIDInt)
}

func CancelConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_, claims, err := auth.AuthObj.VerifyToken(w, r)
	userID := claims.Subject // The "sub" claim is typically used for the user ID
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	dataaccess.UnbookConsultationById(consultationId, userIDInt)
}
