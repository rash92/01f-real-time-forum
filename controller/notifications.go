package controller

import (
	"encoding/json"
	"forum/dbmanagement"
	"html/template"
	"net/http"
)

type NotificationFormData struct {
	NotificationDelete string `json:"deleteNotification"`
}

func Notification(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// Parse the JSON body into FormData struct
	var formData NotificationFormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		// Handle error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	notificationToDelete := formData.NotificationDelete
	if notificationToDelete != "" {
		dbmanagement.DeleteFromTableWithUUID("Notifications", notificationToDelete)
	}
}
