package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

func GetAllTutorialFilesForTAs(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if week < 1 || week > 13 {
		util.ErrorJSON(w, errors.New("invalid week"))
		return
	}

	files, err := dataaccess.GetAllTutorialFilesFromTutorialIDAndWeek(tutorialId, week)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.TutorialFilesResponse{Files: *files}
	util.WriteJSON(w, api.Response{Message: "Tutorial Files fetched successfully!", Data: res}, http.StatusOK)
}

func GetAllTutorialFilesForStudents(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if week < 1 || week > 13 {
		util.ErrorJSON(w, errors.New("invalid week"))
		return
	}

	files, err := dataaccess.GetAllTutorialFilesFromTutorialIDAndWeek(tutorialId, week)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	*files = lo.Filter(*files, func(item models.TutorialFile, index int) bool {
		return item.Visible
	})

	res := api.TutorialFilesResponse{Files: *files}
	util.WriteJSON(w, api.Response{Message: "Tutorial Files fetched successfully!", Data: res}, http.StatusOK)
}

func UploadFilepath(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload api.UploadFilePayload
	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Week < 1 || payload.Week > 13 {
		util.ErrorJSON(w, errors.New("invalid week"))
		return
	}

	err = dataaccess.CheckIfNameExistsForTutorialIDAndWeek(tutorialId, payload.Name, payload.Week)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.CreateTutorialFile(tutorialId, payload.Name, payload.Week, payload.Filepath)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "File uploaded successfully!"}, http.StatusCreated)
}

func DeleteFilepath(w http.ResponseWriter, r *http.Request) {
	var payload api.FilepathPayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.DeleteTutorialFileByFilepath(payload.Filepath)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Filepath removed from table successfully!"}, http.StatusOK)
}

func PrivateFile(w http.ResponseWriter, r *http.Request) {
	var payload api.FilepathPayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.PrivateFileByFilepath(payload.Filepath)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "File privated successfully!"}, http.StatusOK)
}

func UnprivateFile(w http.ResponseWriter, r *http.Request) {
	var payload api.FilepathPayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.UnprivateFileByFilepath(payload.Filepath)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "File unprivated successfully!"}, http.StatusOK)
}