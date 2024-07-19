package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var messageTests = []api.CreateMesssagePayload{
	{UserType: "student", Content: "hello world 123!"},
	{UserType: "teachingAssistant", Content: "hello world 123!"},
}

// Asserts whether the two messages are equal by comparing their fields
func assertEqualMessage(t *testing.T, expected *models.Message, actual *models.Message) {
	assert.Equal(t, expected.DiscussionID, actual.DiscussionID)
	assert.Equal(t, expected.SenderID, actual.SenderID)
	assert.Equal(t, expected.UserType, actual.UserType)
	assert.Equal(t, expected.Content, actual.Content)
}

// Tests Create Message - Tests if both students and TAs can create messages with valid tutorial id
func TestValidCreateMessageForTutorial(t *testing.T) {
	for _, messageTest := range messageTests {
		// Create test Student, TA and Tutorial
		testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create a discussion
		err = dataaccess.CreateDiscussion(int(testTutorial.ID))
		assert.NoError(t, err)
		discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
		assert.NoError(t, err)

		// Send a request to the create message api endpoint
		if messageTest.UserType == "student" {
			messageTest.SenderID = int(testStudent.ID)
			_, status, err := CreateStudentAuthenticatedMockRequest(messageTest, fmt.Sprintf("/api/messages/%d", int(testTutorial.ID)), "POST", testStudent)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, status)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
			_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(messageTest, fmt.Sprintf("/api/messages/%d", int(testTutorial.ID)), "POST", testTeachingAssistant)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, status)
		}

		// Get the actual message that is created
		messages, err := dataaccess.GetMessagesByTutorialId(int(testTutorial.ID))
		assert.NoError(t, err)

		// Compare expected message that should be created with the actual message created
		assert.Equal(t, 1, len(*messages))
		expectedMessage := &models.Message{DiscussionID: discussionId, SenderID: messageTest.SenderID, UserType: messageTest.UserType, Content: messageTest.Content}
		assertEqualMessage(t, expectedMessage, &(*messages)[0])

		// Clean up
		CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteDiscussionById(discussionId)
		dataaccess.DeleteMessagesByDiscussionId(discussionId)
	}
}

// Tests create message - Tests if both students and TAs can create messages with invalid tutorial id
func TestInvalidTutorialIdCreateMessageForTutorial(t *testing.T) {
	for _, messageTest := range messageTests {
		// Create test Student, TA and Tutorial
		testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)

		// Create a discussion
		err = dataaccess.CreateDiscussion(int(testTutorial.ID))
		assert.NoError(t, err)
		discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
		assert.NoError(t, err)

		// Send a request to the create message api endpoint
		if messageTest.UserType == "student" {
			messageTest.SenderID = int(testStudent.ID)
			_, status, err := CreateStudentAuthenticatedMockRequest(messageTest, fmt.Sprintf("/api/messages/%d", -1), "POST", testStudent)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, status)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
			_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(messageTest, fmt.Sprintf("/api/messages/%d", -1), "POST", testTeachingAssistant)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, status)
		}

		// Clean up
		CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
		dataaccess.DeleteDiscussionById(discussionId)
		dataaccess.DeleteMessagesByDiscussionId(discussionId)
	}
}

// Tests get all messages - Tests if TAs can get all messages with valid tutorial id
func TestValidGetAllMessagesForTutorialTeachingAssistant(t *testing.T) {
	// Create test Student, TA and Tutorial
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
		
	// Create a discussion
	err = dataaccess.CreateDiscussion(int(testTutorial.ID))
	assert.NoError(t, err)
	discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
	assert.NoError(t, err)

	// Create messages
	for _, messageTest := range messageTests {
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
		}
		err = dataaccess.CreateMessage(discussionId, messageTest.SenderID, messageTest.UserType, messageTest.Content)
		assert.NoError(t, err)
	}

	// Send a request to the get all messages api endpoint
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/messages/%d", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get all actual messages fetched from response
	var messagesResponse api.MessagesResponse
	err = json.Unmarshal(resData, &messagesResponse)
	assert.NoError(t, err)
	actualMessagesResponse := messagesResponse.Messages

	// Generate the expected messages
	var expectedMessagesResponse []api.MessageResponse
	for _, messageTest := range messageTests {
		var sender string
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
			sender = testStudent.Name
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
			sender = testTeachingAssistant.Name
		}
		expectedMessagesResponse = append(expectedMessagesResponse, api.MessageResponse{
			TutorialID: int(testTutorial.ID), Sender: sender, SenderID: messageTest.SenderID, UserType: messageTest.UserType, Content: messageTest.Content,
		})
	}
	
	// Compare expected messages that should be fetched with the actual messages fetched
	actualMessagesResponseBytes, err := json.Marshal(actualMessagesResponse)
	assert.NoError(t, err)
	expectedMessagesResponseBytes, err := json.Marshal(expectedMessagesResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedMessagesResponseBytes), string(actualMessagesResponseBytes))

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	dataaccess.DeleteDiscussionById(discussionId)
	dataaccess.DeleteMessagesByDiscussionId(discussionId)
}

