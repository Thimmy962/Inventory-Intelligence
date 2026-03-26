package handler

import (
	"html/template"
	// "log"
	"main/internal/database"
	"net/http"
)


func (db *Handler) Index(wr http.ResponseWriter, req *http.Request) {
	list, _ := db.server.Queries.GetFullProductDetail(req.Context())
	tmpl := template.Must(template.ParseFiles(
 	   "template/layout.html",
  	  "template/index.html",
	))
	lists := map[string][]database.GetFullProductDetailRow {
		"products": list,
	}
	tmpl.ExecuteTemplate(wr, "layout.html", lists)
}


func (handler *Handler) Checkout(wr http.ResponseWriter, req *http.Request) {

}