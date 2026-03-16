package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// called on bad request
func ProcessingError(w http.ResponseWriter, code int, Err error) {
	w.Header().Set("Content-Type", "Apllication/json")
	w.WriteHeader(code)
	w.Write([]byte(Err.Error()))
}



// handles both writing a single struct or a list of structs to responseWriter
// all thanks to NewEncoder
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, 500, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(code)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err.Error())
	}
}
