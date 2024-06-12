package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	res := api.MessagesResponse{Messages: *messages}
	util.WriteJSON(w, api.Response{Message: "Fetched messages successfully!", Data: res}, http.StatusOK)
}

func CreateMessageForTutorial(w http.ResponseWriter, r *http.Request) {
	var payload api.CreateMesssagePayload
	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = dataaccess.CreateMessage(payload.DiscussionID, payload.SenderID, payload.UserType)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, api.Response{Message: "Message posted successfully!"}, http.StatusCreated)
}