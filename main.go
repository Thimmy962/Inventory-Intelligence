package main

import (
	"net/http"
)

var port = "8000"


func main() {
	serMux := http.NewServeMux()
	server := http.Server{Addr: ":"+ port, Handler: serMux}
	

	server.ListenAndServe();
}