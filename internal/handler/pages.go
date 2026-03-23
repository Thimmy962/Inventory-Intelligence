package handler

import (
	"html/template"
	"net/http"
)


func (db *Handler) Index(wr http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles(
 	   "template/layout.html",
  	  "template/index.html",
	))
	tmpl.ExecuteTemplate(wr, "layout.html",nil)
}