package tests

import (
	"NUSTuts-Backend/internal/api"
)

var messageTests = []api.CreateMesssagePayload{
	{SenderID: 0, UserType: "student", Content: "hello!"},
	{SenderID: 0, UserType: "student", Content: "hello world 123!"},
	{SenderID: 0, UserType: "teachingAssistant", Content: "hello world 123!"},
}

// func TestValidCreateMessageForTutorial(t *testing.T) {
// 	for i, messageTest := range messageTests {
// 		_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
// 		assert.NoError(t, err)

// 		// Send a request to the create message api endpoint
// 		_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(messageTest, fmt.Sprintf("/api/messages/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusCreated, status)

// 		// Get the actual tutorial file that is created
// 		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload.Name, filePayload.Week)
// 		assert.NoError(t, err)

// 		// Compare expected tutorial file that should be created with the actual file created
// 		expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: validFilesTests[i].Filepath, Name: validFilesTests[i].Name, Visible: true, Week: validFilesTests[i].Week}
// 		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)

// 		// Clean up
// 		CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
// 		dataaccess.DeleteTutorialFileByFilepath(tutorialFile.Filepath)
// 	}
// }