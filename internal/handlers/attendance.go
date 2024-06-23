package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GenerateAttendanceCodeForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.DeleteGeneratedAttendanceString(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	err = dataaccess.DeleteTodayAttendanceByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	err = dataaccess.GenerateTodayAttendanceByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendanceString, err := dataaccess.CreateRandomAttendanceString(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendanceStringResponse := api.AttendanceStringResponse{AttendanceString: *attendanceString}
	util.WriteJSON(w, api.Response{Message: "Code generated successfully!", Data: attendanceStringResponse}, http.StatusCreated)
}

func GetAttendanceCodeForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	attendanceString, err := dataaccess.GetAttendanceStringByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendanceStringExpired, err := dataaccess.VerifyAttendanceCode(tutorialId, attendanceString.Code)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if attendanceStringExpired {
		util.WriteJSON(w, api.Response{Message: "Code has expired!", Data: nil}, http.StatusNotFound)
		return
	}
		
	attendanceStringResponse := api.AttendanceStringResponse{AttendanceString: *attendanceString}
	util.WriteJSON(w, api.Response{Message: "Code retrieved successfully!", Data: attendanceStringResponse}, http.StatusOK)
}

func DeleteAttendanceString(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload api.DeleteAttendanceStringPayload
	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.DeleteGeneratedAttendanceString(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Expired code has been removed successfully!"}, http.StatusOK)
}

func VerifyAndMarkStudentAttendance(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload api.MarkAttendancePayload
	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	ok, err := dataaccess.VerifyAttendanceCode(tutorialId, payload.AttendanceCode)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if ok {
		err = dataaccess.MarkPresent(payload.StudentID, tutorialId)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}
	
	util.WriteJSON(w, api.Response{Message: "Your attendance has been marked successfully!"}, http.StatusOK)
}