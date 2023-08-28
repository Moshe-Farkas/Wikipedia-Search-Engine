package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlequery)
	log.Fatal(http.ListenAndServe(":1234", mux))
}

func handlequery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	query := r.URL.Query().Get("q")
	

	fmt.Fprintf(w, `{
		"query": "%s"
	}`, query)

}











