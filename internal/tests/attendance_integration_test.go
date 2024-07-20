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
