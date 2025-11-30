package main

import (
	"fmt"
	"net/http"
)

var port = 8080

// *http.Request -> location -> User requests and parameters are present -> user provided data
// http.ResponseWriter -> Backend writes its response
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	// this also does the same thing
	w.Write([]byte("Hello World"))
	// Hello world -> w
	fmt.Fprintf(w, "Hello World")
}

func main() {
	// localhost:8080/api -> called -> handler -> function
	http.HandleFunc("/api", apiHandler)
	fmt.Printf("Starting server at port %d", port)
	http.ListenAndServe(":8080", nil)
}
