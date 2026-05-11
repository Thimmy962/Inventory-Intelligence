package main

import (
	"context"
	"database/sql"
	"log"
	"main/internal/app"
	"main/internal/database"
	"main/internal/handler"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var port = "8000"

var requiredEnv = []string{
	"DB_URL",
	"SECRET_STRING",
}

func MustLoadEnv() {
	var enverror = 0
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Printf("env vars '%v', is missing\n", key)
			enverror++
		}
	}

	if enverror > 0 {
		log.Println("The environmental Variables ablove are missing")
		os.Exit(1)
	}
}

func main() {
	MustLoadEnv()
	workers := 3
	var err error // Initial declaration

	// Use = because err is already declared
	if len(os.Args) > 1 {
		sworkers := os.Args[1]
		// Use = here as well to update the existing variables
		workers, err = strconv.Atoi(sworkers)
		if err != nil {
			log.Printf("Could not convert '%s' to int", sworkers)
			os.Exit(1)
		}
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	dbServer := &app.Server{
		DB:        db,
		Queries:   database.New(db),
		SecretKey: os.Getenv("secret_key"),
	}

	// serMux := http.NewServeMux()
	serMux := mux.NewRouter().StrictSlash(true)
	server := http.Server{Addr: ":" + port, Handler: serMux}
	channel := make(chan string, 100)
	var wg sync.WaitGroup
	dbQuery := handler.NewHandler(dbServer, channel, &wg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go dbQuery.StartWorker(ctx, workers)

	serMux.HandleFunc("/bulkcreate", dbServer.CORSMiddleware(dbQuery.BulkCreateProducts)).Methods("POST")
	serMux.HandleFunc("/newproduct", dbServer.CORSMiddleware(dbQuery.CreateProduct)).Methods("POST")
	serMux.HandleFunc("/product/{id}", dbServer.CORSMiddleware(dbQuery.GetProduct)).Methods("GET")
	serMux.HandleFunc("/products", dbServer.CORSMiddleware(dbQuery.GetProducts)).Methods("GET")
	serMux.HandleFunc("/bulkpurchase", dbServer.CORSMiddleware(dbQuery.NewBulkPurchase)).Methods("POST")
	serMux.HandleFunc("/purchase", dbServer.CORSMiddleware(dbQuery.NewPurchase)).Methods("POST")
	serMux.HandleFunc("/sales", dbServer.CORSMiddleware(dbQuery.CreateSales)).Methods("POST")
	serMux.HandleFunc("/index", dbServer.HTMLCORSMiddleware(dbQuery.Index)).Methods("GET")
	serMux.HandleFunc("/", dbServer.CheckOutCORSMiddleware(dbQuery.Checkout)).Methods("GET")
	serMux.HandleFunc("/search", dbServer.CORSMiddleware(dbQuery.Search)).Methods("GET")
	serMux.HandleFunc("/edit/{id}", dbServer.HTMLCORSMiddleware(dbQuery.EditProduct)).Methods("GET", "PUT")
	serMux.HandleFunc("/adjustment", dbServer.CORSMiddleware(dbQuery.Adjustment)).Methods("POST")
	serMux.HandleFunc("/topselling", dbServer.CheckOutCORSMiddleware(dbQuery.TopProducts)).Methods("GET")
	serMux.HandleFunc("/trends/{pid}", dbServer.CORSMiddleware(dbQuery.ProductTrend)).Methods("GET")
	serMux.HandleFunc("/login", dbServer.LoginCORSMiddleware(dbQuery.Login)).Methods("GET", "POST")
	serMux.HandleFunc("/register", dbServer.HTMLCORSMiddleware(dbQuery.Register)).Methods("GET", "POST")

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		log.Println("Received Signal: ", sig)
		log.Println("Shutting Down...")
		server.Shutdown(ctx)
		cancel()
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
