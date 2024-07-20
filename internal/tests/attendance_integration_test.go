package tests

import (
	"NUSTuts-Backend/internal/api"
	"testing"

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
	assert.Equal(t, expected.AttendanceString.Code, actual.AttendanceString.Code)
	assert.Equal(t, expected.AttendanceString.ExpiresAt, actual.AttendanceString.ExpiresAt)
	assert.Equal(t, expected.AttendanceString.TutorialID, actual.AttendanceString.TutorialID)
}
