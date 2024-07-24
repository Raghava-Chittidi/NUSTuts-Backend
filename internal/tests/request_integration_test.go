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

func TestInvalidRequestFormatToJoinTutorial(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty, arbitray where
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	invalidFormatPayload := struct {
		InvalidField string `json:"invalid_field"`
	}{
		InvalidField: "invalid",
	}

	// Create a request to join tutorial with invalid payload
	_, status, _ := CreateStudentAuthenticatedMockRequest(invalidFormatPayload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.Equal(t, http.StatusBadRequest, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestNonExistingModuleRequestToJoinTutorial(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty, arbitray where
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create a request to join tutorial with non-existing module code
	requestToJoinTutPayload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "1",
	}
	_, status, _ := CreateStudentAuthenticatedMockRequest(requestToJoinTutPayload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.Equal(t, http.StatusInternalServerError, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestNonExistingClassNoRequestToJoinTutorial(t *testing.T) {
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
	_, err = dataaccess.GetTutorialByClassAndModuleCode("1", "CS2040S")
	assert.NoError(t, err)

	// Create a request to join tutorial with non-existing class no
	requestToJoinTutPayload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "2",
	}
	_, status, _ := CreateStudentAuthenticatedMockRequest(requestToJoinTutPayload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.Equal(t, http.StatusInternalServerError, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestEmptyRequestToJoinTutorial(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty, arbitray where
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create a request to join tutorial with empty payload
	_, status, _ := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.Equal(t, http.StatusBadRequest, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
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

func TestAcceptNonExistingRequestID(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Accept a non-existing request
	_, status, _ := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d/accept", 1), "PATCH", testStudent)
	assert.Equal(t, http.StatusInternalServerError, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestValidRejectRequest(t *testing.T) {
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

	// Reject the request
	_, status, err = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d/reject", requests[0].ID), "PATCH", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get the request
	request, err := dataaccess.GetRequestById(int(requests[0].ID))
	assert.NoError(t, err)

	// Assert that the request is rejected
	assert.Equal(t, "rejected", request.Status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	// Delete the test tutorial
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Tutorial{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestRejectNonExistingRequestID(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Reject a non-existing request
	_, status, _ := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d/reject", 1), "PATCH", testStudent)
	assert.Equal(t, http.StatusInternalServerError, status)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

func TestValidGetPendingRequests(t *testing.T) {
	testStudentModels := []models.Student{}
	testDefaultStudent, testDefaultTeachingAssistant, testDefaultTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
	for _, student := range testStudents {
		// Create test TeachingAssistant, Student, Tutorial
		student, err := CreateMockStudent(&student, testDefaultTeachingAssistant, testDefaultTutorial)
		assert.NoError(t, err)
		testStudentModels = append(testStudentModels, *student)
	}

	// Make sure requests table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create new TA
	testTeachingAssistant, err := dataaccess.CreateTeachingAssistant("TEST", "TEST@gmail.com", "1234567890")
	assert.NoError(t, err)

	// Create test tutorial
	testTutorial, err := dataaccess.CreateTutorial("1", "CS2040S", int(testTeachingAssistant.ID))
	assert.NoError(t, err)
	// Get the test tutorial
	_, err = dataaccess.GetTutorialByClassAndModuleCode("1", "CS2040S")
	assert.NoError(t, err)

	for _, student := range testStudentModels {
		// Create a request to join tutorial
		requestToJoinTutPayload := api.RequestToJoinTutorialPayload{
			StudentID:  int(student.ID),
			ModuleCode: "CS2040S",
			ClassNo:    "1",
		}
		_, status, err := CreateStudentAuthenticatedMockRequest(requestToJoinTutPayload, fmt.Sprintf("/api/requests/"), "POST", &student)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
	}

	// Current no. of requests in the test db should be number of students
	// Get all requests in the test db
	var requests []*models.Request
	result := database.DB.Table("requests").Find(&requests)
	assert.NoError(t, result.Error)
	// Assert that request count is 1
	assert.Equal(t, len(testStudentModels), len(requests))

	// Get all pending requests
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d", int(testTutorial.ID)), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual attendance string for the tutorial
	var requestsResponse []api.RequestResponse
	err = json.Unmarshal(resData, &requestsResponse)
	assert.NoError(t, err)

	// Assert that request response length is equal to the number of students
	assert.Equal(t, len(testStudentModels), len(requestsResponse))
	// Assert that the request response contains all the student requests and the correct status
	for i, student := range testStudentModels {
		assert.Equal(t, student.Name, requestsResponse[i].Name)
		assert.Equal(t, student.Email, requestsResponse[i].Email)
	}

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	// Delete the test tutorial
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Tutorial{})

	// Clean up students, ta, tutorial
	for _, student := range testStudentModels {
		CleanupCreatedStudent(&student)
	}
	CleanupCreatedStudent(testDefaultStudent)
	CleanupCreatedTeachingAssistant(testDefaultTeachingAssistant)
	CleanupCreatedTeachingAssistant(testTeachingAssistant)
	CleanupCreatedTutorial(testDefaultTutorial)
}

func TestValidUnrequestedTutorialClassNo(t *testing.T) {
	testStudent, _, _, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Make sure requests table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})

	// Current no. of requests in the test db should be 0
	var count int64
	database.DB.Table("requests").Count(&count)
	assert.Equal(t, 0, int(count))

	// Create test tutorial classes
	_, err = dataaccess.CreateTutorial("1", "CS2040S", 50)
	assert.NoError(t, err)
	_, err = dataaccess.CreateTutorial("2", "CS2040S", 100)
	assert.NoError(t, err)
	_, err = dataaccess.CreateTutorial("3", "CS2040S", 101)
	assert.NoError(t, err)

	// Create a request to join tutorial
	requestToJoinTut1Payload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "2",
	}
	requestToJoinTut2Payload := api.RequestToJoinTutorialPayload{
		StudentID:  int(testStudent.ID),
		ModuleCode: "CS2040S",
		ClassNo:    "3",
	}
	_, status, err := CreateStudentAuthenticatedMockRequest(requestToJoinTut1Payload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	_, status, err = CreateStudentAuthenticatedMockRequest(requestToJoinTut2Payload, fmt.Sprintf("/api/requests/"), "POST", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)

	// Get unrequested tutorial classes
	classesRes, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/requests/%d/CS2040S", int(testStudent.ID)), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(classesRes, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual unrequested classes for the tutorial
	var unrequestedClasses []string
	err = json.Unmarshal(resData, &unrequestedClasses)
	assert.NoError(t, err)

	// Log the unrequested classes
	fmt.Println(unrequestedClasses)

	// Assert that unrequested classes length is equal to 2
	assert.Equal(t, 2, len(unrequestedClasses))
	// Assert that the unrequested classes contains the correct class no
	assert.Contains(t, unrequestedClasses, "2")
	assert.Contains(t, unrequestedClasses, "3")

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Request{})
	// Delete the test tutorial
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Tutorial{})
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}
