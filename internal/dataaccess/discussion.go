package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateDiscussion(id int) error {
	discussion := &models.Discussion{TutorialID: id}
	result := database.DB.Table("discussions").Create(discussion)
	return result.Error
}

func CreateDiscussionForEveryTutorial() error {
	tutorialIds, err := GetAllTutorialIDs()
	if err != nil {
		return err
	}

	if tutorialIds == nil {
		return nil
	}

	for _, tutorialId := range *tutorialIds {
		err = CreateDiscussion(tutorialId)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDiscussionIdByTutorialId(id int) (int, error) {
	var discussionId int
	result := database.DB.Table("discussions").Where("tutorial_id = ?", id).First(discussionId)
	if result.Error != nil {
		return -1, result.Error
	}

	return discussionId, nil
}