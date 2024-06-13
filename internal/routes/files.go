package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func FileRoutes(r chi.Router) {
	r.Get("/student/{tutorialId}/{week}", handlers.GetAllTutorialFilesForStudents)
	r.Get("/teachingAssistant/{tutorialId}/{week}", handlers.GetAllTutorialFilesForTAs)
	r.Post("/upload/{tutorialId}", handlers.UploadFilepath)
	r.Patch("/delete/{tutorialId}", handlers.DeleteFilepath)
	r.Patch("/private/{tutorialId}", handlers.PrivateFile)
	r.Patch("/unprivate/{tutorialId}", handlers.UnprivateFile)
}