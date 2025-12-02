package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var port = 8080

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users = make(map[int]User)
	idSeq = 1
	mutex = &sync.Mutex{}
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		mutex.Lock()
		defer mutex.Unlock()
		usersList := make([]User, 0, len(users))
		for _, user := range users {
			usersList = append(usersList, user)
		}
		json.NewEncoder(w).Encode(usersList)
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadGateway)
			return
		}
		mutex.Lock()
		user.ID = idSeq
		idSeq++
		users[user.ID] = user // This line stores the newly created user in the users map using the generated ID as the key.
		mutex.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var id int
	if _, err := fmt.Sscanf(r.URL.Path, "/users/%d", &id); err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	user, ok := users[id]
	if !ok {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(user)
	case "PUT":
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		updatedUser.ID = id
		users[id] = updatedUser
		json.NewEncoder(w).Encode(updatedUser)
	case "DELETE":
		delete(users, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

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

// Understanding query parameters
// now if we get a request on url - localhost:8080/?name=amaan then
// it will print home sweet home amaan
func homeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "home sweet home %s", name)
}

// Extracting path variables
// Example - https://localhost:8080/about/123
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) >= 3 && pathSegments[1] == "about" {
		userID := pathSegments[2]
		fmt.Fprintf(w, "User ID: %s", userID)
	} else {
		http.NotFound(w, r)
	}
	// fmt.Fprintln(w, "about last night")
}

// Combining both query params and path variables
// http://localhost:8080/username/123?includedetails=true
func usernameHandler(w http.ResponseWriter, r *http.Request) {
	pathSeg := strings.Split(r.URL.Path, "/")
	includeDets := r.URL.Query().Get("includedetails")
	if len(pathSeg) >= 3 && pathSeg[1] == "username" {
		userId := pathSeg[2]
		fmt.Fprintf(w, "User id: %s\n", userId)
		if includeDets == "true" {
			w.Write([]byte("Details are included\n"))
		}
	} else {
		http.NotFound(w, r)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", logMiddleware(headerMiddleware(http.HandlerFunc(homeHandler))))
	mux.Handle("/about/", logMiddleware(headerMiddleware(http.HandlerFunc(aboutHandler))))
	mux.Handle("/username/", logMiddleware(headerMiddleware(http.HandlerFunc(usernameHandler))))
	mux.Handle("/users", logMiddleware(http.HandlerFunc(usersHandler)))
	mux.Handle("/users/", logMiddleware(http.HandlerFunc(userHandler)))
	// localhost:8080/api -> called -> handler -> function
	mux.HandleFunc("/api", apiHandler)
	log.Printf("Starting server at port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
