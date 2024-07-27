package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

// Creates a request for a student to join a tutorial
// The request is initially set to pending
func CreateRequest(studentId int, tutorialId int) error {
	request := &models.Request{StudentID: studentId, TutorialID: tutorialId, Status: "pending"}
	result := database.DB.Table("requests").Create(request)
	return result.Error
}

// Gets all pending requests for a tutorial
func GetPendingRequestsByTutorialId(id int) ([]*models.Request, error) {
	var requests []*models.Request
	result := database.DB.Table("requests").Where("tutorial_id = ?", id).Where("status = ?", "pending").Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}

// Gets a request by its id
func GetRequestById(id int) (*models.Request, error) {
	var request models.Request
	result := database.DB.Table("requests").Where("id = ?", id).First(&request)
	if result.Error != nil {
		return nil, result.Error
	}

	return &request, nil
}

// Accepts a request by its id
// The request status is set to accepted
func AcceptRequestById(id int) error {
	request, err := GetRequestById(id)
	if err != nil {
		return err
	}

	request.Status = "accepted"
	database.DB.Save(&request)
	return nil
}

// Rejects a request by its id
// The request status is set to rejected
func RejectRequestById(id int) error {
	request, err := GetRequestById(id)
	if err != nil {
		return err
	}

	request.Status = "rejected"
	database.DB.Save(&request)
	return nil
}

// Gets all requests requested by a student for class no of a module
func GetClassNoByStudentIdAndModuleCode(id int, moduleCode string) (*[]string, error) {
	var classNoArr []string
	result := database.DB.Table("requests").Joins("JOIN tutorials ON requests.tutorial_id = tutorials.id").
		Where("requests.student_id = ?", id).Where("tutorials.module = ?", moduleCode).
		Select("tutorial_code").Find(&classNoArr)

	if result.Error != nil {
		return nil, result.Error
	}

	return &classNoArr, nil
}
