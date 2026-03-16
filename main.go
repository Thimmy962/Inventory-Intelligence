package main

import (
	"database/sql"
	"net/http"
	"os"
	"log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main/internal/database"
)

var port = "8000"

type Server struct {
	db *sql.DB
	queries *database.Queries
}

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

	dbServer := &Server {
		db: db,
		queries: database.New(db),		
	}

	serMux := http.NewServeMux()
	server := http.Server{Addr: ":"+ port, Handler: serMux}
	
	serMux.HandleFunc("GET /", dbServer.CORSMiddleware(dbServer.sayHello))
	serMux.HandleFunc("POST /newproduct/", dbServer.CORSMiddleware(dbServer.CreateProduct))
	serMux.HandleFunc("POST /bulkcreate/", dbServer.CORSMiddleware(dbServer.BulkCreateProducts))
	serMux.HandleFunc("GET /products/", dbServer.CORSMiddleware(dbServer.GetProducts))
	serMux.HandleFunc("GET /product/{id}/", dbServer.CORSMiddleware(dbServer.GetProduct))

	log.Println(server.ListenAndServe());
}

func (dbserver *Server)sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}



func (s *Server) CORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

		next(w, req)
	}
}
