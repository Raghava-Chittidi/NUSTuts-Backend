package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRegistry models.Registry
var testTutorial2 = models.Tutorial{
	TutorialCode: "910021",
	Module:       "test_MA1521",
}

// Asserts whether the two registries are equal by comparing their fields
func assertEqualRegistry(t *testing.T, expected *models.Registry, actual *models.Registry) {
	assert.Equal(t, expected.StudentID, actual.StudentID)
	assert.Equal(t, expected.TutorialID, actual.TutorialID)
}

func TestRegistryDataaccess(t *testing.T) {
	student1, err := dataaccess.CreateStudent(testStudent.Name, testStudent.Email, testStudent.Password, testStudent.Modules)
	assert.NoError(t, err)
	tutorial1, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, int(testTeachingAssistant.ID))
	assert.NoError(t, err)
	testRegistry.StudentID = int(student1.ID)
	testRegistry.TutorialID = int(tutorial1.ID)

	t.Run("Join Tutorial", func(t *testing.T) {
		// Current no. of registries in the test db should be 0
		var count int64
		database.DB.Table("registries").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a registry
		err := dataaccess.JoinTutorial(int(student1.ID), int(tutorial1.ID))
		assert.NoError(t, err)

		// Current no. of registries in the test db should be 1
		database.DB.Table("registries").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Get Registry by Student ID and Tutorial ID", func(t *testing.T) {
		// Get the registry
		registry, err := dataaccess.GetRegistryByStudentIDAndTutorialID(int(student1.ID), int(tutorial1.ID))
		assert.NoError(t, err)

		// Compare expected registry that should be fetched with the actual registry fetched
		assertEqualRegistry(t, &testRegistry, registry)
	})

	// Student joins another tutorial
	tutorial2, err := dataaccess.CreateTutorial(testTutorial2.TutorialCode, testTutorial2.Module, int(testTeachingAssistant.ID))
	assert.NoError(t, err)
	err = dataaccess.JoinTutorial(int(student1.ID), int(tutorial2.ID))
	assert.NoError(t, err)

	t.Run("Get Tutorials by Student ID", func(t *testing.T) {
		// Expected tutorials that the student should be in
		expectedTutorials := []models.Tutorial{*tutorial1, *tutorial2}

		// Get the actual tutorials that the student is in
		actualTutorials, err := dataaccess.GetTutorialsByStudentId(int(student1.ID))
		assert.NoError(t, err)

		// Compare expected tutorials that should be fetched with the actual tutorials fetched
		for j, expectedTutorial := range expectedTutorials {
			assertEqualTutorial(t, &expectedTutorial, &(*actualTutorials)[j])
		}
	})

	// Another Student joins the first tutorial
	student2, err := dataaccess.CreateStudent(testStudents[0].Name, testStudents[0].Email, testStudents[0].Password, testStudents[0].Modules)
	assert.NoError(t, err)
	err = dataaccess.JoinTutorial(int(student2.ID), int(tutorial1.ID))
	assert.NoError(t, err)

	t.Run("Get all Students IDs of Students in a Tutorial", func(t *testing.T) {
		// Expected student ids
		expectedStudentIds := []int{int(student1.ID), int(student2.ID)}

		// Get the student ids
		actualStudentIds, err := dataaccess.GetAllStudentIdsOfStudentsInTutorial(int(tutorial1.ID))
		assert.NoError(t, err)

		// Compare expected student ids that should be fetched with the actual student ids fetched
		assert.ElementsMatch(t, expectedStudentIds, *actualStudentIds)
	})

	t.Run("Check if valid Student in Tutorial by ID", func(t *testing.T) {
		// Get the boolean value
		isInside, err := dataaccess.CheckIfStudentInTutorialById(int(student1.ID), int(tutorial1.ID))
		assert.NoError(t, err)

		// Compare expected boolean value with the actual boolean value
		assert.Equal(t, true, isInside)
	})

	t.Run("Check if invalid Student in Tutorial by ID", func(t *testing.T) {
		// Get the boolean value
		isInside, err := dataaccess.CheckIfStudentInTutorialById(-1, int(tutorial1.ID))
		assert.Error(t, err)

		// Compare expected boolean value with the actual boolean value
		assert.Equal(t, false, isInside)
	})

	t.Run("Delete Registry by Student and Tutorial", func(t *testing.T) {
		// Current no. of registries in the test db should be 3
		var count int64
		database.DB.Table("registries").Count(&count)
		assert.Equal(t, 3, int(count))

		// Delete the registries created
		err = dataaccess.DeleteRegistryByStudentAndTutorial(student1, tutorial1)
		assert.NoError(t, err)
		err := dataaccess.DeleteRegistryByStudentAndTutorial(student1, tutorial2)
		assert.NoError(t, err)
		err = dataaccess.DeleteRegistryByStudentAndTutorial(student2, tutorial1)
		assert.NoError(t, err)

		// Current no. of registries in the test db should be 0
		database.DB.Table("registries").Count(&count)
		assert.Equal(t, 0, int(count))
	})

	// Clean up
	err = dataaccess.DeleteStudentByEmail(student1.Email)
	assert.NoError(t, err)
	err = dataaccess.DeleteStudentByEmail(student2.Email)
	assert.NoError(t, err)
	err = dataaccess.DeleteTutorialById(int(tutorial1.ID))
	assert.NoError(t, err)
	err = dataaccess.DeleteTutorialById(int(tutorial2.ID))
	assert.NoError(t, err)
}