package handler

import (
	"net/http"
)


func (db *Handler) Index(wr http.ResponseWriter, req *http.Request) {
	db.tmpl.ExecuteTemplate(wr, "index.html", nil)
}