package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDiscussion = models.Discussion{TutorialID: 0}

func TestDiscussionDataaccess(t *testing.T) {
	t.Run("Create Discussion", func(t *testing.T) {
		// Current no. of discussions in the test db should be 0
		var count int64
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a discussion
		err := dataaccess.CreateDiscussion(testDiscussion.TutorialID)
		assert.NoError(t, err)

		// Current no. of discussions in the test db should be 1
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	// Get id of test discussion created
	var testDiscussionID int64
	result := database.DB.Table("discussions").Where("tutorial_id = ?", testDiscussion.TutorialID).Select("id").Find(&testDiscussionID)
	assert.NoError(t, result.Error)

	// Create two tutorials
	for i := 0; i < 2; i++ {
		_, err := dataaccess.CreateTutorial(testTutorial.TutorialCode, testTutorial.Module, int(testTeachingAssistant.ID))
		assert.NoError(t, err)
	}

	t.Run("Create Discussion for every Tutorial", func(t *testing.T) {
		// Current no. of discussions in the test db should be 1
		var count int64
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 1, int(count))

		// Create discussions for every tutorial - Should create 2 more discussions since there are 2 tutorials in the test db
		err := dataaccess.CreateDiscussionForEveryTutorial()
		assert.NoError(t, err)

		// Current no. of discussions in the test db should be 3
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 3, int(count))
	})

	t.Run("Get Discussion ID by Tutorial ID", func(t *testing.T) {
		// Get the discussion id
		discussionId, err := dataaccess.GetDiscussionIdByTutorialId(testDiscussion.TutorialID)
		assert.NoError(t, err)

		// Compare expected discussion id that should be fetched with the actual discussion id fetched
		assert.Equal(t, int(testDiscussionID), discussionId)
	})

	t.Run("Get Discussion by ID", func(t *testing.T) {
		// Get the discussion
		discussion, err := dataaccess.GetDiscussionById(int(testDiscussionID))
		assert.NoError(t, err)

		// Compare expected discussion that should be fetched with the actual discussion fetched
		assert.Equal(t, testDiscussion.TutorialID, discussion.TutorialID)
	})

	t.Run("Delete Discussion by ID", func(t *testing.T) {
		// Current no. of discussions in the test db should be 3
		var count int64
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 3, int(count))

		// Delete the discussion created by the first test
		err := dataaccess.DeleteDiscussionById(int(testDiscussionID))
		assert.NoError(t, err)

		// Current no. of discussions in the test db should be 2
		database.DB.Table("discussions").Count(&count)
		assert.Equal(t, 2, int(count))
	})

	// Clean up
	tutorialIds, err := dataaccess.GetAllTutorialIds()
	assert.NoError(t, err)
	for _, tutorialId := range *tutorialIds {
		dataaccess.DeleteTutorialById(tutorialId)
		dicussionId, err := dataaccess.GetDiscussionIdByTutorialId(tutorialId)
		assert.NoError(t, err)
		dataaccess.DeleteDiscussionById(dicussionId)
	}
	dataaccess.DeleteDiscussionById(int(testDiscussionID))
}