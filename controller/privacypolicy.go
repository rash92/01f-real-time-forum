package controller

import (
	"html/template"
	"net/http"
)

func PrivacyPolicy(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.ExecuteTemplate(w, "privacypolicy.html", nil)
}
