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
		util.ErrorJSON(w, err)
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
		util.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	requests, err := dataaccess.GetPendingRequestsByTutorialId(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var data []api.RequestResponse
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
		util.ErrorJSON(w, err, http.StatusNotFound)
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
		util.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	err = dataaccess.RejectRequestById(requestId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Rejected request successfully!"}, http.StatusOK)
}