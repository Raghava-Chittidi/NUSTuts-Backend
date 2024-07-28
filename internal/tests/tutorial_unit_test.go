package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts whether the two tutorials are equal by comparing their fields
func assertEqualTutorial(t *testing.T, expected *models.Tutorial, actual *models.Tutorial) {
	assert.Equal(t, expected.TutorialCode, actual.TutorialCode)
	assert.Equal(t, expected.Module, actual.Module)
	assert.Equal(t, expected.TeachingAssistantID, actual.TeachingAssistantID)
}

func TestTutorialDataaccess(t *testing.T) {
	var testTutorialIds []int
	t.Run("Create Tutorial", func(t *testing.T) {
		// Current no. of tutorials in the test db should be 0
		var count int64
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a tutorial
		tutorial, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, testTutorial.TeachingAssistantID)
		assert.NoError(t, err)
		testTutorialIds = append(testTutorialIds, int(tutorial.ID))

		// Current no. of tutorials in the test db should be 1
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Get Tutorial by ID", func(t *testing.T) {
		// Get the tutorial
		tutorial, err := dataaccess.GetTutorialById(testTutorialIds[0])
		assert.NoError(t, err)

		// Compare expected tutorial that should be fetched with the actual tutorial fetched
		assertEqualTutorial(t, &testTutorial, tutorial)
	})

	t.Run("Get Tutorial by Class and Module Code", func(t *testing.T) {
		// Get the tutorial
		tutorial, err := dataaccess.GetTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		assert.NoError(t, err)

		// Compare expected tutorial that should be fetched with the actual tutorial fetched
		assertEqualTutorial(t, &testTutorial, tutorial)
	})

	t.Run("Check if valid Teaching Assistant in Tutorial by ID", func(t *testing.T) {
		// Get the boolean value
		isInside, err := dataaccess.CheckIfTeachingAssistantInTutorialById(int(testTeachingAssistant.ID), testTutorialIds[0])
		assert.NoError(t, err)

		// Compare expected boolean value with the actual boolean value
		assert.Equal(t, true, isInside)
	})

	t.Run("Check if invalid Teaching Assistant in Tutorial by ID", func(t *testing.T) {
		// Get the boolean value
		isInside, err := dataaccess.CheckIfTeachingAssistantInTutorialById(-1, testTutorialIds[0])
		assert.Error(t, err)

		// Compare expected boolean value with the actual boolean value
		assert.Equal(t, false, isInside)
	})

	// Create another tutorial
	tutorial, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, testTutorial.TeachingAssistantID)
	assert.NoError(t, err)
	testTutorialIds = append(testTutorialIds, int(tutorial.ID))

	t.Run("Get all Tutorial IDs", func(t *testing.T) {
		// Get all tutorial ids
		tutorialIds, err := dataaccess.GetAllTutorialIds()
		assert.NoError(t, err)
		assert.Equal(t, len(testTutorialIds), len(*tutorialIds))

		// Compare expected tutorial ids that should be fetched with the actual tutorial ids fetched
		expectedTutorialIdsBytes, err := json.Marshal(testTutorialIds)
		assert.NoError(t, err)
		actualTutorialIdsBytes, err := json.Marshal(tutorialIds)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedTutorialIdsBytes), string(actualTutorialIdsBytes))
	})

	t.Run("Delete Tutorial by ID", func(t *testing.T) {
		// Current no. of tutorials in the test db should be 2
		var count int64
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 2, int(count))

		// Delete the tutorial created by the first test
		err := dataaccess.DeleteTutorialById(testTutorialIds[0])
		assert.NoError(t, err)

		// Current no. of tutorials in the test db should be 1
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Delete Tutorial by Class and Module Code", func(t *testing.T) {
		// Current no. of tutorials in the test db should be 1
		var count int64
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 1, int(count))

		// Delete the second tutorial created
		err := dataaccess.DeleteTutorialByClassAndModuleCode(tutorial.TutorialCode, tutorial.Module)
		assert.NoError(t, err)

		// Current no. of tutorials in the test db should be 0
		database.DB.Table("tutorials").Count(&count)
		assert.Equal(t, 0, int(count))
	})
}