package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var validfilesTests = []api.UploadFilePayload{
// 	{Name: "test_filename1", Week: 1, Filepath: "test_filepath2"},
// 	{Name: "test_filename2", Week: 6, Filepath: "test_filepath2"},
// }

// var invalidWeekfilesTests = []api.UploadFilePayload{
// 	{Name: "test_filename1", Week: -1, Filepath: "test_filepath2"},
// 	{Name: "test_filename2", Week: 20, Filepath: "test_filepath2"},
// }

// Asserts whether the two tutorial files are equal by comparing their fields
// func assertEqualTutorialFile(t *testing.T, expected *models.TutorialFile, student *models.TutorialFile) {
// 	assert.Equal(t, expected.Email, student.Email)
// 	assert.Equal(t, expected.Suspended, student.Suspended)
// 	assert.Equal(t, expected.RegisteredTeachers, student.RegisteredTeachers)
// }

// Unique name with valid week number
func TestValidUploadFilepath(t *testing.T) {
		// Create test TAs and tutorials in the database
		// log.Println(validfilesTests)
		// log.Println("tets")
		_, err := dataaccess.CreateTeachingAssistant("test_ta", "test_ta@gmail.com", "test_ta")
		assert.NoError(t, err)

		testTeachingAssistant, err := dataaccess.GetTeachingAssistantByEmail("test_ta@gmail.com")
		assert.NoError(t, err)
		log.Println(testTeachingAssistant.ID)
		assert.NoError(t, err)
	
		// testTutorial, err := dataaccess.CreateTutorial("123456", "test_CS1101S", int(testTeachingAssistant.ID))
		// assert.NoError(t, err)
	
		// Send a request to the upload file api endpoint
		// _, status, err := CreateMockRequest(uploadFilePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST")
		// assert.NoError(t, err)
		// assert.Equal(t, http.StatusCreated, status)
	
		// // Get the actual tutorial file that is created
		// tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		// assert.NoError(t, err)
	
		// // Compare expected tutorial file that should be created with the actual
		// expectedTutorialFile := models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		// assert.Equal(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		// err = dataaccess.DeleteTeachingAssistantById(int(testTeachingAssistant.ID))
		// assert.NoError(t, err)
		// dataaccess.DeleteTutorialById(int(testTutorial.ID))
		// dataaccess.DeleteTutorialFileById(int(tutorialFile.ID))
}
