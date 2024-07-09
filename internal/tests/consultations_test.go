package tests

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// var testTeachingAssistant = models.TeachingAssistant{
// 	Name:     "test_ta",
// 	Email:    "test_ta@gmail.com",
// 	Password: "test_ta",
// }

// var testStudent = models.Student{
// 	Name:     "test_student",
// 	Email:    "test_student@gmail.com",
// 	Password: "test_student",
// 	Modules:  []string{"test_CS1101S"},
// }

// var testTutorial = models.Tutorial{
// 	TutorialCode: "123456",
// 	Module:       "test_CS1101S",
// }

// var testStudents = []models.Student{
// 	{
// 		Name:     "test_student1",
// 		Email:    "test_student1@gmail.com",
// 		Password: "test_student1",
// 		Modules:  []string{"test_CS1101S"},
// 	},
// 	{
// 		Name:     "test_student2",
// 		Email:    "test_student2@gmail.com",
// 		Password: "test_student2",
// 		Modules:  []string{"test_CS1101S"},
// 	},
// 	{
// 		Name:     "test_student3",
// 		Email:    "test_student3@gmail.com",
// 		Password: "test_student3",
// 		Modules:  []string{"test_CS1101S"},
// 	},
// }

// Asserts whether the two consultations response are equal by comparing their fields
func assertEqualConsultationResponse(t *testing.T, expected *api.ConsultationResponse, actual *api.ConsultationResponse) {
	assert.Equal(t, expected.Tutorial.ID, actual.Tutorial.ID)
	assert.Equal(t, expected.TeachingAssistant.ID, actual.TeachingAssistant.ID)
	assert.Equal(t, expected.Student.ID, actual.Student.ID)

	assert.Equal(t, expected.Date, actual.Date)
	assert.Equal(t, expected.StartTime, actual.StartTime)
	assert.Equal(t, expected.EndTime, actual.EndTime)
	assert.Equal(t, expected.Booked, actual.Booked)
}

// Asserts whether the two booked consultations response are equal by comparing their fields
func assertEqualBookedConsultationsResponse(t *testing.T, expected *api.BookedConsultationsResponse, actual *api.BookedConsultationsResponse) {
	assert.Equal(t, len(expected.BookedConsultations), len(actual.BookedConsultations))
	for i, expectedBookedConsultationsByDate := range expected.BookedConsultations {
		actualBookedConsultationsByDate := actual.BookedConsultations[i]
		assert.Equal(t, expectedBookedConsultationsByDate.Date, actualBookedConsultationsByDate.Date)
		assert.Equal(t, len(expectedBookedConsultationsByDate.Consultations), len(actualBookedConsultationsByDate.Consultations))
		for j, expectedConsultation := range expectedBookedConsultationsByDate.Consultations {
			assertEqualConsultationResponse(t, &expectedConsultation, &actualBookedConsultationsByDate.Consultations[j])
		}
	}
}

// Test valid consultations fetch for date
func TestValidGetConsultations(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get consulations for the tutorial on the date
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.ConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

	// Get the tutorial
	tutorial, err := dataaccess.GetTutorialById(int(testTutorial.ID))
	assert.NoError(t, err)

	// Get the teaching assistant
	teachingAssistant, err := dataaccess.GetTeachingAssistantById(int(testTutorial.TeachingAssistantID))
	assert.NoError(t, err)

	var student models.Student
	// Compare expected consultations that should be fetched with the actual consultations on the date
	expectedConsultations := &[]api.ConsultationResponse{
		{
			ID:                1,
			Tutorial:          *tutorial,
			Student:           student,
			TeachingAssistant: *teachingAssistant,
			Date:              date,
			StartTime:         "10:00",
			EndTime:           "11:00",
			Booked:            false,
		},
		{
			ID:                2,
			Tutorial:          *tutorial,
			Student:           student,
			TeachingAssistant: *teachingAssistant,
			Date:              date,
			StartTime:         "11:00",
			EndTime:           "12:00",
			Booked:            false,
		},
	}

	assert.Equal(t, len(*expectedConsultations), len(actualConsultationsForTutorialForDate))
	for i, expectedConsultation := range *expectedConsultations {
		assertEqualConsultationResponse(t, &expectedConsultation, &actualConsultationsForTutorialForDate[i])
	}

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	// Clear consultations in actualConsultationsForTutorialForDate
	for _, consultation := range actualConsultationsForTutorialForDate {
		dataaccess.DeleteConsultationById(int(consultation.ID))
	}
}

