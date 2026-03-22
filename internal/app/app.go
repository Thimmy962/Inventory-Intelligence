package app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"main/internal/database"
	"net/http"
)


type Server struct {
	DB *sql.DB
	Queries *database.Queries
}

func (dbserver *Server)SayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}



func (s *Server) CORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		next(w, req)
	}
}