package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var testTeachingAssistant = models.TeachingAssistant{
// 	Name: "test_ta",
// 	Email: "test_ta@gmail.com",
// 	Password: "test_ta",
// }

// var testStudent = models.Student{
// 	Name: "test_student",
// 	Email: "test_student@gmail.com",
// 	Password: "test_student",
// 	Modules: []string{"test_CS1101S"},
// }

// var testTutorial = models.Tutorial{
// 	TutorialCode: "123456",
// 	Module: "test_CS1101S",
// }

var validFilesTests = []api.UploadFilePayload{
	{Name: "test_filename1", Week: 1, Filepath: "test_filepath1"},
	{Name: "test_filename2", Week: 6, Filepath: "test_filepath2"},
	{Name: "test_filename3", Week: 6, Filepath: "test_filepath3"},
}

var invalidWeekFilesTests = []api.UploadFilePayload{
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
	for i, filePayload := range validFilesTests {
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Send a request to the upload file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(filePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)

		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: true, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests upload file - Unique name with invalid week number
func TestInvalidWeekUploadFilepath(t *testing.T) {
	for _, filePayload := range invalidWeekFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Send a request to the upload file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(filePayload, fmt.Sprintf("/api/files/upload/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, status)

		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.Error(t, err)
		assert.Nil(t, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
	}
}

// Tests delete file - Unique name with valid filepath
func TestValidDeleteFilepath(t *testing.T) {
	for _, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the delete file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: filePayload.Filepath}, fmt.Sprintf("/api/files/delete/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.Error(t, err)
		assert.Nil(t, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
	}
}

// Tests delete file - Unique name with invalid filepath
func TestInvalidFilepathDeleteFilepath(t *testing.T) {
	for i, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the delete file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/delete/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)

		// Get the tutorial file that should not have been deleted due to invalid path
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: true, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests private file - Unique name with valid filepath
func TestValidPrivateFile(t *testing.T) {
	for i, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: filePayload.Filepath}, fmt.Sprintf("/api/files/private/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: false, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests private file - Unique name with invalid filepath
func TestInvalidFilepathPrivateFile(t *testing.T) {
	for i, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/private/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)

		// Get the tutorial file
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: true, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests unprivate file - Unique name with valid filepath
func TestValidUnprivateFile(t *testing.T) {
	for i, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file and set visibility to false
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)
		tutorialFile.Visible = false
		database.DB.Save(tutorialFile)

		// Send a request to the private file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: filePayload.Filepath}, fmt.Sprintf("/api/files/unprivate/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get the tutorial file
		actualTutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: true, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, actualTutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests unprivate file - Unique name with invalid filepath
func TestInvalidFilepathUnprivateFile(t *testing.T) {
	for i, filePayload := range validFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create test tutorial file
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)

		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)
		tutorialFile.Visible = false
		database.DB.Save(tutorialFile)

		// Send a request to the unprivate file api endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(api.FilepathPayload{Filepath: "Invalid filepath"}, fmt.Sprintf("/api/files/unprivate/%d", int(testTutorial.ID)), "PATCH", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)

		// Get the tutorial file
		actualTutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: false, Week: validFilesTests[i].Week}
		assertEqualTutorialFile(t, expectedTutorialFile, actualTutorialFile)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
	}
}

// Tests get files for Student - Valid URLParams System test?
func TestValidGetFilesForStudent(t *testing.T) {
	// Create test Student, TA and Tutorial
	testStudent, _, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	for _, filePayload := range validFilesTests {
		// Create test tutorial files
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)
	}

	// Send a request to the get files for Student endpoint
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/student/%d/%d", int(testTutorial.ID), validFilesTests[1].Week), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual tutorial files fetched from response
	var tutorialFilesResponse api.TutorialFilesResponse
	err = json.Unmarshal(resData, &tutorialFilesResponse)
	assert.NoError(t, err)
	actualTutorialFiles := tutorialFilesResponse.Files

	// Compare expected tutorial file that should be created with the actual file created
	expectedTutorialFiles := &[]models.TutorialFile{
		{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[1].Filepath, Name: validFilesTests[1].Name, Visible: true, Week: validFilesTests[1].Week},
		{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[2].Filepath, Name: validFilesTests[2].Name, Visible: true, Week: validFilesTests[2].Week},
	}

	assert.Equal(t, len(*expectedTutorialFiles), len(actualTutorialFiles))
	for j, expectedTutorialFile := range *expectedTutorialFiles {
		assertEqualTutorialFile(t, &expectedTutorialFile, &actualTutorialFiles[j])
	}

	// Clean up
	CleanupCreatedStudentTeachingAssistantAndTutorial()
	for _, filePayload := range validFilesTests {
		dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
	}
}

// Tests get files for Student - Invalid URLParams System test?
func TestInvalidWeekGetFilesForStudent(t *testing.T) {
	for _, filePayload := range invalidWeekFilesTests {
		// Create test Student, TA and Tutorial
		testStudent, _, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		for _, uploadFilePayload := range validFilesTests {
			// Create test tutorial files
			err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
			assert.NoError(t, err)
		}

		// Send a request to the get files for Student endpoint
		_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/student/%d/%d", int(testTutorial.ID), filePayload.Week), "GET", testStudent)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, status)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		for _, filePayload := range validFilesTests {
			dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
		}
	}
}

