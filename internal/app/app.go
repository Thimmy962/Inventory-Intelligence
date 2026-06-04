package app

import (
	"context"
	"database/sql"
	"html/template"
	"main/internal/auth"
	"main/internal/database"
	"net/http"

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
		manager_cookies, err := req.Cookie("admin")
		if err != nil || manager_cookies.Value == "false"{
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		access_cookie, err := req.Cookie("access_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// get the token and id form request
		id, err := auth.GetBearerToken(access_cookie.Value, s.SecretKey)
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


		ctx := context.WithValue(req.Context(), "id", id)
		ctx = context.WithValue(ctx, "level", level)
		req = req.WithContext(ctx)
        	// CORS headers
	        w.Header().Set("Access-Control-Allow-Origin", "*")
        	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        	w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		next(w, req)
	}
}


func (s *Server) renderError(wr http.ResponseWriter, statusCode int) {

	tmpl := template.Must(template.ParseFiles(
		"template/layout.html","template/error2.html",
	))
	errMsg := struct{
		StatusCode int
		Message string
	} {StatusCode: statusCode, Message: "Unauthorized"}
	wr.WriteHeader(statusCode)

	tmpl.ExecuteTemplate(wr, "layout.html", errMsg)
}


func (s *Server) HTMLCORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
		managerCookies, errAdmin := req.Cookie("admin")
		accessCookie, errAccess := req.Cookie("access_token")
		
		if errAdmin != nil || managerCookies.Value == "false" || errAccess != nil {
			s.renderError(w, http.StatusUnauthorized)
			return
		}	


		// get the token and id form request
		id, err := auth.GetBearerToken(accessCookie.Value, s.SecretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(req.Context(), "id", id)
		req = req.WithContext(ctx)
        	// CORS headers
	        w.Header().Set("Access-Control-Allow-Origin", "*")
        	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        	w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/html")
		next(w, req)
	}
}

func (s *Server) LoginCORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
        	// CORS headers
	        w.Header().Set("Access-Control-Allow-Origin", "*")
        	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        	w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/html")
		next(w, req)
	}
}

func (s *Server) CheckOutCORSMiddleware(next http.HandlerFunc)  http.HandlerFunc{
	return  func(w http.ResponseWriter, req *http.Request) {
		accessCookie, errAccess := req.Cookie("access_token")


		if errAccess != nil {
			s.renderError(w, http.StatusUnauthorized)
			return
		}

		// get the token and id form request
		userID, err := auth.GetBearerToken(accessCookie.Value, s.SecretKey)

		if err != nil {
			s.renderError(w, http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(req.Context(), "id", userID)
		req = req.WithContext(ctx)

        	// CORS headers
	        w.Header().Set("Access-Control-Allow-Origin", "*")
        	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Apikey")
        	w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/html")
		next(w, req)
	}
}


