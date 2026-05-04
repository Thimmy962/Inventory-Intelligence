package app

import (
	"database/sql"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"os"
	"context"
	_ "github.com/lib/pq"
)


type Server struct {
	DB *sql.DB
	Queries *database.Queries
	SecretKey string
}

func (dbserver *Server)SayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}



func (s *Server) CORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
		// get the token and id form request
		id, err := auth.GetBearerToken(req.Header, os.Getenv("SECRET_STRING"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// if userid is_manage
		level, err := s.Queries.GetLevel(req.Context(), id) // get the level of a staff 
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}


		ctx := context.WithValue(req.Context(), id, id)
		ctx = context.WithValue(ctx, level, level)
		req = req.WithContext(ctx)
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		next(w, req)
	}
}

func (s *Server) HTMLCORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
        // CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/html")

		next(w, req)
	}
}

