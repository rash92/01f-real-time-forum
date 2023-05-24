package auth

import (
	"forum/dbmanagement"
	"forum/utils"
	"net/http"
	"time"
)

/*
Returns the cookie value of the current session that gives a sessions ID.  Used to determine which user is using the program.
*/

const timeout = 30 * time.Minute

func GetSessionFromBrowser(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := r.Cookie("session")
	utils.HandleError("Unable to get session from browser Err:", err)

	if err != nil {
		return "", err
	}

	value := session.Value
	return value, err
}

/*
Creates session that gives a sessions ID, used to determine which user is using the program.
*/
func CreateUserSession(w http.ResponseWriter, r *http.Request, user dbmanagement.User) error {
	session, err := user.CreateSession()
	utils.HandleError("Unable to create user session err:", err)
	cookie := http.Cookie{
		Name:     "session",
		Value:    session.UUID,
		Expires:  time.Now().Add(timeout),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	return err
}
