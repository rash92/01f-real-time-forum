package auth

import (
	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

/*
Changes a plain text password into a hashed key password and returns that key.  Used to store password as hash key in the database.
*/
func HashPassword(p string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	utils.HandleError("Password hashing error", err)
	return string(h)
}

/*
Compares plain text to a hashed key and returns a boolean.  Used for authentification of user's password.
*/
func CompareHash(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
