package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"math/rand"
	"time"
)

const PossibleChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const AttendanceCodeDuration = 5

// Creates a random attendance string for a tutorial, and stores it in the database
func CreateRandomAttendanceString(tutorialId int) (*models.AttendanceString, error) {
	bytesArr := make([]byte, 10)
	for i := range bytesArr {
		char := PossibleChars[rand.Intn(len(PossibleChars))]
		bytesArr[i] = char
	}

	code := string(bytesArr)
	attendanceString := &models.AttendanceString{Code: code, TutorialID: tutorialId, ExpiresAt: time.Now().UTC().Add(time.Minute * AttendanceCodeDuration)}
	result := database.DB.Table("attendance_strings").Create(attendanceString)
	if result.Error != nil {
		return nil, result.Error
	}

	return attendanceString, nil
}

// Gets the ttendance string for a tutorial for the current date
func GetAttendanceStringByTutorialID(tutorialId int) (*models.AttendanceString, error) {
	var attendanceString models.AttendanceString
	result := database.DB.Table("attendance_strings").Where("tutorial_id = ?", tutorialId).First(&attendanceString)
	if result.Error != nil {
		return nil, result.Error
	}

	return &attendanceString, nil
}

// Gets all attendance for a tutorial on a specific date sorted by student ID in ascending order
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

// Gets all attendance for a tutorial on the current date sorted by student ID in ascending order
func GetTodayAttendanceByTutorialID(tutorialId int) (*[]models.Attendance, error) {
	date := time.Now().UTC().Format("2006-01-02")
	return GetAttendanceByDateAndTutorialID(date, tutorialId)
}

// Gets all attendance for a tutorial on all tutorial sessions,
// Sorted by date in descending order, and student ID in ascending order
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

// Gets all attendance of a student in a tutorial for all tutorial sessions,
// Sorted by date in descending order
func GetStudentAttendance(tutorialId int, studentId int) (*[]models.Attendance, error) {
	var attendances []models.Attendance
	result := database.DB.Table("attendances").Where("tutorial_id = ?", tutorialId).
		Where("student_id = ?", studentId).
		Order("date DESC").
		Find(&attendances)

	if result.Error != nil {
		return nil, result.Error
	}

	return &attendances, nil
}

// Deletes the attendance string for a tutorial
func DeleteGeneratedAttendanceString(tutorialId int) error {
	attendanceString, err := GetAttendanceStringByTutorialID(tutorialId)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("attendance_strings").Delete(attendanceString)
	return result.Error
}

// Generates attendance for a tutorial for the current date for all students in the tutorial
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

// Generates attendance for a tutorial for a specific date for all students in the tutorial
func GenerateAttendanceForDateByTutorialID(date string, tutorialId int) error {
	studentIds, err := GetAllStudentIdsOfStudentsInTutorial(tutorialId)
	if err != nil {
		return err
	}

	for _, studentId := range *studentIds {
		attendance := &models.Attendance{StudentID: studentId, TutorialID: tutorialId, Date: date}
		result := database.DB.Table("attendances").Create(attendance)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Deletes attendance for a tutorial for the current date
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

// Deletes attendance for a tutorial for a specific date
func DeleteAttendanceForDateByTutorialID(date string, tutorialId int) error {
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

// Verifies the attendance code for a tutorial
// By checking if the attendance code exists and has not expired
// And checks if the attendance code matches the attendance code for the attendance string
func VerifyAttendanceCode(tutorialId int, attendanceCode string) (bool, error) {
	var attendanceString models.AttendanceString
	result := database.DB.Table("attendance_strings").Where("code = ?", attendanceCode).
		Where("tutorial_id = ?", tutorialId).First(&attendanceString)
	if result.Error != nil {
		return false, result.Error
	}

	if attendanceString.ExpiresAt.Before(time.Now().UTC()) {
		return false, nil
	}

	return true, nil
}

// Marks a student as present for a tutorial on the current date
func MarkPresent(studentId int, tutorialId int) error {
	var attendance models.Attendance
	date := time.Now().UTC().Format("2006-01-02")
	result := database.DB.Table("attendances").Where("student_id = ?", studentId).
		Where("tutorial_id = ?", tutorialId).Where("date = ?", date).First(&attendance)
	if result.Error != nil {
		return result.Error
	}

	attendance.Present = true
	database.DB.Table("attendances").Save(attendance)
	return nil
}

// Gets the attendance of a student for a tutorial on the current date
func GetTodayAttendanceByStudentId(studentId int, tutorialId int) (*models.Attendance, error) {
	var attendance models.Attendance
	date := time.Now().UTC().Format("2006-01-02")
	result := database.DB.Table("attendances").Where("student_id = ?", studentId).
		Where("tutorial_id = ?", tutorialId).Where("date = ?", date).First(&attendance)
	if result.Error != nil {
		return nil, result.Error
	}
	return &attendance, nil
}
