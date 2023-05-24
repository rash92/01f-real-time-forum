package controller

import (
	"html/template"
	"net/http"
)

type Errors struct {
	Errorcode    int
	Errormessage string
}

func PageErrors(w http.ResponseWriter, r *http.Request, tmpl *template.Template, errorcode int, errormessage string) {
	errors := Errors{
		Errorcode:    errorcode,
		Errormessage: errormessage,
	}
	w.WriteHeader(errorcode)
	tmpl.ExecuteTemplate(w, "error.html", errors)
}