// Tests get all messages - Tests if TAs can get all messages with invalid tutorial id
func TestInvalidTutorialIDGetAllMessagesForTutorialTeachingAssistant(t *testing.T) {
	// Create test Student, TA and Tutorial
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
		
	// Create a discussion
	err = dataaccess.CreateDiscussion(int(testTutorial.ID))
	assert.NoError(t, err)
	discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
	assert.NoError(t, err)

	// Create messages
	for _, messageTest := range messageTests {
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
		}
		err = dataaccess.CreateMessage(discussionId, messageTest.SenderID, messageTest.UserType, messageTest.Content)
		assert.NoError(t, err)
	}

	// Send a request to the get all messages api endpoint
	_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/messages/%d", -1), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	dataaccess.DeleteDiscussionById(discussionId)
	dataaccess.DeleteMessagesByDiscussionId(discussionId)
}

// Tests get all messages - Tests if Students can get all messages with valid tutorial id
func TestValidGetAllMessagesForTutorialStudent(t *testing.T) {
	// Create test Student, TA and Tutorial
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
		
	// Create a discussion
	err = dataaccess.CreateDiscussion(int(testTutorial.ID))
	assert.NoError(t, err)
	discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
	assert.NoError(t, err)

	// Create messages
	for _, messageTest := range messageTests {
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
		}
		err = dataaccess.CreateMessage(discussionId, messageTest.SenderID, messageTest.UserType, messageTest.Content)
		assert.NoError(t, err)
	}

	// Send a request to the get all messages api endpoint
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/messages/%d", int(testTutorial.ID)), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get all actual messages fetched from response
	var messagesResponse api.MessagesResponse
	err = json.Unmarshal(resData, &messagesResponse)
	assert.NoError(t, err)
	actualMessagesResponse := messagesResponse.Messages

	// Generate the expected messages
	var expectedMessagesResponse []api.MessageResponse
	for _, messageTest := range messageTests {
		var sender string
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
			sender = testStudent.Name
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
			sender = testTeachingAssistant.Name
		}
		expectedMessagesResponse = append(expectedMessagesResponse, api.MessageResponse{
			TutorialID: int(testTutorial.ID), Sender: sender, SenderID: messageTest.SenderID, UserType: messageTest.UserType, Content: messageTest.Content,
		})
	}
	
	// Compare expected messages that should be fetched with the actual messages fetched
	actualMessagesResponseBytes, err := json.Marshal(actualMessagesResponse)
	assert.NoError(t, err)
	expectedMessagesResponseBytes, err := json.Marshal(expectedMessagesResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedMessagesResponseBytes), string(actualMessagesResponseBytes))

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	dataaccess.DeleteDiscussionById(discussionId)
	dataaccess.DeleteMessagesByDiscussionId(discussionId)
}

// Tests get all messages - Tests if Students can get all messages with invalid tutorial id
func TestInvalidTutorialIDGetAllMessagesForTutorialStudent(t *testing.T) {
	// Create test Student, TA and Tutorial
	testStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
		
	// Create a discussion
	err = dataaccess.CreateDiscussion(int(testTutorial.ID))
	assert.NoError(t, err)
	discussionId, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
	assert.NoError(t, err)

	// Create messages
	for _, messageTest := range messageTests {
		if (messageTest.UserType == "student") {
			messageTest.SenderID = int(testStudent.ID)
		} else {
			messageTest.SenderID = int(testTeachingAssistant.ID)
		}
		err = dataaccess.CreateMessage(discussionId, messageTest.SenderID, messageTest.UserType, messageTest.Content)
		assert.NoError(t, err)
	}

	// Send a request to the get all messages api endpoint
	_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/messages/%d", -1), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	dataaccess.DeleteDiscussionById(discussionId)
	dataaccess.DeleteMessagesByDiscussionId(discussionId)
}
