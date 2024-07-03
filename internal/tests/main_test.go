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
	"net/http/httptest"
	"os"
	"testing"
)

var TestRouter = router.TestSetup()

func TestMain(m *testing.M) {
	err := database.Connect(true)
	if err != nil {
		log.Fatalln(err)
	}

	// err = util.Migrate()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	os.Exit(m.Run())
}

func CreateMockRequest(payload interface{}, url string, method string, tokens ...string) ([]byte, int, error) {
	var requestBody *bytes.Buffer = nil
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, -1, err
		}

		requestBody = bytes.NewBuffer(payloadBytes)
	}

	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if len(tokens) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens[0]))
	}

	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)
	return w.Body.Bytes(), w.Code, nil
}

func CreateStudentAuthenticatedMockRequest(payload interface{}, url string, method string, student *models.Student) ([]byte, int, error) {
	authUser := auth.AuthenticatedUser{
		ID:          int(student.ID),
		Name:        student.Name,
		Email:       student.Email,
		Role:        auth.RoleStudent,
	}
	
	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		return nil, -1, err
	}
	
	return CreateMockRequest(payload, url, method, tokens.AccessToken)
}

func CreateTeachingAssistantAuthenticatedMockRequest(payload interface{}, url string, method string, teachingAssistant *models.TeachingAssistant) ([]byte, int, error) {
	authUser := auth.AuthenticatedUser{
		ID:          int(teachingAssistant.ID),
		Name:        teachingAssistant.Name,
		Email:       teachingAssistant.Email,
		Role:        auth.RoleTeachingAssistant,
	}
	
	tokens, err := auth.AuthObj.GenerateTokens(&authUser)
	if err != nil {
		return nil, -1, err
	}
	
	return CreateMockRequest(payload, url, method, tokens.AccessToken)
}

func CreateMockTeachingAssistantAndMockTutorial() (*models.TeachingAssistant, *models.Tutorial, error) {
	testTeachingAssistant, err := dataaccess.CreateTeachingAssistant(testTeachingAssistant.Name, testTeachingAssistant.Email, testTeachingAssistant.Password)
	if err != nil {
		return nil, nil, err
	}

	testTutorial, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, int(testTeachingAssistant.ID))
	if err != nil {
		return nil, nil, err
	}

	testTeachingAssistant.TutorialID = int(testTutorial.ID)
	database.DB.Save(testTeachingAssistant)
	return testTeachingAssistant, testTutorial, nil
}
