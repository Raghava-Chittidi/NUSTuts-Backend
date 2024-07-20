package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts whether the two teaching assistants are equal by comparing their fields
func assertEqualTeachingAssistant(t *testing.T, expected *models.TeachingAssistant, actual *models.TeachingAssistant) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.TutorialID, actual.TutorialID)
}

func TestTeachingAssistantDataaccess(t *testing.T) {
	var testTeachingAssistantId int
	var testTeachingAssistantEmail string
	t.Run("Create teaching assistant", func(t *testing.T) {
		// Current no. of teaching assistants in the test db should be 0
		var count int64
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a teaching assistant
		teachingAssistant, err := dataaccess.CreateTeachingAssistant(testTeachingAssistant.Name, testTeachingAssistant.Email, testTeachingAssistant.Password)
		assert.NoError(t, err)
		testTeachingAssistantId = int(teachingAssistant.ID)
		testTeachingAssistantEmail = teachingAssistant.Email

		// Current no. of teaching assistants in the test db should be 1
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Get Teaching Assistant by ID", func(t *testing.T) {
		// Get the teaching assistant
		teachingAssistant, err := dataaccess.GetTeachingAssistantById(testTeachingAssistantId)
		assert.NoError(t, err)

		// Compare expected teaching assistant that should be fetched with the actual teaching assistant fetched
		assertEqualTeachingAssistant(t, &testTeachingAssistant, teachingAssistant)
	})

	t.Run("Get Teaching Assistant by Email", func(t *testing.T) {
		// Get the teaching assistant
		teachingAssistant, err := dataaccess.GetTeachingAssistantByEmail(testTeachingAssistantEmail)
		assert.NoError(t, err)

		// Compare expected teaching assistant that should be fetched with the actual teaching assistant fetched
		assertEqualTeachingAssistant(t, &testTeachingAssistant, teachingAssistant)
	})

	t.Run("Delete Teaching Assistant by ID", func(t *testing.T) {
		// Current no. of teaching assistants in the test db should be 1
		var count int64
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 1, int(count))

		// Delete the teaching assistant created by the first test
		err := dataaccess.DeleteTeachingAssistantById(testTeachingAssistantId)
		assert.NoError(t, err)

		// Current no. of teaching assistants in the test db should be 0
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 0, int(count))
	})

	// Recreate the teaching assistant
	_, err := dataaccess.CreateTeachingAssistant(testTeachingAssistant.Name, testTeachingAssistant.Email, testTeachingAssistant.Password)
	assert.NoError(t, err)

	t.Run("Delete Teaching Assistant by Email", func(t *testing.T) {
		// Current no. of teaching assistants in the test db should be 1
		var count int64
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 1, int(count))

		// Delete the teaching assistant that has been recreated
		err := dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistantEmail)
		assert.NoError(t, err)

		// Current no. of teaching assistants in the test db should be 0
		database.DB.Table("teaching_assistants").Count(&count)
		assert.Equal(t, 0, int(count))
	})
}