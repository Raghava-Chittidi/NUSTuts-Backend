package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateStudent(name string, email string, password string, modules []string) (*models.Student, error) {
	student := &models.Student{Name: name, Email: email, Password: password, Modules: modules}
	result := database.DB.Table("students").Create(student)
	return student, result.Error
}

func GetStudentById(id int) (*models.Student, error) {
	var student models.Student
	result := database.DB.Table("students").Where("id = ?", id).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}

func GetStudentByEmail(email string) (*models.Student, error) {
	var student models.Student
	result := database.DB.Table("students").Where("email = ?", email).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}

func DeleteStudentByEmail(email string) error {
	student, err := GetStudentByEmail(email)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("students").Delete(&student)
	if result.Error != nil {
		return result.Error
	}

	return nil
}