package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateRequest(studentId int, tutorialId int) error {
	request := &models.Request{StudentID: studentId, TutorialID: tutorialId, Status: "pending"}
	result := database.DB.Table("requests").Create(request)
	return result.Error
}

func GetPendingRequestsByTutorialId(id int) ([]*models.Request, error) {
	var requests []*models.Request
	result := database.DB.Table("requests").Where("tutorialId = ?", id).Where("status = ?", "pending").Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}

func GetRequestById(id int) (*models.Request, error) {
	var request models.Request
	result := database.DB.Table("requests").Where("id = ?", id).First(&request)
	if result.Error != nil {
		return nil, result.Error
	}

	return &request, nil
}

func AcceptRequestById(id int) error {
	request, err := GetRequestById(id)
	if err != nil {
		return err
	}

	request.Status = "accepted"
	database.DB.Save(&request)
	return nil
}

func RejectRequestById(id int) error {
	request, err := GetRequestById(id)
	if err != nil {
		return err
	}

	request.Status = "rejected"
	database.DB.Save(&request)
	return nil
}
