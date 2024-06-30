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

type UploadFilePayload struct {
	Name string `json:"name"`
	Week int `json:"week"`
	Filepath string `json:"filepath"`
}

type RequestResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type FilepathPayload struct {
	Filepath string `json:"filepath"`
}

type CreateMesssagePayload struct {
	SenderID int `json:"senderId"`
	UserType string `json:"userType"`
	Content string `json:"content"`
}

type MarkAttendancePayload struct {
	StudentID int `json:"studentId"`
	AttendanceCode string `json:"attendanceCode"`
}

type DeleteAttendanceStringPayload struct {
	AttendanceCode string `json:"attendanceCode"`
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

type TutorialFilesResponse struct {
	Files []models.TutorialFile `json:"files"`
}

type FilepathResponse struct {
	Filepath string `json:"filepath"`
}

type MessageResponse struct {
	TutorialID int `json:"tutorialId"`
	Sender string `json:"sender"`
	SenderID int `json:"senderId"`
	UserType string `json:"userType"`
	Content string `json:"content"`
}

type MessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}

type ConsultationResponse struct {
	ID uint `json:"id"`
	Tutorial models.Tutorial `json:"tutorial"`
	Student models.Student `json:"student"`
	TeachingAssistant models.TeachingAssistant `json:"teachingAssistant"`
	Date string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Booked bool `json:"booked"`
}

type ConsultationsResponse struct {
	Consultations []ConsultationResponse `json:"consultations"`
}

type BookedConsultationsByDate struct {
	Date string `json:"date"`
	Consultations []ConsultationResponse `json:"consultations"`
}

type BookedConsultationsResponse struct {
	BookedConsultations []BookedConsultationsByDate `json:"bookedConsultations"`
}

type AttendanceStringResponse struct {
	AttendanceString models.AttendanceString `json:"attendanceString"`
}

type AttendanceResponse struct {
	ID uint `json:"id"`
	Student models.Student `json:"student"`
	TutorialID int `json:"tutorialId"`
	Date string `json:"date"`
	Present bool `json:"present"`
}

type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
}

type AttendanceListByDate struct {
	Date string `json:"date"`
	Attendance []AttendanceResponse `json:"attendance"`
}

type AttendanceListsByDateResponse struct {
	AttendanceLists []AttendanceListByDate `json:"attendanceLists"`
}

type StudentAttendanceResponse struct {
	Attendance []models.Attendance `json:"attendance"`
}