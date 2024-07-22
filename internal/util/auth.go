package util

import "golang.org/x/crypto/bcrypt"

// Checks whether the password a user has typed in is valid, by comparing it with the hash stored in the database
func VerifyPassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
