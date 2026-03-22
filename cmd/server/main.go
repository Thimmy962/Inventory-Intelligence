package main

import (
	"database/sql"
	"net/http"
	"os"
	"log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main/internal/database"
	"main/internal/app"
	"main/internal/handler"
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

	serMux := http.NewServeMux()
	server := http.Server{Addr: ":"+ port, Handler: serMux}
	dbQuery := handler.NewHandler(dbServer)
	
	serMux.HandleFunc("GET /", dbServer.CORSMiddleware(dbServer.SayHello))
	serMux.HandleFunc("POST /newproduct", dbServer.CORSMiddleware(dbQuery.CreateProduct))
	serMux.HandleFunc("GET /product/{id}", dbServer.CORSMiddleware(dbQuery.GetProduct))
	serMux.HandleFunc("GET /products", dbServer.CORSMiddleware(dbQuery.GetProducts))
	serMux.HandleFunc("POST /bulkpurchase", dbServer.CORSMiddleware(dbQuery.NewBulkPurchase))
	serMux.HandleFunc("POST /purchase", dbServer.CORSMiddleware(dbQuery.NewPurchase))
	serMux.HandleFunc("POST /sales", dbServer.CORSMiddleware(dbQuery.CreateSales))

	log.Println(server.ListenAndServe());
}
