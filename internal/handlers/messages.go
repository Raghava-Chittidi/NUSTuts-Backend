package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

func GetAllMessagesForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	messages, err := dataaccess.GetMessagesByTutorialId(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	messagesWithSenders := lo.Map(*messages, func(item models.Message, index int) api.MessageResponse {
		if (item.UserType == "student") {
			sender, err := dataaccess.GetStudentById(item.SenderID)
			if err != nil {
				util.ErrorJSON(w, err, http.StatusInternalServerError)
				return api.MessageResponse{}
			}

			return api.MessageResponse{Sender: sender.Name, TutorialID: tutorialId, UserType: item.UserType, Content: item.Content}
		} else {
			sender, err := dataaccess.GetTeachingAssistantById(item.SenderID)
			if err != nil {
				util.ErrorJSON(w, err, http.StatusInternalServerError)
				return api.MessageResponse{}
			}

			return api.MessageResponse{Sender: sender.Name, TutorialID: tutorialId, UserType: item.UserType, Content: item.Content}
		}

	})
	
	res := api.MessagesResponse{Messages: messagesWithSenders}
	util.WriteJSON(w, api.Response{Message: "Fetched messages successfully!", Data: res}, http.StatusOK)
}

func CreateMessageForTutorial(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload api.CreateMesssagePayload
	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	discussionId, err := dataaccess.GetDiscussionIdByTutorialId(tutorialId)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = dataaccess.CreateMessage(discussionId, payload.SenderID, payload.UserType, payload.Content)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Message sent successfully!"}, http.StatusCreated)
}