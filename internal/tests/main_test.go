package tests

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/router"
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var TestRouter = router.Setup()

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

func CreateMockRequest(payload interface{}, url string, method string) ([]byte, int, error) {
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
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)
	return w.Body.Bytes(), w.Code, nil
}