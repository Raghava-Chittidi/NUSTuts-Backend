package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"math/rand"
	"time"
)

const PossibleChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const AttendanceCodeDuration = 5

func CreateRandomAttendanceString(tutorialId int) (*models.AttendanceString, error) {
	bytesArr := make([]byte, 10)
	for i := range bytesArr {
		char := PossibleChars[rand.Intn(len(PossibleChars))]
		bytesArr[i] = char
	}

	code := string(bytesArr)
	attendanceString := &models.AttendanceString{Code: code, TutorialID: tutorialId, ExpiresAt: time.Now().Add(time.Minute * AttendanceCodeDuration)}
	result := database.DB.Table("attendance_strings").Create(attendanceString)
	if result.Error != nil {
		return nil, result.Error
	}

	return attendanceString, nil
}

func GetAttendanceStringByTutorialID(tutorialId int) (*models.AttendanceString, error) {
	var attendanceString models.AttendanceString
	result := database.DB.Table("attendance_strings").Where("tutorial_id = ?", tutorialId).First(&attendanceString)
	if result.Error != nil {
		return nil, result.Error
	}

	return &attendanceString, nil
}

func GetAttendanceByDateAndTutorialID(date string, tutorialId int) (*[]models.Attendance, error) {
	var attendances []models.Attendance
	result := database.DB.Table("attendances").Where("date = ?", date).Where("tutorial_id = ?", tutorialId).
		Order("student_id ASC").
		Find(&attendances)
	if result.Error != nil {
		return nil, result.Error
	}

	return &attendances, nil
}

func GetTodayAttendanceByTutorialID(tutorialId int) (*[]models.Attendance, error) {
	date := time.Now().UTC().Format("2006-01-02")
	return GetAttendanceByDateAndTutorialID(date, tutorialId)
}

func GetAllAttendanceByTutorialID(tutorialId int) (*[]models.Attendance, error) {
	var attendances []models.Attendance
	result := database.DB.Table("attendances").
		Where("tutorial_id = ?", tutorialId).
		Order("date DESC").
		Order("student_id ASC").
		Find(&attendances)

	if result.Error != nil {
		return nil, result.Error
	}

	return &attendances, nil
}

func DeleteGeneratedAttendanceString(tutorialId int) error {
	attendanceString, err := GetAttendanceStringByTutorialID(tutorialId)
	if err != nil {
		return err
	}

	result := database.DB.Table("attendance_strings").Delete(attendanceString)
	return result.Error
}

func GenerateTodayAttendanceByTutorialID(tutorialId int) error {
	studentIds, err := GetAllStudentIdsOfStudentsInTutorial(tutorialId)
	if err != nil {
		return err
	}

	date := time.Now().UTC().Format("2006-01-02")
	for _, studentId := range *studentIds {
		attendance := &models.Attendance{StudentID: studentId, TutorialID: tutorialId, Date: date}
		result := database.DB.Table("attendances").Create(attendance)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func DeleteTodayAttendanceByTutorialID(tutorialId int) error {
	date := time.Now().UTC().Format("2006-01-02")
	attendances, err := GetAttendanceByDateAndTutorialID(date, tutorialId)
	if err != nil {
		return err
	}

	for _, attendance := range *attendances {
		result := database.DB.Table("attendances").Delete(&attendance)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func VerifyAttendanceCode(tutorialId int, attendanceCode string) (bool, error) {
	var attendanceString models.AttendanceString
	result := database.DB.Table("attendance_strings").Where("code = ?", attendanceCode).
				Where("tutorial_id = ?", tutorialId).First(&attendanceString)
	if result.Error != nil {
		return false, result.Error
	}

	if attendanceString.ExpiresAt.Before(time.Now())  {
		return false, nil
	}

	return true, nil
}

func MarkPresent(studentId int, tutorialId int) error {
	var attendance models.Attendance
	result := database.DB.Table("attendances").Where("student_id = ", studentId).Where("tutorial_id = ?", tutorialId).First(&attendance)
	if result.Error != nil {
		return result.Error
	}

	attendance.Present = true
	database.DB.Table("attendances").Save(attendance)
	return nil
}