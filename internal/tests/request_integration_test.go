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

func TestValidRequestToJoinTutorial(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty, arbitray where
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create test tutorial
	_, err = dataaccess.CreateTutorial("1", "CS2040S", 50)
	assert.NoError(t, err)
	// Get the test tutorial
	testTutorial, err := dataaccess.GetTutorialByClassAndModuleCode("1", "CS2040S")
	assert.NoError(t, err)

	// Create a request to join tutorial
	requestToJoinTutPayload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "1",
	}
	_, status, err := CreateStudentAuthenticatedMockRequest(requestToJoinTutPayload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Current no. of requests in the test db should be 1
	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)
	// Assert that request count is 1
	assert.Equal(t, 1, len(requests))

	// Assert that the request is correct
	assert.Equal(t, int(testStudent.ID), requests[0].StudentID)
	assert.Equal(t, int(testTutorial.ID), int(requests[0].TutorialID))
	assert.Equal(t, "pending", requests[0].Status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	// Delete the test tutorial
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Tutorial{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestValidAcceptRequest(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create test tutorial
	_, err = dataaccess.CreateTutorial("1", "CS2040S", 50)
	assert.NoError(t, err)
	// Get the test tutorial
	_, err = dataaccess.GetTutorialByClassAndModuleCode("1", "CS2040S")
	assert.NoError(t, err)

	// Create a request to join tutorial
	requestToJoinTutPayload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "1",
	}
	_, status, err := CreateStudentAuthenticatedMockRequest(requestToJoinTutPayload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Current no. of requests in the test db should be 1
	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)
	// Assert that request count is 1
	assert.Equal(t, 1, len(requests))

	// Accept the request
	_, status, err = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d/accept", requests[0].ID), "PATCH", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get the request
	request, err := dataaccess.GetRequestById(int(requests[0].ID))
	assert.NoError(t, err)

	// Assert that the request is accepted
	assert.Equal(t, "accepted", request.Status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	// Delete the test tutorial
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Tutorial{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}
