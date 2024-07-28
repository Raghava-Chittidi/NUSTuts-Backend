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

// Creates a dicussion for each tutorial present in the database
func CreateDiscussionForEveryTutorial() error {
	tutorialIds, err := GetAllTutorialIds()
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
	var discussion models.Discussion
	result := database.DB.Table("discussions").Where("tutorial_id = ?", id).First(&discussion)
	if result.Error != nil {
		return -1, result.Error
	}

	return int(discussion.ID), nil
}

func GetDiscussionById(id int) (*models.Discussion, error) {
	var discussion models.Discussion
	result := database.DB.Table("discussions").Where("id = ?", id).First(&discussion)
	if result.Error != nil {
		return nil, result.Error
	}

	return &discussion, nil
}

func DeleteDiscussionById(id int) (error) {
	discussion, err := GetDiscussionById(id)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("discussions").Delete(&discussion)
	if result.Error != nil {
		return result.Error
	}

	return nil
}