// Tests get files for Student - Invalid Tutorial Id System test?
func TestInvalidTutorialIDGetFilesForStudent(t *testing.T) {
	for _, filePayload := range invalidWeekFilesTests {
		// Create test Student, TA and Tutorial
		testStudent, _, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		for _, uploadFilePayload := range validFilesTests {
			// Create test tutorial files
			err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
			assert.NoError(t, err)
		}

		// Send a request to the get files for Student endpoint
		_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/student/%d/%d", -1, filePayload.Week), "GET", testStudent)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, status)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		for _, filePayload := range validFilesTests {
			dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
		}
	}
}

// Tests get files for TA - Valid URLParams System test?
func TestValidGetFilesForTeachingAssistant(t *testing.T) {
	// Create test Student, TA and Tutorial
	_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	for _, filePayload := range validFilesTests {
		// Create test tutorial files
		err = dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload.Name, filePayload.Week, filePayload.Filepath)
		assert.NoError(t, err)
	}

	// Send a request to the get files for TA endpoint
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/teachingAssistant/%d/%d", int(testTutorial.ID), validFilesTests[1].Week), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual tutorial files fetched from response
	var tutorialFilesResponse api.TutorialFilesResponse
	err = json.Unmarshal(resData, &tutorialFilesResponse)
	assert.NoError(t, err)
	actualTutorialFiles := tutorialFilesResponse.Files

	// Compare expected tutorial file that should be created with the actual file created
	expectedTutorialFiles := &[]models.TutorialFile{
		{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[1].Filepath, Name: validFilesTests[1].Name, Visible: true, Week: validFilesTests[1].Week},
		{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[2].Filepath, Name: validFilesTests[2].Name, Visible: true, Week: validFilesTests[2].Week},
	}

	assert.Equal(t, len(*expectedTutorialFiles), len(actualTutorialFiles))
	for j, expectedTutorialFile := range *expectedTutorialFiles {
		assertEqualTutorialFile(t, &expectedTutorialFile, &actualTutorialFiles[j])
	}

	// Clean up
	CleanupCreatedStudentTeachingAssistantAndTutorial()
	for _, filePayload := range validFilesTests {
		dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
	}
}

// Tests get files for Student - Invalid URLParams System test?
func TestInvalidWeekGetFilesForTeachingAssistant(t *testing.T) {
	for _, filePayload := range invalidWeekFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		for _, uploadFilePayload := range validFilesTests {
			// Create test tutorial files
			err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
			assert.NoError(t, err)
		}

		// Send a request to the get files for TA endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/teachingAssistant/%d/%d", int(testTutorial.ID), filePayload.Week), "GET", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, status)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		for _, filePayload := range validFilesTests {
			dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
		}
	}
}

// Tests get files for Student - Invalid Tutorial Id System test?
func TestInvalidTutorialIDGetFilesForTeachingAssistant(t *testing.T) {
	for _, filePayload := range invalidWeekFilesTests {
		// Create test Student, TA and Tutorial
		_, testTeachingAssistant, testTutorial, err := CreateMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		for _, uploadFilePayload := range validFilesTests {
			// Create test tutorial files
			err = dataaccess.CreateTutorialFile(int(testTutorial.ID), uploadFilePayload.Name, uploadFilePayload.Week, uploadFilePayload.Filepath)
			assert.NoError(t, err)
		}

		// Send a request to the get files for TA endpoint
		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/files/teachingAssistant/%d/%d", -1, filePayload.Week), "GET", testTeachingAssistant)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, status)

		// Clean up
		CleanupCreatedStudentTeachingAssistantAndTutorial()
		for _, filePayload := range validFilesTests {
			dataaccess.DeleteTutorialFileByFilepath(filePayload.Filepath)
		}
	}
}

func TestC(t *testing.T) {
	CleanupCreatedStudentTeachingAssistantAndTutorial()
	// fmt.Print("cleaning")
}
