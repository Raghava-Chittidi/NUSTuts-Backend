package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequest(t *testing.T) {
	// Make sure requests table is empty
	database.DB.Unscoped().Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create a request
	err := dataaccess.CreateRequest(1, 1)
	assert.NoError(t, err)

	// Current no. of requests in the test db should be 1
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 1, int(count))

	// Cleanup
	database.DB.Unscoped().Delete(&models.Request{})
}

func TestGetRequest(t *testing.T) {
	// Make sure requests table is empty
	database.DB.Unscoped().Delete(&models.Request{})

	// Create a request
	err := dataaccess.CreateRequest(1, 1)
	assert.NoError(t, err)

	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)

	// Get the request
	expectedRequest := requests[0]

	// Get the request by ID
	request, err := dataaccess.GetRequestById(int(expectedRequest.ID))
	assert.NoError(t, err)

	// Assert request is not nil
	assert.NotNil(t, request)
	// Assert request is correct
	assert.Equal(t, expectedRequest.ID, request.ID)
	assert.Equal(t, expectedRequest.StudentID, request.StudentID)
	assert.Equal(t, expectedRequest.TutorialID, request.TutorialID)
	assert.Equal(t, expectedRequest.Status, request.Status)

	// Cleanup
	database.DB.Unscoped().Delete(&models.Request{})
}

func TestAcceptRequest(t *testing.T) {
	// Make sure requests table is empty
	database.DB.Unscoped().Delete(&models.Request{})

	// Create a request
	err := dataaccess.CreateRequest(1, 1)
	assert.NoError(t, err)

	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)

	// Get the request
	expectedRequest := requests[0]

	// Accept the request
	err = dataaccess.AcceptRequestById(int(expectedRequest.ID))
	assert.NoError(t, err)

	// Get the request
	request, err := dataaccess.GetRequestById(int(expectedRequest.ID))
	assert.NoError(t, err)

	// Assert request is not nil
	assert.NotNil(t, request)
	// Assert request is accepted
	assert.Equal(t, "accepted", request.Status)

	// Cleanup
	database.DB.Unscoped().Delete(&models.Request{})
}

func TestRejectRequest(t *testing.T) {
	// Make sure requests table is empty
	database.DB.Unscoped().Delete(&models.Request{})

	// Create a request
	err := dataaccess.CreateRequest(1, 1)
	assert.NoError(t, err)

	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)

	// Get the request
	expectedRequest := requests[0]

	// Reject the request
	err = dataaccess.RejectRequestById(int(expectedRequest.ID))
	assert.NoError(t, err)

	// Get the request
	request, err := dataaccess.GetRequestById(int(expectedRequest.ID))
	assert.NoError(t, err)

	// Assert request is not nil
	assert.NotNil(t, request)
	// Assert request is rejected
	assert.Equal(t, "rejected", request.Status)

	// Cleanup
	database.DB.Unscoped().Delete(&models.Request{})
}

func TestGetPendingRequest(t *testing.T) {
	// Make sure requests table is empty
	database.DB.Unscoped().Delete(&models.Request{})

	// Create 3 requests, 1 pending, 1 accepted, 1 rejected
	err := dataaccess.CreateRequest(1, 1)
	assert.NoError(t, err)
	err = dataaccess.CreateRequest(2, 1)
	assert.NoError(t, err)
	err = dataaccess.CreateRequest(3, 1)
	assert.NoError(t, err)

	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)

	// Get the expected pending request, expected accepted request, expected rejected request
	var expectedPendingRequest *models.Request
	var expectedAcceptedRequest *models.Request
	var expectedRejectedRequest *models.Request
	for _, request := range requests {
		if request.StudentID == 1 {
			expectedPendingRequest = request
		} else if request.StudentID == 2 {
			expectedAcceptedRequest = request
		} else if request.StudentID == 3 {
			expectedRejectedRequest = request
		}
	}

	// Accept the accepted request and reject the rejected request
	err = dataaccess.AcceptRequestById(int(expectedAcceptedRequest.ID))
	assert.NoError(t, err)
	err = dataaccess.RejectRequestById(int(expectedRejectedRequest.ID))
	assert.NoError(t, err)

	// Get all pending requests in the test db
	pendingRequests, err := dataaccess.GetPendingRequestsByTutorialId(1)
	assert.NoError(t, err)
	// Assert there is only 1 pending request
	assert.Equal(t, 1, len(pendingRequests))
	// Assert pending request is correct
	assert.Equal(t, expectedPendingRequest.ID, pendingRequests[0].ID)
	assert.Equal(t, expectedPendingRequest.StudentID, pendingRequests[0].StudentID)
	assert.Equal(t, expectedPendingRequest.TutorialID, pendingRequests[0].TutorialID)
	assert.Equal(t, "pending", pendingRequests[0].Status)

	// Cleanup
	database.DB.Unscoped().Delete(&models.Request{})
}
