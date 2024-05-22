package auth

import "strings"

type Role struct {
	UserType string
	Privilege int
}

var (
	RoleTeachingAssistant = Role {
		UserType: "teachingAssistant",
		Privilege: 1,
	}

	RoleStudent = Role {
		UserType: "student",
		Privilege: 2,
	}
)

func GetRoleByEmail(email string) Role {
	// Use regex function
	if strings.Contains(email, "student") {
		return RoleStudent
	} 
	
	return RoleTeachingAssistant
}