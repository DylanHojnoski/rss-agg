package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
    if code > 499 {
        log.Println("Responding with 5XX error:", msg)
    }
    type errorResponse struct {
        Error string `json:"error"`
    }

    respondWithJSON(w, code, errorResponse {
        Error: msg,
    })
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("failed to marshal JSON response: %v", payload)
        w.WriteHeader(500)
        return
    }
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func respondWithFile(w http.ResponseWriter, r *http.Request, code int, payload []byte) {
     err := os.WriteFile("podcast.opml", payload, 0666)
     if err != nil {
        respondWithError(w,400, fmt.Sprintf("Couldn't make file: %v", err))
        return
     }

    w.Header().Add("Content-Type", "application/xml")
    w.WriteHeader(code)
    http.ServeFile(w, r , "podcast.opml")

    os.Remove("podcast.opml")

}