// Test invalid date consultations fetch for date
func TestInvalidDateGetConsultations(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Invalid date
	date := "01-01-2021"

	// Send a request to get consulations for the tutorial on the date
	_, status, _ := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.Equal(t, http.StatusBadRequest, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test valid book consultation for student
func TestValidBookConsultation(t *testing.T) {
	// Test for 10 dates in the future
	for i := 1; i <= 10; i++ {
		// Create test TeachingAssistant, Student, Tutorial
		testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)
		// Get date string for day after tomorrow
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")

		// Send a request to get consulations for the tutorial on the date
		res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get response in json
		var response api.Response
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)
		resData, _ := json.Marshal(response.Data)

		// Get actual consultations for the tutorial on the date
		var consultationsForTutorialForDate api.ConsultationsResponse
		err = json.Unmarshal(resData, &consultationsForTutorialForDate)
		assert.NoError(t, err)
		actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

		for _, consultation := range actualConsultationsForTutorialForDate {
			// Send a request to the book consultation
			_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%d?userId=%d", int(testTutorial.ID),
				int(consultation.ID), int(testStudent.ID)), "PUT", testStudent)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			consultation, err := dataaccess.GetConsultationById(int(consultation.ID))
			assert.NoError(t, err)
			assert.True(t, consultation.Booked)
		}

		// Clean up
		CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
		// Clear consultations in actualConsultationsForTutorialForDate
		for _, consultation := range actualConsultationsForTutorialForDate {
			dataaccess.DeleteConsultationById(int(consultation.ID))
		}
	}
}

// Test invalid consultation id book consultation for student
// Test valid book consultation for student
func TestInvalidConsultationIdBookConsultation(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get consulations for the tutorial on the date
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.ConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

	for _, _ = range actualConsultationsForTutorialForDate {
		// Send a request to the book consultation
		_, status, _ = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%s?userId=%d", int(testTutorial.ID),
			"notintid", int(testStudent.ID)), "PUT", testStudent)
		assert.Equal(t, http.StatusBadRequest, status)
	}

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	// Clear consultations in actualConsultationsForTutorialForDate
	for _, consultation := range actualConsultationsForTutorialForDate {
		dataaccess.DeleteConsultationById(int(consultation.ID))
	}
}

// Test invalid user id book consultation for student
// Test valid book consultation for student
func TestInvalidUserIdBookConsultation(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get consulations for the tutorial on the date
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.ConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

	for _, consultation := range actualConsultationsForTutorialForDate {
		// Send a request to the book consultation
		_, status, _ = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%d?userId=%s", int(testTutorial.ID),
			int(consultation.ID), "invaliduserid"), "PUT", testStudent)
		assert.Equal(t, http.StatusBadRequest, status)
	}

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	// Clear consultations in actualConsultationsForTutorialForDate
	for _, consultation := range actualConsultationsForTutorialForDate {
		dataaccess.DeleteConsultationById(int(consultation.ID))
	}
}

