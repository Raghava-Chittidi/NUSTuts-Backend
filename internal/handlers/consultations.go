package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

func GetConsultationsForTutorialForDate(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	consultDate, err := chi.URLParam(r, "date")
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	consultations, err := dataaccess.GetAllConsultationsForTutorialForDate(consultDate)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.ConsultationsResponse{Consultations: *consultations}
	util.WriteJSON(w, api.Response{Message: "Consultations for tutorial for date fetched successfully!", 
		Data: res}, http.StatusOK)
}

func GetConsultationsForTutorialForStudent(w http.ResponseWriter, r *http.Request) {
	studentId, err := strconv.Atoi(chi.URLParam(r, "studentId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	consultations, err := dataaccess.GetAllConsultationsForTutorialForStudent(studentId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.ConsultationsResponse{Consultations: *consultations}
	util.WriteJSON(w, api.Response{Message: "Consultations for student fetched successfully!", Data: res}, http.StatusOK)
}

func UnbookConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = dataaccess.UnbookConsultationById(consultationId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Consultation unbooked successfully!"}, http.StatusOK)
}

func BookConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = dataaccess.BookConsultationById(consultationId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Consultation booked successfully!"}, http.StatusOK)
}

