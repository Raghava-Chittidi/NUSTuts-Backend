package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts whether the two students are equal by comparing their fields
func assertEqualStudent(t *testing.T, expected *models.Student, actual *models.Student) {
	expectedModules, err := json.Marshal(expected.Modules)
	assert.NoError(t, err)
	actualModules, err := json.Marshal(actual.Modules)
	assert.NoError(t, err)

	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
	assert.JSONEq(t, string(expectedModules), string(actualModules))
}

func TestStudentDataaccess(t *testing.T) {
	var testStudentId int
	var testStudentEmail string
	t.Run("Create Student", func(t *testing.T) {
		// Current no. of students in the test db should be 0
		var count int64
		database.DB.Table("students").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a student
		student, err := dataaccess.CreateStudent(testStudent.Name, testStudent.Email, testStudent.Password, testStudent.Modules)
		assert.NoError(t, err)
		testStudentId = int(student.ID)
		testStudentEmail = student.Email

		// Current no. of students in the test db should be 1
		database.DB.Table("students").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Get Student by ID", func(t *testing.T) {
		// Get the student
		student, err := dataaccess.GetStudentById(testStudentId)
		assert.NoError(t, err)

		// Compare expected student that should be fetched with the actual student fetched
		assertEqualStudent(t, &testStudent, student)
	})

	t.Run("Get Student by Email", func(t *testing.T) {
		// Get the student
		student, err := dataaccess.GetStudentByEmail(testStudentEmail)
		assert.NoError(t, err)

		// Compare expected student that should be fetched with the actual student fetched
		assertEqualStudent(t, &testStudent, student)
	})

	t.Run("Delete Student by Email", func(t *testing.T) {
		// Current no. of students in the test db should be 1
		var count int64
		database.DB.Table("students").Count(&count)
		assert.Equal(t, 1, int(count))

		// Delete the student created by the first test
		err := dataaccess.DeleteStudentByEmail(testStudentEmail)
		assert.NoError(t, err)

		// Current no. of students in the test db should be 0
		database.DB.Table("students").Count(&count)
		assert.Equal(t, 0, int(count))
	})
}