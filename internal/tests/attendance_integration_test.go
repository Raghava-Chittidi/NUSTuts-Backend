package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Asserts whether the two attendance response are equal by comparing their fields
func assertEqualAttendanceResponse(t *testing.T, expected *api.AttendanceResponse, actual *api.AttendanceResponse) {
	assert.Equal(t, expected.Date, actual.Date)
	assert.Equal(t, expected.Present, actual.Present)
	assert.Equal(t, expected.Student.ID, actual.Student.ID)
	assert.Equal(t, expected.TutorialID, actual.TutorialID)
}

// Asserts whether the two attendance lists by date response are equal by comparing their fields
func assertEqualAttendanceListsByDateResponse(t *testing.T, expected *api.AttendanceListsByDateResponse, actual *api.AttendanceListsByDateResponse) {
	assert.Equal(t, len(expected.AttendanceLists), len(actual.AttendanceLists))
	for i, expectedAttendanceListByDate := range expected.AttendanceLists {
		actualAttendanceListByDate := actual.AttendanceLists[i]
		assert.Equal(t, expectedAttendanceListByDate.Date, actualAttendanceListByDate.Date)
		assert.Equal(t, len(expectedAttendanceListByDate.Attendance), len(actualAttendanceListByDate.Attendance))
		for j, expectedAttendance := range expectedAttendanceListByDate.Attendance {
			assertEqualAttendanceResponse(t, &expectedAttendance, &actualAttendanceListByDate.Attendance[j])
		}
	}
}

// Asserts whether the two attendance lists response are equal by comparing their fields
func assertEqualAttendanceListsResponse(t *testing.T, expected *api.AttendanceListResponse, actual *api.AttendanceListResponse) {
	assert.Equal(t, len(expected.Attendances), len(actual.Attendances))

	for i, expectedAttendance := range expected.Attendances {
		actualAttendance := actual.Attendances[i]
		assertEqualAttendanceResponse(t, &expectedAttendance, &actualAttendance)
	}
}

// Asserts whether the two attendance strings are equal by comparing their fields
func assertEqualAttendanceStrings(t *testing.T, expected *api.AttendanceStringResponse, actual *api.AttendanceStringResponse) {
	// assert.Equal(t, expected.AttendanceString.Code, actual.AttendanceString.Code)
	// assert expire at within margin of error
	assert.InDelta(t, expected.AttendanceString.ExpiresAt.Unix(), actual.AttendanceString.ExpiresAt.Unix(), 1)
	assert.Equal(t, expected.AttendanceString.TutorialID, actual.AttendanceString.TutorialID)
}

// Test valid generate attendance code for tutorial
func TestValidGenerateAttendanceCodeForTutorial(t *testing.T) {
	_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var attendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &attendanceStringResponse)
	assert.NoError(t, err)

	const AttendanceCodeDuration = 5
	// Generate expected attendance string for the tutorial
	expectedAttendanceString := api.AttendanceStringResponse{
		AttendanceString: models.AttendanceString{
			Code:       attendanceStringResponse.AttendanceString.Code,
			ExpiresAt:  time.Now().Add(time.Minute * AttendanceCodeDuration), // within margin of error testing
			TutorialID: int(testTutorial.ID),
		},
	}
	assertEqualAttendanceStrings(t, &expectedAttendanceString, &attendanceStringResponse)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid get attendance code for tutorial
func TestValidGetAttendanceCodeForTutorial(t *testing.T) {
	_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	// Send a request to get attendance code for the tutorial
	res, status, err = CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ = json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var attendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &attendanceStringResponse)
	assert.NoError(t, err)
	assertEqualAttendanceStrings(t, &generatedAttendanceStringResponse, &attendanceStringResponse)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid delete attendance code for tutorial
func TestValidDeleteAttendanceCodeForTutorial(t *testing.T) {
	_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Send a request to delete attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/delete", int(testTutorial.ID)), "POST", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Expired code has been removed successfully!", response.Message)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid student mark attendance
