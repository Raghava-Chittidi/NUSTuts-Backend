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

// Retrieves all attendance records for a tutorial of a student
func GetStudentAttendance(w http.ResponseWriter, r *http.Request) {
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

	studentAttendance, err := dataaccess.GetStudentAttendance(tutorialId, studentId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.StudentAttendanceResponse{Attendance: *studentAttendance}
	util.WriteJSON(w, api.Response{Message: "Student attendance retrieved successfully!", Data: res}, http.StatusOK)
}

// Checks if a student is present in a tutorial
func CheckStudentAttendance(w http.ResponseWriter, r *http.Request) {
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

	attendance, err := dataaccess.GetTodayAttendanceByStudentId(studentId, tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Attendance record found!", Data: attendance.Present}, http.StatusOK)
}

// Retrieves all attendance records for a tutorial
func GetAllAttendanceForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	attendances, err := dataaccess.GetAllAttendanceByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendancesResponse, err := transformAttendancesToAttendancesResponse(attendances)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := getAttendanceListsByDateResponse(attendancesResponse)
	util.WriteJSON(w, api.Response{Message: "Attendance list retrieved successfully!", Data: res}, http.StatusOK)
}

// Retrieves attendance for a tutorial on the current date
func GetTodayAttendanceForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	attendances, err := dataaccess.GetTodayAttendanceByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendancesResponse, err := transformAttendancesToAttendancesResponse(attendances)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := api.AttendanceListResponse{Attendances: *attendancesResponse}
	util.WriteJSON(w, api.Response{Message: "Attendance list retrieved successfully!", Data: res}, http.StatusOK)
}

// Generates an attendance code for a tutorial on the current date
// When called, deletes current attendance records and string for the tutorial
func GenerateAttendanceCodeForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = dataaccess.DeleteGeneratedAttendanceString(tutorialId)
	if err != nil {
		// util.ErrorJSON(w, err, http.StatusInternalServerError)
		// return
	}

	err = dataaccess.DeleteTodayAttendanceByTutorialID(tutorialId)
	if err != nil {
		// util.ErrorJSON(w, err, http.StatusInternalServerError)
		// return
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

// Retrieves the attendance code for a tutorial on the current date
func GetAttendanceCodeForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	attendanceString, err := dataaccess.GetAttendanceStringByTutorialID(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	attendanceStringNotExpired, err := dataaccess.VerifyAttendanceCode(tutorialId, attendanceString.Code)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !attendanceStringNotExpired {
		util.WriteJSON(w, api.Response{Message: "Code has expired!", Data: nil}, http.StatusOK)
		return
	}

	attendanceStringResponse := api.AttendanceStringResponse{AttendanceString: *attendanceString}
	util.WriteJSON(w, api.Response{Message: "Code retrieved successfully!", Data: attendanceStringResponse}, http.StatusOK)
}

// Deletes the attendance string for a tutorial on the current date
func DeleteAttendanceString(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// var payload api.DeleteAttendanceStringPayload
	// err = util.ReadJSON(w, r, &payload)
	// if err != nil {
	// 	util.ErrorJSON(w, err, http.StatusBadRequest)
	// 	return
	// }

	err = dataaccess.DeleteGeneratedAttendanceString(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Expired code has been removed successfully!"}, http.StatusOK)
}

// Verifies and marks a student's attendance for a tutorial on the current date
func VerifyAndMarkStudentAttendance(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var payload api.MarkAttendancePayload
	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
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
		util.WriteJSON(w, api.Response{Message: "Your attendance has been marked successfully!"}, http.StatusOK)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Invalid code!", Data: nil}, http.StatusOK)
}

// Transforms an attendance object to an attendance response object
func transformAttendanceToAttendanceResponse(attendance models.Attendance) (*api.AttendanceResponse, error) {
	student, err := dataaccess.GetStudentById(attendance.StudentID)
	if err != nil {
		return nil, err
	}

	attendanceResponse := api.AttendanceResponse{
		ID:         attendance.ID,
		Student:    *student,
		TutorialID: attendance.TutorialID,
		Date:       attendance.Date,
		Present:    attendance.Present,
	}

	return &attendanceResponse, nil
}

// Transforms a list of attendance objects to a list of attendance response objects
func transformAttendancesToAttendancesResponse(attendances *[]models.Attendance) (*[]api.AttendanceResponse, error) {
	attendancesResponse := make([]api.AttendanceResponse, len(*attendances))
	for i, attendance := range *attendances {
		transformedAttendance, err := transformAttendanceToAttendanceResponse(attendance)
		if err != nil {
			return nil, err
		}

		attendancesResponse[i] = *transformedAttendance
	}

	return &attendancesResponse, nil
}

// Returns a response object containing the attendance lists grouped by date
func getAttendanceListsByDateResponse(attendances *[]api.AttendanceResponse) api.AttendanceListsByDateResponse {
	// Group attendances by date
	var groupedAttendances = make(map[string][]api.AttendanceResponse)
	for _, attendance := range *attendances {
		groupedAttendances[attendance.Date] = append(groupedAttendances[attendance.Date], attendance)
	}

	// Sort attendances by date of consultation, in format dd-mm-yyyy
	// Each element in the array is an object containing the date, and the attendances array for that date
	var sortedAttendances []api.AttendanceListByDate = make([]api.AttendanceListByDate, 0)
	for date, attendance := range groupedAttendances {
		sortedAttendances = append(sortedAttendances, api.AttendanceListByDate{Date: date, Attendance: attendance})
	}

	// // Sort the array by date in descending order
	sort.Slice(sortedAttendances, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02", sortedAttendances[i].Date)
		date2, _ := time.Parse("2006-01-02", sortedAttendances[j].Date)
		return date1.After(date2)
	})

	return api.AttendanceListsByDateResponse{AttendanceLists: sortedAttendances}
}
