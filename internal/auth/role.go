package auth

type Role struct {
	UserType string `json:"userType"`
	Privilege int `json:"privilege"`
}

// The roles for the 2 types of users that use our application
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