func TestValidStudentMarkAttendance(t *testing.T) {
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	markAttendancePayload := api.MarkAttendancePayload{
		StudentID:      int(testStudent.ID),
		AttendanceCode: generatedAttendanceStringResponse.AttendanceString.Code,
	}

	res, status, err = CreateStudentAuthenticatedMockRequest(markAttendancePayload, fmt.Sprintf("/api/attendance/student/%d/mark", int(testTutorial.ID)), "POST", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response2 api.Response
	err = json.Unmarshal(res, &response2)
	assert.NoError(t, err)
	assert.Nil(t, response2.Data)
	assert.Equal(t, "Your attendance has been marked successfully!", response2.Message)

	attendance, err := dataaccess.GetTodayAttendanceByStudentId(int(testStudent.ID), int(testTutorial.ID))
	assert.NoError(t, err)
	assert.NotNil(t, attendance)
	assert.Equal(t, int(testStudent.ID), attendance.StudentID)
	assert.Equal(t, int(testTutorial.ID), attendance.TutorialID)
	assert.Equal(t, time.Now().Format("2006-01-02"), attendance.Date)
	assert.Equal(t, true, attendance.Present)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test student mark attendance with wrong code
func TestIncorrectCodeStudentMarkAttendance(t *testing.T) {
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	markAttendancePayload := api.MarkAttendancePayload{
		StudentID:      int(testStudent.ID),
		AttendanceCode: generatedAttendanceStringResponse.AttendanceString.Code + "wrong",
	}

	_, status, _ = CreateStudentAuthenticatedMockRequest(markAttendancePayload, fmt.Sprintf("/api/attendance/student/%d/mark", int(testTutorial.ID)), "POST", testStudent)
	assert.Equal(t, http.StatusInternalServerError, status)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid check student attendance
func TestValidCheckStudentAttendance(t *testing.T) {
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.Nil(t, err)
	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	// Check student attendance
	res, status, err = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/student/%d/attended/%d", int(testTutorial.ID), int(testStudent.ID)), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var attendedResponse api.Response
	err = json.Unmarshal(res, &attendedResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Attendance record found!", attendedResponse.Message)
	assert.Equal(t, false, attendedResponse.Data)

	// Mark attendance
	markAttendancePayload := api.MarkAttendancePayload{
		StudentID:      int(testStudent.ID),
		AttendanceCode: generatedAttendanceStringResponse.AttendanceString.Code,
	}

	_, status, _ = CreateStudentAuthenticatedMockRequest(markAttendancePayload, fmt.Sprintf("/api/attendance/student/%d/mark", int(testTutorial.ID)), "POST", testStudent)
	assert.Equal(t, http.StatusOK, status)

	// Check student attendance
	res, status, err = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/student/%d/attended/%d", int(testTutorial.ID), int(testStudent.ID)), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	err = json.Unmarshal(res, &attendedResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Attendance record found!", attendedResponse.Message)
	assert.Equal(t, true, attendedResponse.Data)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid get today attendance for tutorial
func TestValidGetTodayAttendanceForTutorial(t *testing.T) {
	testStudentModels := []models.Student{}
	testDefaultStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	for _, student := range testStudents {
		// Create test TeachingAssistant, Student, Tutorial
		student, err := CreateMockStudent(&student, testTeachingAssistant, testTutorial)
		assert.NoError(t, err)
		testStudentModels = append(testStudentModels, *student)
	}

	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	// Send a request to get today attendance for the tutorial
	res, status, err = CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/list", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	// Get response in json
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ = json.Marshal(response.Data)

	var attendanceListResponse api.AttendanceListResponse
	err = json.Unmarshal(resData, &attendanceListResponse)
	assert.NoError(t, err)

	// Expected today attendance list for the tutorial
	expectedAttendanceList := api.AttendanceListResponse{
		Attendances: []api.AttendanceResponse{
			{
				Student:    *testDefaultStudent,
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    false,
			},
			{
				Student:    testStudentModels[0],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    false,
			},
			{
				Student:    testStudentModels[1],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    false,
			},
			{
				Student:    testStudentModels[2],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    false,
			},
		},
	}

	assertEqualAttendanceListsResponse(t, &expectedAttendanceList, &attendanceListResponse)

	for i, student := range testStudentModels {
		// Mark attendance
		markAttendancePayload := api.MarkAttendancePayload{
			StudentID:      int(testStudentModels[i].ID),
			AttendanceCode: generatedAttendanceStringResponse.AttendanceString.Code,
		}

		_, status, _ = CreateStudentAuthenticatedMockRequest(markAttendancePayload, fmt.Sprintf("/api/attendance/student/%d/mark", int(testTutorial.ID)), "POST", &student)
		assert.Equal(t, http.StatusOK, status)
	}

	// Send a request to get today attendance for the tutorial
	res, status, err = CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/list", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ = json.Marshal(response.Data)

	var attendanceListResponseMarked api.AttendanceListResponse
	err = json.Unmarshal(resData, &attendanceListResponseMarked)
	assert.NoError(t, err)

	// Expected today attendance list for the tutorial
	expectedAttendanceListMarked := api.AttendanceListResponse{
		Attendances: []api.AttendanceResponse{
			{
				Student:    *testDefaultStudent,
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    false,
			},
			{
				Student:    testStudentModels[0],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    true,
			},
			{
				Student:    testStudentModels[1],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    true,
			},
			{
				Student:    testStudentModels[2],
				TutorialID: int(testTutorial.ID),
				Date:       time.Now().Format("2006-01-02"),
				Present:    true,
			},
		},
	}

	assertEqualAttendanceListsResponse(t, &expectedAttendanceListMarked, &attendanceListResponseMarked)

	// Clean up
	dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	dataaccess.DeleteTodayAttendanceByTutorialID(int(testTutorial.ID))

	// Clean up students, ta, tutorial
	for _, student := range testStudentModels {
		CleanupCreatedStudent(&student)
	}
	CleanupCreatedStudent(testDefaultStudent)
	CleanupCreatedTeachingAssistant(testTeachingAssistant)
	CleanupCreatedTutorial(testTutorial)
}

// Test valid get all attendance for tutorial
func TestValidGetAllAttendanceForTutorial(t *testing.T) {
	testStudentModels := []models.Student{}
	testDefaultStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	for _, student := range testStudents {
		// Create test TeachingAssistant, Student, Tutorial
		student, err := CreateMockStudent(&student, testTeachingAssistant, testTutorial)
		assert.NoError(t, err)
		testStudentModels = append(testStudentModels, *student)
	}

	// Send a request to generate attendance code for the tutorial
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/generate", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var generatedAttendanceStringResponse api.AttendanceStringResponse
	err = json.Unmarshal(resData, &generatedAttendanceStringResponse)
	assert.NoError(t, err)

	for i, student := range testStudentModels {
		// Mark attendance
		markAttendancePayload := api.MarkAttendancePayload{
			StudentID:      int(testStudentModels[i].ID),
			AttendanceCode: generatedAttendanceStringResponse.AttendanceString.Code,
		}

		_, status, _ = CreateStudentAuthenticatedMockRequest(markAttendancePayload, fmt.Sprintf("/api/attendance/student/%d/mark", int(testTutorial.ID)), "POST", &student)
		assert.Equal(t, http.StatusOK, status)
	}

	// Generate 2 attendance records for date 1 day before today, 2 days before today
	for i := 1; i <= 2; i++ {
		dataaccess.GenerateAttendanceForDateByTutorialID(time.Now().AddDate(0, 0, -i).Format("2006-01-02"), int(testTutorial.ID))
	}

	// Send a request to get all attendance for the tutorial
	res, status, err = CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/attendance/%d/list/all", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ = json.Marshal(response.Data)

	var attendanceListsResponse api.AttendanceListsByDateResponse
	err = json.Unmarshal(resData, &attendanceListsResponse)
	assert.NoError(t, err)

	// Expected today attendance list for the tutorial
	expectedAttendanceList := api.AttendanceListsByDateResponse{
		AttendanceLists: []api.AttendanceListByDate{
			{
				Date: time.Now().Format("2006-01-02"),
				Attendance: []api.AttendanceResponse{
					{
						Student:    *testDefaultStudent,
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[0],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().Format("2006-01-02"),
						Present:    true,
					},
					{
						Student:    testStudentModels[1],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().Format("2006-01-02"),
						Present:    true,
					},
					{
						Student:    testStudentModels[2],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().Format("2006-01-02"),
						Present:    true,
					},
				},
			},
			{
				Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
				Attendance: []api.AttendanceResponse{
					{
						Student:    *testDefaultStudent,
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[0],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[1],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[2],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
						Present:    false,
					},
				},
			},
			{
				Date: time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
				Attendance: []api.AttendanceResponse{
					{
						Student:    *testDefaultStudent,
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[0],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[1],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
						Present:    false,
					},
					{
						Student:    testStudentModels[2],
						TutorialID: int(testTutorial.ID),
						Date:       time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
						Present:    false,
					},
				},
			},
		},
	}

	assertEqualAttendanceListsByDateResponse(t, &expectedAttendanceList, &attendanceListsResponse)
}
