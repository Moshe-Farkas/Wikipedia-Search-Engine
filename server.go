package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	fmt.Println("hello, world")
	http.Handle("/", http.FileServer(http.Dir("./files")))
	log.Fatal(http.ListenAndServe(":1234", nil))
}








// package main

// import (
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// var globalRequestCount int

// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", Root)
// 	mux.HandleFunc("/hello", Hello)
// 	mux.HandleFunc("/headers", headers)
// 	mux.HandleFunc("/requests", requestCount)
// 	mux.HandleFunc("/jsontest", jsonTest)
// 	err := http.ListenAndServe(":1234", mux)
// 	if errors.Is(err, http.ErrServerClosed) {
// 		fmt.Printf("server closed\n")
// 	} else if err != nil {
// 		fmt.Printf("could not start server: %s\n", err)	
// 		os.Exit(1)
// 	}
// }

// func jsonTest(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintln(w, `{
// 	"person": {
// 		"name": "John Doe",
// 		"age": 30,
// 		"email": "john@example.com",
// 		"address": {
// 		"street": "123 Main St",
// 		"city": "Anytown",
// 		"state": "CA",
// 		"postal_code": "12345"
// 		},
// 		"hobbies": ["reading", "hiking", "cooking"],
// 		"is_student": false
// 		}
// 	} `)
// }

// func requestCount(w http.ResponseWriter, r *http.Request) {
// 	globalRequestCount++
// 	fmt.Fprintf(w, "you are the %d request\n", globalRequestCount)	
// }

// func headers(w http.ResponseWriter, r *http.Request) {
// 	for name, headers := range r.Header {
// 		for _, h := range headers {
// 			fmt.Fprintf(w, "%v: %v\n", name, h)
// 		}
// 	}
// 	globalRequestCount++
// }

// func Root(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "this is the root\n")
// 	globalRequestCount++
// }

// func Hello(w http.ResponseWriter, r *http.Request) {
// 	// ctx := r.Context()
// 	// fmt.Println("function handler")
// 	// defer fmt.Println("function handler ended")
	
// }