// Test valid cancel consultation for student
func TestValidCancelConsultation(t *testing.T) {
	// Test for 10 dates in the future
	for i := 1; i <= 5; i++ {
		// Create test TeachingAssistant, Student, Tutorial
		testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
		assert.NoError(t, err)
		// Get date string for day after tomorrow
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")

		// Send a request to get consulations for the tutorial on the date
		res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get response in json
		var response api.Response
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)
		resData, _ := json.Marshal(response.Data)

		// Get actual consultations for the tutorial on the date
		var consultationsForTutorialForDate api.ConsultationsResponse
		err = json.Unmarshal(resData, &consultationsForTutorialForDate)
		assert.NoError(t, err)
		actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

		for _, consultation := range actualConsultationsForTutorialForDate {
			// Send a request to the book consultation
			_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%d?userId=%d", int(testTutorial.ID),
				int(consultation.ID), int(testStudent.ID)), "PUT", testStudent)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			consultation, err := dataaccess.GetConsultationById(int(consultation.ID))
			assert.NoError(t, err)
			assert.True(t, consultation.Booked)

			// Send a request to the cancel consultation
			_, status, err = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/cancel/%d?userId=%d", int(testTutorial.ID),
				int(consultation.ID), int(testStudent.ID)), "PUT", testStudent)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			consultation, err = dataaccess.GetConsultationById(int(consultation.ID))
			assert.NoError(t, err)
			assert.False(t, consultation.Booked)
		}

		// Clean up
		CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
		// Clear consultations in actualConsultationsForTutorialForDate
		for _, consultation := range actualConsultationsForTutorialForDate {
			dataaccess.DeleteConsultationById(int(consultation.ID))
		}
	}
}

// Test invalid consultation id cancel consultation for student
func TestInvalidConsultationIdCancelConsultation(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get consulations for the tutorial on the date
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.ConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

	for _, _ = range actualConsultationsForTutorialForDate {
		// Send a request to the cancel consultation
		_, status, _ = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/cancel/%s?userId=%d", int(testTutorial.ID),
			"notintid", int(testStudent.ID)), "PUT", testStudent)
		assert.Equal(t, http.StatusBadRequest, status)
	}

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	// Clear consultations in actualConsultationsForTutorialForDate
	for _, consultation := range actualConsultationsForTutorialForDate {
		dataaccess.DeleteConsultationById(int(consultation.ID))
	}
}

// Test invalid user id cancel consultation for student
func TestInvalidUserIdCancelConsultation(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get consulations for the tutorial on the date
	res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.ConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

	for _, consultation := range actualConsultationsForTutorialForDate {
		// Send a request to the cancel consultation
		_, status, _ = CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/cancel/%d?userId=%s", int(testTutorial.ID),
			int(consultation.ID), "invaliduserid"), "PUT", testStudent)
		assert.Equal(t, http.StatusBadRequest, status)
	}

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
	// Clear consultations in actualConsultationsForTutorialForDate
	for _, consultation := range actualConsultationsForTutorialForDate {
		dataaccess.DeleteConsultationById(int(consultation.ID))
	}
}

