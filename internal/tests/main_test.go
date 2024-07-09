package tests

import (
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/router"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var TestRouter = router.Setup()

func TestMain(m *testing.M) {
	err := database.Connect(true)
	if err != nil {
		log.Fatalln("Failed to connect to database!", err)
	}

	// err = util.Migrate()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// Initialise auth obj
	err = auth.InitialiseAuthObj()
	if err != nil {
		log.Fatalln("Failed to initialise auth obj!", err)
	}

	// Initialise auth obj
	err = auth.InitialiseAuthObj()
	if err != nil {
		log.Fatalln("Failed to initialise auth obj!", err)
	}

	os.Exit(m.Run())
}

func CreateMockRequest(payload interface{}, url string, method string, tokens ...string) ([]byte, int, error) {
	var req *http.Request

	// Check if given payload is nil
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, -1, err
		}

		requestBody := bytes.NewBuffer(payloadBytes)
		req = httptest.NewRequest(method, url, requestBody)
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	// Add access token to Authorization header if token is present
	req.Header.Set("Content-Type", "application/json")
	if len(tokens) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens[0]))
	}

	// Dispatches request to the correct handler
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)
	return w.Body.Bytes(), w.Code, nil
}

func CreateStudentAuthenticatedMockRequest(payload interface{}, url string, method string, student *models.Student) ([]byte, int, error) {
	// Create an authenticated Student user
	authUser := auth.AuthenticatedUser{
		ID:          int(student.ID),
		Name:        student.Name,
		Email:       student.Email,
		Role:        auth.RoleStudent,
	}
	
	// Generate access and refresh tokens
	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		return nil, -1, err
	}
	
	return CreateMockRequest(payload, url, method, tokens.AccessToken)
}

func CreateTeachingAssistantAuthenticatedMockRequest(payload interface{}, url string, method string, teachingAssistant *models.TeachingAssistant) ([]byte, int, error) {
	// Create an authenticated TeachingAssistant user
	authUser := auth.AuthenticatedUser{
		ID:          int(teachingAssistant.ID),
		Name:        teachingAssistant.Name,
		Email:       teachingAssistant.Email,
		Role:        auth.RoleTeachingAssistant,
	}
	
	// Generate access and refresh tokens
	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		return nil, -1, err
	}
	
	return CreateMockRequest(payload, url, method, tokens.AccessToken)
}

func CreateMockStudentTeachingAssistantAndTutorial() (*models.Student, *models.TeachingAssistant, *models.Tutorial, error) {
	// Create test TeachingAssistant, Student, Tutorial
	testTeachingAssistant, err := dataaccess.CreateTeachingAssistant(testTeachingAssistant.Name, testTeachingAssistant.Email, testTeachingAssistant.Password)
	if err != nil {
		return nil, nil, nil, err
	}

	testTutorial, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, int(testTeachingAssistant.ID))
	if err != nil {
		return nil, nil, nil, err
	}

	testStudent, err := dataaccess.CreateStudent(testStudent.Name, testStudent.Email, testStudent.Password, testStudent.Modules)
	if err != nil {
		return nil, nil, nil, err
	}

	// Assign TA to the tutorial
	testTeachingAssistant.TutorialID = int(testTutorial.ID)
	database.DB.Save(testTeachingAssistant)

	// Enable Student to join the tutorial
	err = dataaccess.JoinTutorial(int(testStudent.ID), int(testTutorial.ID))
	if err != nil {
		return nil, nil, nil, err
	}

	student, err := dataaccess.GetStudentByEmail(testStudent.Email)
	if err != nil {
		return nil, nil, nil, err
	}

	teachingAssistant, err := dataaccess.GetTeachingAssistantByEmail(testTeachingAssistant.Email)
	if err != nil {
		return nil, nil, nil, err
	}

	return student, teachingAssistant, testTutorial, nil
}

// Cleanup all the created Teaching Assistants, Tutorials and Students for tests
func CleanupCreatedStudentTeachingAssistantAndTutorial() {
	dataaccess.DeleteTeachingAssistantByEmail(testTeachingAssistant.Email)
	dataaccess.DeleteTutorialByClassAndModuleCode(testTutorial.TutorialCode, testTutorial.Module)
	dataaccess.DeleteStudentByEmail(testStudent.Email)
	// dataaccess.DeleteRegistryByStudentIDAndTutorialID(testTeachingAssistant.TutorialID, testStudent.ID)
}
