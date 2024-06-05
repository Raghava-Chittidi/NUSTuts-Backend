package api

import (
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/models"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Message string `json:"message"`
	Error error `json:"error,omitempty"`
}

type RequestToJoinTutorialPayload struct {
	StudentID int `json:"studentId"`
	ModuleCode string `json:"moduleCode"`
	ClassNo string `json:"classNo"`
}

type RequestResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type StudentAuthResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role auth.Role `json:"role"`
	Modules []string `json:"modules"`
	Tutorials []models.Tutorial `json:"tutorials"`
	Tokens auth.Tokens `json:"tokens"`
}

type TeachingAssistantAuthResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role auth.Role `json:"role"`
	Tutorial models.Tutorial `json:"tutorial"`
	Tokens auth.Tokens `json:"tokens"`
}