// Test get booked consultations for student
// Test by booking consultations for different students
// Each student should only see their booked consultations
func TestGetBookedConsultationsForStudent(t *testing.T) {
	consultationIds := []uint{}
	testStudentModels := []models.Student{}
	dates := []string{}
	testDefaultStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
	// Book 2 consultation slots for each student in testStudents
	for i, student := range testStudents {
		// Create test TeachingAssistant, Student, Tutorial
		student, err := CreateMockStudent(&student, testTeachingAssistant, testTutorial)
		assert.NoError(t, err)
		testStudentModels = append(testStudentModels, *student)
		date := time.Now().AddDate(0, 0, i+1).Format("2006-01-02")
		dates = append(dates, date)

		// Send a request to get consulations for the tutorial on the date
		res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", student)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get response in json
		var response api.Response
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)
		resData, _ := json.Marshal(response.Data)

		// Get actual consultations for the tutorial on the date
		var consultationsForTutorialForDate api.ConsultationsResponse
		err = json.Unmarshal(resData, &consultationsForTutorialForDate)
		assert.NoError(t, err)
		actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

		for _, consultation := range actualConsultationsForTutorialForDate {
			// Add consultation ID to consultationIds
			consultationIds = append(consultationIds, consultation.ID)
			// Send a request to the book consultation
			_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%d?userId=%d", int(testTutorial.ID),
				int(consultation.ID), int(student.ID)), "PUT", student)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			consultation, err := dataaccess.GetConsultationById(int(consultation.ID))
			assert.NoError(t, err)
			assert.True(t, consultation.Booked)
		}
	}

	// Compare expected booked consultations for each student with the actual booked consultations
	for i, student := range testStudentModels {
		res, status, err := CreateStudentAuthenticatedMockRequest(nil,
			fmt.Sprintf("/api/consultations/student/%d/%d?date=%s&time=%s", int(testTutorial.ID), int(student.ID), dates[i], "10:00"), "GET", &student)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		// Get response in json
		var response api.Response
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)
		resData, _ := json.Marshal(response.Data)

		// Get actual consultations for the tutorial on the date
		var consultationsForTutorialForDate api.BookedConsultationsResponse
		err = json.Unmarshal(resData, &consultationsForTutorialForDate)
		assert.NoError(t, err)
		actualConsultationsForTutorialForDate := consultationsForTutorialForDate

		expectedBookedConsultations := api.BookedConsultationsResponse{
			BookedConsultations: []api.BookedConsultationsByDate{
				{
					Date: dates[i],
					Consultations: []api.ConsultationResponse{
						{
							ID:                1,
							Tutorial:          *testTutorial,
							Student:           student,
							TeachingAssistant: *testTeachingAssistant,
							Date:              dates[i],
							StartTime:         "10:00",
							EndTime:           "11:00",
							Booked:            true,
						},
						{
							ID:                2,
							Tutorial:          *testTutorial,
							Student:           student,
							TeachingAssistant: *testTeachingAssistant,
							Date:              dates[i],
							StartTime:         "11:00",
							EndTime:           "12:00",
							Booked:            true,
						},
					},
				},
			},
		}

		assertEqualBookedConsultationsResponse(t, &expectedBookedConsultations, &actualConsultationsForTutorialForDate)
	}

	// Clean up
	// Clean up consultations
	for _, consultationId := range consultationIds {
		dataaccess.DeleteConsultationById(int(consultationId))
	}

	// Clean up students, ta, tutorial
	for _, student := range testStudentModels {
		CleanupCreatedStudent(&student)
	}
	CleanupCreatedStudent(testDefaultStudent)
	CleanupCreatedTeachingAssistant(testTeachingAssistant)
	CleanupCreatedTutorial(testTutorial)
}

