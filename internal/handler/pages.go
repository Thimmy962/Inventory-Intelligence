package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"time"
	"golang.org/x/text/cases"
    	"golang.org/x/text/language"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// func (db *Handler) GetUser()

func (db *Handler) TopProducts(wr http.ResponseWriter, req *http.Request) {
	list, _ := db.server.Queries.GetTopProduct(req.Context())
	respondWithJSON(wr, http.StatusOK, list)
}

func (db *Handler) ProductTrend(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	product_id := vars["pid"]

	trends, err := db.server.Queries.ProductTrend(req.Context(), product_id)
	if err != nil {
		ProcessingError(wr, http.StatusNotFound, nil)
		return
	}
	respondWithJSON(wr, http.StatusOK, trends)
}

func (db *Handler) Register(wr http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		staff := struct {
			Fname      string    `json:"first_name"`
			Lname      string    `json:"last_name"`
			Username   string    `json:"username"`
			Password   string    `json:"password"`
			Manager bool      `json:"is_manager"`
		}{}

		err := json.NewDecoder(req.Body).Decode(&staff)
		if err != nil {
			log.Println(err)
			ProcessingError(wr, http.StatusBadRequest, err)
			return
		}
		
		caser := cases.Title(language.English)
		
		staff.Fname = caser.String(staff.Fname)
		staff.Lname = caser.String(staff.Lname)
		staff.Username = caser.String(staff.Username)
		
		staff.Password, err = auth.CreateHash(staff.Password)
		if err != nil {
			ProcessingError(wr, http.StatusBadRequest, err)
			return
		}
		err = db.server.Queries.CreateStaff(req.Context(),
			database.CreateStaffParams{FirstName: staff.Fname, LastName: staff.Lname,
			 Pword: staff.Password, IsManager: staff.Manager, Username: staff.Username})
		if err != nil {
			ProcessingError(wr, http.StatusNotFound, err)
			return
		}
		respondWithJSON(wr, http.StatusOK, nil)
		return
	}
	tmpl := template.Must(template.ParseFiles(
		"template/layout.html", "template/register.html",
	))
	tmpl.ExecuteTemplate(wr, "layout.html", nil)}

func (db *Handler) Login(wr http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		loginDetail := struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Pword string
		}{}
		err := json.NewDecoder(req.Body).Decode(&loginDetail)
		if err == io.EOF {
			log.Println(err)
			ProcessingError(wr, http.StatusBadRequest, err)
			return
		}		
		caser := cases.Title(language.English)
		
		loginDetail.Username = caser.String(loginDetail.Username)
		// get username details
		data, err := db.server.Queries.LoginUser(req.Context(), loginDetail.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = fmt.Errorf("Username or password not valid")
			}
			ProcessingError(wr, http.StatusNotFound, err)
			return
		}
		// check hashed password for that username and password coming from frontend
		// converts password to hash and compares to the hashed passed to it
		// return error or bool
		match, err := auth.CheckHash(loginDetail.Password, data.Pword)
		if err != nil {
			ProcessingError(wr, http.StatusBadRequest, err)
			return
		}
		if !match {
			err = fmt.Errorf("Staff with username and password not found")
			ProcessingError(wr, http.StatusNotFound, err)
			return
		}

		// convert id in string to uuid
		id := uuid.MustParse(data.ID)

		// get authentication tokens
		auth_token, err := auth.MakeTokens(db.server.SecretKey, id, 5 * time.Hour)
		if err != nil {
			ProcessingError(wr, http.StatusInternalServerError, err)
			return
		}
		data1 := struct{
			Username string
			ID uuid.UUID
			Is_Manager bool
			AuthToken string
		}{
			Username: data.Username, ID: id, Is_Manager: data.IsManager, AuthToken: auth_token,
		}
		respondWithJSON(wr, http.StatusOK, data1)
		return
	}
	tmpl := template.Must(template.ParseFiles(
		"template/login.html",
	))

	tmpl.ExecuteTemplate(wr, "login.html", nil)
}

func (db *Handler) Index(wr http.ResponseWriter, req *http.Request) {
	list, _ := db.server.Queries.GetTopProduct(req.Context())
	tmpl := template.Must(template.ParseFiles(
		"template/layout.html",
		"template/index.html",
	))
	lists := map[string][]database.GetTopProductRow{
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

func (handler *Handler) Search(wr http.ResponseWriter, req *http.Request) {
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

func (handler *Handler) EditProduct(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if req.Method == "PUT" {
		var product Product
		err := json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			ProcessingError(wr, 400, err)
			return
		}
		err = handler.server.Queries.EditOneProduct(req.Context(), database.EditOneProductParams{
			ID: product.ID, ProductName: product.ProductName, Price: product.Price, ReorderLevel: product.ReorderLevel,
		})
		if err != nil {
			log.Println(err)
			ProcessingError(wr, 400, err)
			return
		}
		respondWithJSON(wr, http.StatusAccepted, nil)
		return
	}

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
		// database.GetOneFullProductDetailRow
		tmpl.ExecuteTemplate(wr, "layout.html", product)
	}
}
