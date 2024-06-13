package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateMessage(discussionId int, senderId int, userType string, content string) error {
	message := &models.Message{DiscussionID: discussionId, SenderID: senderId, UserType: userType, Content: content}
	result := database.DB.Table("messages").Create(message)
	return result.Error
}

func GetMessagesByDiscussionId(id int) (*[]models.Message, error) {
	var messages []models.Message
	result := database.DB.Table("messages").Where("discussion_id = ?", id).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return &messages, nil
}

func GetMessagesByTutorialId(id int) (*[]models.Message, error) {
	discussionId, err := GetDiscussionIdByTutorialId(id)
	if err != nil {
		return nil, err
	}

	messages, err := GetMessagesByDiscussionId(discussionId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}