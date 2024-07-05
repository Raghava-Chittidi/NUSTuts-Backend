package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
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
	{Name: "test_filename3", Week: 6, Filepath: "test_filepath3"},
}

var invalidWeekfilesTests = []api.UploadFilePayload{
	{Name: "test_filename1", Week: -1, Filepath: "test_filepath1"},
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

// Tests upload file - Unique name with valid week number
func TestValidUploadFilepath(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)
	
		// Send a request to the upload file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(uploadFilePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
	
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)
	
		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests upload file - Unique name with invalid week number
func TestInvalidWeekUploadFilepath(t *testing.T) {
	for _, uploadFilePayload := range invalidWeekfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)
	
		// Send a request to the upload file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(uploadFilePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
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

// Tests delete file - Unique name with valid filepath
func TestValidDeleteFilepath(t *testing.T) {
	for _, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the delete file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: uploadFilePayload.Filepath}, fmt.Sprintf("/api/files/delete/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
	
		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.Error(t, err)
		assert.Nil(t, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
	}
}

// Tests delete file - Unique name with invalid filepath
func TestInvalidFilepathDeleteFilepath(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the delete file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/delete/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	
		// Get the tutorial file that should not have been deleted due to invalid path
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests private file - Unique name with valid filepath
func TestValidPrivateFile(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: uploadFilePayload.Filepath}, fmt.Sprintf("/api/files/private/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
	
		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: false, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests private file - Unique name with invalid filepath
func TestInvalidFilepathPrivateFile(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/private/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	
		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests unprivate file - Unique name with valid filepath
func TestValidUnprivateFile(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file and set visibility to false
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)
		tutorialFile.Visible = false
		database.DB.Save(tutorialFile)
		
		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: uploadFilePayload.Filepath}, fmt.Sprintf("/api/files/unprivate/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
	
		// Get the tutorial file
		actualTutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: true, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, actualTutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests unprivate file - Unique name with invalid filepath
func TestInvalidFilepathUnprivateFile(t *testing.T) {
	for i, uploadFilePayload := range validfilesTests {
		// Create test tutorial and TA
		testTeachingAssistant, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
		assert.NoError(t, err)

		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)
		tutorialFile.Visible = false
		database.DB.Save(tutorialFile)

		// Send a request to the unprivate file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/unprivate/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	
		// Get the tutorial file
		actualTutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[i].Filepath, Name: validfilesTests[i].Name, Visible: false, Week: validfilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, actualTutorialFile)
	
		// Clean up
		dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
		dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests get files for Student - Valid URLParams System test?
// func TestValidGetFilesForStudent(t *testing.T) {
// 	// Create test tutorial and TA
// 	_, testTutorial, err := CreateMockTeachingAssistantAndMockTutorial()
// 	assert.NoError(t, err)
	
// 	for _, uploadFilePayload := range validfilesTests {
// 		// Create test tutorial file
// 		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
// 		assert.NoError(t, err)
// 	}

// 	// Send a request to the get files for student endpoint
// 	res, status, err := CreateMockRequest(nil, fmt.Sprintf("/api/files/student/%d/%d", int(testTutorial.ID), testTutorial.ID), "PATCH")
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, status)

// 	// Get actual tutorial files fetched from response
// 	var actualTutorialFiles = &[]models.TutorialFile{}
// 	err = json.Unmarshal(res, actualTutorialFiles)
// 	assert.NoError(t, err)

// 	// Compare expected tutorial file that should be created with the actual file created
// 	expectedTutorialFiles := &[]models.TutorialFile{
// 		{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[1].Filepath, Name: validfilesTests[1].Name, Visible: true, Week: validfilesTests[1].Week},
// 		{TutorialID: int(testTutorial.ID), Filepath: validfilesTests[2].Filepath, Name: validfilesTests[2].Name, Visible: true, Week: validfilesTests[2].Week},
// 	}

// 	for j, actualTutorialFile := range *actualTutorialFiles {
// 		assertEqualTutorialFile(t, &(*expectedTutorialFiles)[j], &actualTutorialFile)
// 	}

// 	// Clean up
// 	dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
// 	dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
// 	for _, uploadFilePayload := range validfilesTests {
// 		dataaccess.DeleteTutorialFileByFilepath(uploadFilePayload.Filepath)
// 	}
	
// }