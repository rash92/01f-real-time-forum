package auth

import (
	"forum/dbmanagement"
	"forum/utils"
	"net/http"
	"time"
)

const Limit = 100

func LimitRequests(w http.ResponseWriter, r *http.Request, user dbmanagement.User) dbmanagement.User {
	limitTime := time.Minute * 1
	userSession, err := user.ReturnSession(user.UUID)
	if err != nil {
		utils.HandleError("Unable to get user or no user", err)
		return dbmanagement.User{}
	}
	startTime := userSession.CreatedAt
	endTime := startTime.Add(limitTime)

	go func() {
		for {
			time.Sleep(15 * time.Second)
			_, err := dbmanagement.SelectUserFromSession(userSession.UUID)
			if err != nil {
				utils.HandleError("Unable to get user or no user", err)
				continue
			}
			if CheckTime(endTime) {
				dbmanagement.UpdateUserToken(userSession.UserId, Limit)
				startTime = endTime
				endTime = startTime.Add(limitTime)
			}
		}
	}()

	return user
}

func CheckTime(endTime time.Time) bool {
	return time.Now().After(endTime)
}
