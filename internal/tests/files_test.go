package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTeachingAssistant = models.TeachingAssistant{
	Name: "test_ta",
	Email: "test_ta@gmail.com",
	Password: "test_ta",
}

var testTutorial = models.Tutorial{
	TutorialCode: "123456",
	Module: "test_CS1101S",
}

var validfilesTests = []api.UploadFilePayload{
	{Name: "test_filename1", Week: 1, Filepath: "test_filepath1"},
	{Name: "test_filename2", Week: 6, Filepath: "test_filepath2"},
}

var invalidWeekfilesTests = []api.UploadFilePayload{
	{Name: "test_filename1", Week: -1, Filepath: "test_filepath2"},
	{Name: "test_filename2", Week: 20, Filepath: "test_filepath2"},
}

// Asserts whether the two tutorial files are equal by comparing their fields
func assertEqualTutorialFile(t *testing.T, expected *models.TutorialFile, actual *models.TutorialFile) {
	assert.Equal(t, expected.Filepath, actual.Filepath)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Visible, actual.Visible)
	assert.Equal(t, expected.Week, actual.Week)
	assert.Equal(t, expected.TutorialID, actual.TutorialID)
}

// Unique name with valid week number
func TestValidUploadFilepath(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)
	
		// Send a request to the upload file api endpoint
		_, status, err := CreateMockRequest(uploadFilePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
	
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)
	
		// Compare expected tutorial file that should be created with the actual
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Unique name with valid week number
func TestInvalidWeekUploadFilepath(t *testing.T) {
	for _, uploadFilePayload := range invalidWeekfilesTests {
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)
	
		// Send a request to the upload file api endpoint
		_, status, err := CreateMockRequest(uploadFilePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.Error(t, err)
		assert.Nil(t, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
	}
}