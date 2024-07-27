package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// Get all consultations for a tutorial for a specific date
func GetConsultationsForTutorialForDate(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultDate := r.URL.Query().Get("date")

	// Check if date is in the correct format
	_, err = time.Parse("2006-01-02", consultDate)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultations, err := dataaccess.GetAllConsultationsForTutorialForDate(tutorialId, consultDate)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// If there are no consultations for the tutorial for the date, generate them
	if len(*consultations) < 2 {
		util.GenerateConsultationsForDate(tutorialId, consultDate)
		consultations, err = dataaccess.GetAllConsultationsForTutorialForDate(tutorialId, consultDate)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	// Transform the consultations to include student and teaching assistant details
	consultationsResponse, err := transformConsultationsToConsultationsResponse(consultations)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.ConsultationsResponse{Consultations: *consultationsResponse}
	util.WriteJSON(w, api.Response{Message: "Consultations for tutorial for date fetched successfully!",
		Data: res}, http.StatusOK)
}

// Get all booked consultations booked by any student
// for a tutorial for a specific date for a teaching assistant
func GetBookedConsultationsForTutorialForTA(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	dateParam := r.URL.Query().Get("date")
	timeParam := r.URL.Query().Get("time")

	// Check if date is in the correct format
	_, err = time.Parse("2006-01-02", dateParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if time is in the correct format
	_, err = time.Parse("15:04", timeParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultations, err := dataaccess.GetBookedConsultationsForTutorialForTA(tutorialId, dateParam, timeParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Transform the consultations to include student and teaching assistant details
	consultationsResponse, err := transformConsultationsToConsultationsResponse(consultations)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := getBookedConsultationsResponse(consultationsResponse)
	util.WriteJSON(w, api.Response{Message: "Booked consultations for TA fetched successfully!", Data: res}, http.StatusOK)
}

// Get all booked consultations booked by a student in a tutorial
func GetBookedConsultationsForTutorialForStudent(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	studentId, err := strconv.Atoi(chi.URLParam(r, "studentId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	dateParam := r.URL.Query().Get("date")
	timeParam := r.URL.Query().Get("time")

	// Check if date is in the correct format
	_, err = time.Parse("2006-01-02", dateParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if time is in the correct format
	_, err = time.Parse("15:04", timeParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultations, err := dataaccess.GetBookedConsultationsForTutorialForStudent(tutorialId, studentId, dateParam, timeParam)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Transform the consultations to include student and teaching assistant details
	consultationsResponse, err := transformConsultationsToConsultationsResponse(consultations)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := getBookedConsultationsResponse(consultationsResponse)
	util.WriteJSON(w, api.Response{Message: "Booked consultations for student fetched successfully!", Data: res}, http.StatusOK)
}

// Transform consultations to include student and teaching assistant details
func transformConsultationsToConsultationsResponse(consultations *[]models.Consultation) (*[]api.ConsultationResponse, error) {
	consultationsResponse := make([]api.ConsultationResponse, len(*consultations))
	for i, consultation := range *consultations {
		transformedConsultation, err := transformConsultationToConsultationResponse(consultation)
		if err != nil {
			return nil, err
		}

		consultationsResponse[i] = *transformedConsultation
	}

	return &consultationsResponse, nil
}

// Get student and teaching assistant details for a consultation
func transformConsultationToConsultationResponse(consultation models.Consultation) (*api.ConsultationResponse, error) {
	// Get the student
	// If studentID is 0, then no student has booked the consultation
	var student models.Student
	if consultation.StudentID != 0 {
		studentPointer, err := dataaccess.GetStudentById(consultation.StudentID)
		if err != nil {
			return nil, err
		}
		student = *studentPointer
	}

	// Get the tutorial
	tutorial, err := dataaccess.GetTutorialById(consultation.TutorialID)
	if err != nil {
		return nil, err
	}

	// Get the teaching assistant
	teachingAssistant, err := dataaccess.GetTeachingAssistantById(tutorial.TeachingAssistantID)
	if err != nil {
		return nil, err
	}

	return &api.ConsultationResponse{
		ID:                consultation.ID,
		Tutorial:          *tutorial,
		Student:           student,
		TeachingAssistant: *teachingAssistant,
		Date:              consultation.Date,
		StartTime:         consultation.StartTime,
		EndTime:           consultation.EndTime,
		Booked:            consultation.Booked,
	}, nil
}

// Returns a response object containing the booked consultations grouped by date
func getBookedConsultationsResponse(consultations *[]api.ConsultationResponse) api.BookedConsultationsResponse {
	// Group consultations by date
	var groupedConsultations = make(map[string][]api.ConsultationResponse)
	for _, consultation := range *consultations {
		groupedConsultations[consultation.Date] = append(groupedConsultations[consultation.Date], consultation)
	}

	// Sort consultations by date of consultation, in format dd-mm-yyyy
	// Each element in the array is an object containing the date, and the consultations array for that date
	var sortedConsultations []api.BookedConsultationsByDate = make([]api.BookedConsultationsByDate, 0)
	for date, consults := range groupedConsultations {
		sortedConsultations = append(sortedConsultations, api.BookedConsultationsByDate{Date: date, Consultations: consults})
	}

	// Sort the array by date
	sort.Slice(sortedConsultations, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02", sortedConsultations[i].Date)
		date2, _ := time.Parse("2006-01-02", sortedConsultations[j].Date)
		return date1.Before(date2)
	})

	return api.BookedConsultationsResponse{BookedConsultations: sortedConsultations}
}

// Book a consultation by its id for a student
func BookConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Get the student id
	userID := r.URL.Query().Get("userId")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultation, err := dataaccess.BookConsultationById(consultationId, userIDInt)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Transform the consultation to include student and teaching assistant details
	consultationResponse, err := transformConsultationToConsultationResponse(*consultation)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := *consultationResponse
	util.WriteJSON(w, api.Response{Message: "Consultation succesfully booked", Data: res}, http.StatusOK)
}

// Cancel a consultation by its id for a student
func CancelConsultationById(w http.ResponseWriter, r *http.Request) {
	consultationId, err := strconv.Atoi(chi.URLParam(r, "consultationId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Get the student id
	userID := r.URL.Query().Get("userId")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	consultation, err := dataaccess.UnbookConsultationById(consultationId, userIDInt)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Transform the consultation to include student and teaching assistant details
	consultationResponse, err := transformConsultationToConsultationResponse(*consultation)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := *consultationResponse
	util.WriteJSON(w, api.Response{Message: "Consultation succesfully cancelled", Data: res}, http.StatusOK)
}
