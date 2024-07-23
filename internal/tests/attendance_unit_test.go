package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAttendanceString(t *testing.T) {
	// Make sure attendance strings table is empty
	database.DB.Unscoped().Delete(&models.AttendanceString{})

	// Current no. of attendance strings in the test db should be 0
	var count int64
	database.DB.Table("attendance_strings").Count(&count)
	assert.Equal(t, 0, int(count))

	// Generate an attendance string
	attendanceString, err := dataaccess.CreateRandomAttendanceString(int(testTutorial.ID))
	assert.NoError(t, err)

	// Current no. of attendance strings in the test db should be 1
	database.DB.Table("attendance_strings").Count(&count)
	assert.Equal(t, 1, int(count))

	// Assert attendance string is not nil
	assert.NotNil(t, attendanceString)
	// Assert code is proper length
	assert.Equal(t, 10, len(attendanceString.Code))
	// Assert tutorial ID is correct
	assert.Equal(t, int(testTutorial.ID), attendanceString.TutorialID)
	// Assert expiry time is same as current time + expiry duration (within margin of error)
	assert.InDelta(t, attendanceString.ExpiresAt.Unix(), time.Now().Add(time.Minute*dataaccess.AttendanceCodeDuration).Unix(), 1)

	// Cleanup
	database.DB.Unscoped().Delete(attendanceString)
}

func TestDeleteAttendanceString(t *testing.T) {
	// Make sure attendance strings table is empty
	database.DB.Unscoped().Delete(&models.AttendanceString{})

	// Generate an attendance string
	_, err := dataaccess.CreateRandomAttendanceString(int(testTutorial.ID))
	assert.NoError(t, err)

	// Current no. of attendance strings in the test db should be 1
	var count int64
	database.DB.Table("attendance_strings").Count(&count)
	assert.Equal(t, 1, int(count))

	// Delete the attendance string
	err = dataaccess.DeleteGeneratedAttendanceString(int(testTutorial.ID))
	assert.NoError(t, err)

	// Current no. of attendance strings in the test db should be 0
	database.DB.Table("attendance_strings").Count(&count)
	assert.Equal(t, 0, int(count))

	// Cleanup
	database.DB.Unscoped().Delete(&models.AttendanceString{})
}
