package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var filePayload1 = validFilesTests[1]
var filePayload2 = validFilesTests[2]

func TestFilesDataaccess(t *testing.T) {
	var tutorialFileId int
	expectedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: filePayload1.Filepath, Name: filePayload1.Name, Visible: true, Week: filePayload1.Week}
	expectedPrivatedTutorialFile := &models.TutorialFile{TutorialID: int(testTutorial.ID), Filepath: filePayload1.Filepath, Name: filePayload1.Name, Visible: false, Week: filePayload1.Week}

	t.Run("Create Tutorial File", func(t *testing.T) {
		// Current no. of files in the test db should be 0
		var count int64
		database.DB.Table("tutorial_files").Count(&count)
		assert.Equal(t, 0, int(count))

		// Create a tutorial file
		err := dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload1.Name, filePayload1.Week, filePayload1.Filepath)
		assert.NoError(t, err)

		// Current no. of files in the test db should be 1
		database.DB.Table("tutorial_files").Count(&count)
		assert.Equal(t, 1, int(count))
	})

	t.Run("Get Tutorial File from TutorialID and Filename", func(t *testing.T) {
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileFromTutorialIDAndFilename(int(testTutorial.ID), filePayload1.Name, filePayload1.Week)
		assert.NoError(t, err)
		tutorialFileId = int(tutorialFile.ID)

		// Compare expected tutorial file that should be created with the actual file created
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	})

	// Create one more file
	err := dataaccess.CreateTutorialFile(int(testTutorial.ID), filePayload2.Name, filePayload2.Week, filePayload2.Filepath)
	assert.NoError(t, err)

	t.Run("Get all Tutorial Files from TutorialID and Week", func(t *testing.T) {
		// Get all tutorial files for that tutorial and week
		actualTutorialFiles, err := dataaccess.GetAllTutorialFilesFromTutorialIDAndWeek(int(testTutorial.ID), filePayload1.Week)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		expectedTutorialFiles := &[]models.TutorialFile{
			{TutorialID: int(testTutorial.ID), Filepath: filePayload1.Filepath, Name: filePayload1.Name, Visible: true, Week: filePayload1.Week},
			{TutorialID: int(testTutorial.ID), Filepath: filePayload2.Filepath, Name: filePayload2.Name, Visible: true, Week: filePayload2.Week},
		}

		assert.Equal(t, len(*expectedTutorialFiles), len(*actualTutorialFiles))
		for j, expectedTutorialFile := range *expectedTutorialFiles {
			assertEqualTutorialFile(t, &expectedTutorialFile, &(*actualTutorialFiles)[j])
		}
	})

	t.Run("Get Tutorial File by ID", func(t *testing.T) {
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileById(tutorialFileId)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	})

	t.Run("Get Tutorial File by Filepath", func(t *testing.T) {
		// Get the actual tutorial file that is created
		tutorialFile, err := dataaccess.GetTutorialFileByFilepath(filePayload1.Filepath)
		assert.NoError(t, err)

		// Compare expected tutorial file that should be created with the actual file created
		assertEqualTutorialFile(t, expectedTutorialFile, tutorialFile)
	})

	t.Run("Private Tutorial File by Filepath", func(t *testing.T) {
		// Private the first file
		err := dataaccess.PrivateFileByFilepath(filePayload1.Filepath)
		assert.NoError(t, err)

		// Get the privated file
		var tutorialFile models.TutorialFile
		database.DB.Table("tutorial_files").Where("filepath = ?", filePayload1.Filepath).Find(&tutorialFile)

		// Compare expected tutorial file with the actual file
		assertEqualTutorialFile(t, expectedPrivatedTutorialFile, &tutorialFile)
	})

	t.Run("Unprivate Tutorial File by Filepath", func(t *testing.T) {
		// Unprivate the first file
		err := dataaccess.UnprivateFileByFilepath(filePayload1.Filepath)
		assert.NoError(t, err)

		// Get the unprivated file
		var tutorialFile models.TutorialFile
		database.DB.Table("tutorial_files").Where("filepath = ?", filePayload1.Filepath).Find(&tutorialFile)

		// Compare expected tutorial file with the actual file
		assertEqualTutorialFile(t, expectedTutorialFile, &tutorialFile)
	})

	t.Run("Delete Tutorial File by Filepath", func(t *testing.T) {
		// Current no. of files in the test db should be 2
		var count int64
		database.DB.Table("tutorial_files").Count(&count)
		assert.Equal(t, 2, int(count))

		// Delete the tutorial files created
		err := dataaccess.DeleteTutorialFileByFilepath(filePayload1.Filepath)
		assert.NoError(t, err)
		err = dataaccess.DeleteTutorialFileByFilepath(filePayload2.Filepath)
		assert.NoError(t, err)

		// Current no. of files in the test db should be 0
		database.DB.Table("tutorial_files").Count(&count)
		assert.Equal(t, 0, int(count))
	})
}