package main

// import (
// 	"database/sql"
// 	"html/template"
// 	"log"
// 	"main/internal/app"
// 	"main/internal/database"
// 	"main/internal/handler"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// )

// var port = "8000"

// func main() {
// 	err := godotenv.Load(); if err != nil {
// 		log.Println(err)
// 		os.Exit(1)
// 	}

// 	dbURL := os.Getenv("DB_URL")
// 	db, err := sql.Open("postgres", dbURL)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 		os.Exit(1)
// 	}

// 	dbServer := &app.Server {
// 		DB: db,
// 		Queries: database.New(db),		
// 	}

// 	serMux := http.NewServeMux()
	
// 	server := http.Server{Addr: ":"+ port, Handler: serMux}
// 	temp, err := template.ParseGlob("template/*.html")
// 	if err != nil {
// 		log.Fatal((err.Error()))
// 		os.Exit(1)
// 	}
// 	dbQuery := handler.NewHandler(dbServer, temp)
	

// 	serMux.HandleFunc("POST /newproduct", dbServer.CORSMiddleware(dbQuery.CreateProduct))
// 	serMux.HandleFunc("GET /product/{id}", dbServer.CORSMiddleware(dbQuery.GetProduct))
// 	serMux.HandleFunc("GET /products", dbServer.CORSMiddleware(dbQuery.GetProducts))
// 	serMux.HandleFunc("POST /bulkpurchase", dbServer.CORSMiddleware(dbQuery.NewBulkPurchase))
// 	serMux.HandleFunc("POST /purchase", dbServer.CORSMiddleware(dbQuery.NewPurchase))
// 	serMux.HandleFunc("POST /sales", dbServer.CORSMiddleware(dbQuery.CreateSales))
// 	// serMux.HandleFunc("GET /index", dbServer.HTMLCORSMiddleware(dbQuery.Index))

// 	log.Println(server.ListenAndServe());
// }


import (
	"database/sql"
	"html/template"
	"log"
	"main/internal/app"
	"main/internal/database"
	"main/internal/handler"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var port = "8000"

func main() {
	err := godotenv.Load(); if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	dbServer := &app.Server {
		DB: db,
		Queries: database.New(db),		
	}

	// serMux := http.NewServeMux()
	serMux := mux.NewRouter().StrictSlash(true)
	server := http.Server{Addr: ":"+ port, Handler: serMux}
	temp, err := template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatal((err.Error()))
		os.Exit(1)
	}
	dbQuery := handler.NewHandler(dbServer, temp)
	

	serMux.HandleFunc("/newproduct", dbServer.CORSMiddleware(dbQuery.CreateProduct)).Methods("POST")
	serMux.HandleFunc("/product/{id}", dbServer.CORSMiddleware(dbQuery.GetProduct)).Methods("GET")
	serMux.HandleFunc("/products", dbServer.CORSMiddleware(dbQuery.GetProducts)).Methods("GET")
	serMux.HandleFunc("/bulkpurchase", dbServer.CORSMiddleware(dbQuery.NewBulkPurchase)).Methods("POST")
	serMux.HandleFunc("/purchase", dbServer.CORSMiddleware(dbQuery.NewPurchase)).Methods("POST")
	serMux.HandleFunc("/sales", dbServer.CORSMiddleware(dbQuery.CreateSales)).Methods("POST")
	serMux.HandleFunc("/index", dbServer.HTMLCORSMiddleware(dbQuery.Index)).Methods("GET")

	log.Println(server.ListenAndServe());
}