// Test invalid date get booked consultations for student
func TestInvalidDateGetBookedConsultationsForStudent(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Send a request to get booked consulations for the tutorial on the date
	_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/student/%d/%d?date=%s&time=%s", int(testTutorial.ID), int(testStudent.ID), "01-01-2021", "10:00"), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test invalid time get booked consultations for student
func TestInvalidTimeGetBookedConsultationsForStudent(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	testStudent, _, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get booked consulations for the tutorial on the date
	_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/student/%d/%d?date=%s&time=%s", int(testTutorial.ID), int(testStudent.ID), date, "invalidtime"), "GET", testStudent)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test get booked consultations for teaching assistant
// Test by booking consultations for different students
// TA should see all booked consultations booked by all students
func TestGetBookedConsultationsForTeachingAssistant(t *testing.T) {
	consultationIds := []uint{}
	testStudentModels := []models.Student{}
	dates := []string{}
	testDefaultStudent, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)
	// Book 2 consultation slots for each student in testStudents
	for i, student := range testStudents {
		// Create test TeachingAssistant, Student, Tutorial
		student, err := CreateMockStudent(&student, testTeachingAssistant, testTutorial)
		assert.NoError(t, err)
		testStudentModels = append(testStudentModels, *student)
		date := time.Now().AddDate(0, 0, i+1).Format("2006-01-02")
		dates = append(dates, date)

		// Send a request to get consulations for the tutorial on the date
		res, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d?date=%s", int(testTutorial.ID), date), "GET", student)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)

		// Get response in json
		var response api.Response
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)
		resData, _ := json.Marshal(response.Data)

		// Get actual consultations for the tutorial on the date
		var consultationsForTutorialForDate api.ConsultationsResponse
		err = json.Unmarshal(resData, &consultationsForTutorialForDate)
		assert.NoError(t, err)
		actualConsultationsForTutorialForDate := consultationsForTutorialForDate.Consultations

		for _, consultation := range actualConsultationsForTutorialForDate {
			// Add consultation ID to consultationIds
			consultationIds = append(consultationIds, consultation.ID)
			// Send a request to the book consultation
			_, status, err := CreateStudentAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/%d/book/%d?userId=%d", int(testTutorial.ID),
				int(consultation.ID), int(student.ID)), "PUT", student)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			consultation, err := dataaccess.GetConsultationById(int(consultation.ID))
			assert.NoError(t, err)
			assert.True(t, consultation.Booked)
		}
	}

	// Compare expected booked consultations for TA with the actual booked consultations
	res, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil,
		fmt.Sprintf("/api/consultations/teachingAssistant/%d?date=%s&time=%s", int(testTutorial.ID), dates[0], "10:00"), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	// Get response in json
	var response api.Response
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)
	resData, _ := json.Marshal(response.Data)

	// Get actual consultations for the tutorial on the date
	var consultationsForTutorialForDate api.BookedConsultationsResponse
	err = json.Unmarshal(resData, &consultationsForTutorialForDate)
	assert.NoError(t, err)
	actualConsultationsForTutorialForDate := consultationsForTutorialForDate

	expectedBookedConsultations := api.BookedConsultationsResponse{
		BookedConsultations: []api.BookedConsultationsByDate{
			{
				Date: dates[0],
				Consultations: []api.ConsultationResponse{
					{
						ID:                1,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[0],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[0],
						StartTime:         "10:00",
						EndTime:           "11:00",
						Booked:            true,
					},
					{
						ID:                2,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[0],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[0],
						StartTime:         "11:00",
						EndTime:           "12:00",
						Booked:            true,
					},
				},
			},
			{
				Date: dates[1],
				Consultations: []api.ConsultationResponse{
					{
						ID:                1,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[1],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[1],
						StartTime:         "10:00",
						EndTime:           "11:00",
						Booked:            true,
					},
					{
						ID:                2,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[1],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[1],
						StartTime:         "11:00",
						EndTime:           "12:00",
						Booked:            true,
					},
				},
			},
			{
				Date: dates[2],
				Consultations: []api.ConsultationResponse{
					{
						ID:                1,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[2],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[2],
						StartTime:         "10:00",
						EndTime:           "11:00",
						Booked:            true,
					},
					{
						ID:                2,
						Tutorial:          *testTutorial,
						Student:           testStudentModels[2],
						TeachingAssistant: *testTeachingAssistant,
						Date:              dates[2],
						StartTime:         "11:00",
						EndTime:           "12:00",
						Booked:            true,
					},
				},
			},
		},
	}

	assertEqualBookedConsultationsResponse(t, &expectedBookedConsultations, &actualConsultationsForTutorialForDate)

	// Clean up
	// Clean up consultations
	for _, consultationId := range consultationIds {
		dataaccess.DeleteConsultationById(int(consultationId))
	}

	// Clean up students, ta, tutorial
	for _, student := range testStudentModels {
		CleanupCreatedStudent(&student)
	}
	CleanupCreatedStudent(testDefaultStudent)
	CleanupCreatedTeachingAssistant(testTeachingAssistant)
	CleanupCreatedTutorial(testTutorial)
}

// Test invalid date get booked consultations for teaching assistant
func TestInvalidDateGetBookedConsultationsForTeachingAssistant(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Send a request to get booked consulations for the tutorial on the date
	_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/teachingAssistant/%d?date=%s&time=%s", int(testTutorial.ID), "01-01-2021", "10:00"), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}

// Test invalid time get booked consultations for teaching assistant
func TestInvalidTimeGetBookedConsultationsForTeachingAssistant(t *testing.T) {
	// Create test TeachingAssistant, Student, Tutorial
	_, testTeachingAssistant, testTutorial, err := CreateSingleMockStudentTeachingAssistantAndTutorial()
	assert.NoError(t, err)

	// Get date string for day after tomorrow
	date := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	// Send a request to get booked consulations for the tutorial on the date
	_, status, err := CreateTeachingAssistantAuthenticatedMockRequest(nil, fmt.Sprintf("/api/consultations/teachingAssistant/%d?date=%s&time=%s", int(testTutorial.ID), date, "invalidtime"), "GET", testTeachingAssistant)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// Clean up
	CleanupSingleCreatedStudentTeachingAssistantAndTutorial()
}
