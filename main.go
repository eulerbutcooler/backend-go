package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
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

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement logic
		w.Header().Set("X-Custom-Header", "Pav bhaji ka kya bhav paaji")
		// End of middleware logic
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, time.Since(start))
		next.ServeHTTP(w, r)
	})
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home sweet home")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "about last night")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", logMiddleware(headerMiddleware(http.HandlerFunc(homeHandler))))
	mux.Handle("/about", logMiddleware(headerMiddleware(http.HandlerFunc(aboutHandler))))

	// localhost:8080/api -> called -> handler -> function
	mux.HandleFunc("/api", apiHandler)
	log.Printf("Starting server at port %d", port)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("khel khatam", err)
	}
}
