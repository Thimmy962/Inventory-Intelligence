package handler

import (
	"fmt"
	"html/template"
	"log"
	"main/internal/database"
	"net/http"

	"github.com/gorilla/mux"
)


func (db *Handler) Index(wr http.ResponseWriter, req *http.Request) {
	list, _ := db.server.Queries.GetFullProductDetailView(req.Context())
	tmpl := template.Must(template.ParseFiles(
 	   "template/layout.html",
  	  "template/index.html",
	))
	lists := map[string][]database.Productdetail {
		"products": list,
	}
	
	tmpl.ExecuteTemplate(wr, "layout.html", lists)
}

func (handler *Handler) Checkout(wr http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"template/layout.html", "template/checkout.html",
	))
	tmpl.ExecuteTemplate(wr, "layout.html", nil)
}

func (handler *Handler)Search(wr http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("search")
	res, err := handler.server.Queries.SearchProductForCheckout(req.Context(), query)
	if err != nil {
		log.Println(err)
		ProcessingError(wr, http.StatusNotFound, fmt.Errorf("%s not found", query))
		return
	}

	if len(res) == 0 {
 	   ProcessingError(wr, http.StatusNotFound, fmt.Errorf("%s not found", query))
  	  return
	}

	respondWithJSON(wr, 200, res)
}

func (handler *Handler)EditProduct(wr http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    id := vars["id"]
    product, err := handler.server.Queries.GetOneFullProductDetail(req.Context(), id)
    if err != nil {
        tmpl := template.Must(template.ParseFiles(
            "template/layout.html", "template/error.html",
        ))
        tmpl.ExecuteTemplate(wr, "layout.html", nil)
    } else {
        tmpl := template.Must(template.ParseFiles(
            "template/layout.html", "template/edit.html",
        ))
        tmpl.ExecuteTemplate(wr, "layout.html", product)
    }
}
