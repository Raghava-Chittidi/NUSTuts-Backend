package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RequestToJoinTutorial(w http.ResponseWriter, r *http.Request) {
	var payload api.RequestToJoinTutorialPayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	tutorial, err := dataaccess.GetTutorialByClassAndModuleCode(payload.ClassNo, payload.ModuleCode)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = dataaccess.CreateRequest(payload.StudentID, int(tutorial.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Request sent successfully!"}, http.StatusCreated)
}

func AllPendingRequestsForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	requests, err := dataaccess.GetPendingRequestsByTutorialId(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Used to force data to at least be an empty array instead of null due to how json handles nil slices
	// https://stackoverflow.com/questions/56200925/return-an-empty-array-instead-of-null-with-golang-for-json-return-with-gin
	var data []api.RequestResponse = make([]api.RequestResponse, 0)
	for _, request := range requests {
		student, err := dataaccess.GetStudentById(request.StudentID)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		data = append(data, api.RequestResponse{ID: int(request.ID), Name: student.Name, Email: student.Email})
	}

	res := api.Response{Message: "Requests fetched successfully", Data: data}

	util.WriteJSON(w, res, http.StatusOK)
}

func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	requestId, err := strconv.Atoi(chi.URLParam(r, "requestId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = dataaccess.AcceptRequestById(requestId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	request, err := dataaccess.GetRequestById(requestId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = dataaccess.JoinTutorial(request.StudentID, request.TutorialID)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Accepted request successfully!"}, http.StatusOK)
}

func RejectRequest(w http.ResponseWriter, r *http.Request) {
	requestId, err := strconv.Atoi(chi.URLParam(r, "requestId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = dataaccess.RejectRequestById(requestId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Rejected request successfully!"}, http.StatusOK)
}

func GetUnrequestedClassNo(w http.ResponseWriter, r *http.Request) {
	studentId, err := strconv.Atoi(chi.URLParam(r, "studentId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	moduleCode := chi.URLParam(r, "moduleCode")
	classNoArr, err := dataaccess.GetClassNoByStudentIdAndModuleCode(studentId, moduleCode)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Fetched class no. successfully!", Data: classNoArr}, http.StatusOK)
}
