package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var messagePayload1 = messageTests[0]
var messagePayload2 = messageTests[1]

func TestMessagesDataaccess(t *testing.T) {
	err := dataaccess.CreateDiscussion(int(testTutorial.ID))
	assert.NoError(t, err)
	testDiscussionID, err := dataaccess.GetDiscussionIdByTutorialId(int(testTutorial.ID))
	assert.NoError(t, err)
	t.Run("Create Message", func(t *testing.T) {
		// Current no. of messages in the test db should be 0
		var count int64
		database.DB.Table("messages").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a message
		err := dataaccess.CreateMessage(testDiscussionID, messagePayload1.SenderID, messagePayload1.UserType, messagePayload1.Content)
		assert.NoError(t, err)

		// Current no. of messages in the test db should be 1
		database.DB.Table("messages").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	// Create one more message
	err = dataaccess.CreateMessage(testDiscussionID, messagePayload2.SenderID, messagePayload2.UserType, messagePayload2.Content)
	assert.NoError(t, err)

	t.Run("Get Messages by Discussion ID", func(t *testing.T) {
		// Get the messages
		messages, err := dataaccess.GetMessagesByDiscussionId(testDiscussionID)
		assert.NoError(t, err)

		// Current no. of fetched messages should be 2
		assert.Equal(t, 2, len(*messages))

		// Compare expected messages that should be fetched with the actual messages fetched
		for i, message := range *messages {
			expectedMessage := &models.Message{DiscussionID: testDiscussionID, SenderID: messageTests[i].SenderID, UserType: messageTests[i].UserType, Content: messageTests[i].Content}
			assertEqualMessage(t, expectedMessage, &message)
		}
	})

	t.Run("Get Messages by Tutorial ID", func(t *testing.T) {
		// Get the messages
		messages, err := dataaccess.GetMessagesByTutorialId(int(testTutorial.ID))
		assert.NoError(t, err)

		// Current no. of fetched messages should be 2
		assert.Equal(t, 2, len(*messages))

		// Compare expected messages that should be fetched with the actual messages fetched
		for i, message := range *messages {
			expectedMessage := &models.Message{DiscussionID: testDiscussionID, SenderID: messageTests[i].SenderID, UserType: messageTests[i].UserType, Content: messageTests[i].Content}
			assertEqualMessage(t, expectedMessage, &message)
		}
	})

	t.Run("Delete Messages by Discussion ID", func(t *testing.T) {
		// Current no. of messages in the test db should be 2
		var count int64
		database.DB.Table("messages").Count(&count)
		assert.Equal(t, 2, int(count))

		// Delete the messages created
		err := dataaccess.DeleteMessagesByDiscussionId(testDiscussionID)
		assert.NoError(t, err)

		// Current no. of messages in the test db should be 0
		database.DB.Table("messages").Count(&count)
		assert.Equal(t, 0, int(count))
	})

	// Clean up
	dataaccess.DeleteDiscussionById(testDiscussionID)
}