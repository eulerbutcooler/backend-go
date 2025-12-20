package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eulerbutcooler/backend-go/internal/db"
	_ "github.com/lib/pq"
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

func (user User) Validate() error {
	if user.Name == "" {
		return fmt.Errorf("missing field: name")
	}
	if user.Email == "" {
		return fmt.Errorf("missing field: email")
	}
	return nil
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	err = user.Validate()
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User data is valid"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		rows, err := db.GetDB().Query("SELECT id, name, email FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadGateway)
			return
		}
		if err := user.Validate(); err != nil {
			response := map[string]string{"error": err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		err := db.GetDB().QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}
		// mutex.Lock()
		// user.ID = idSeq
		// idSeq++
		// users[user.ID] = user // This line stores the newly created user in the users map using the generated ID as the key.
		// mutex.Unlock()
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
	// mutex.Lock()
	// defer mutex.Unlock()
	// user, ok := users[id]
	// if !ok {
	// 	http.Error(w, "User Not Found", http.StatusNotFound)
	// 	return
	// }
	switch r.Method {
	case "GET":
		var user User
		err := db.GetDB().QueryRow("SELECT * FROM users where id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
		}
		json.NewEncoder(w).Encode(user)
	case "PUT":
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err := db.GetDB().QueryRow("UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING name, email, id", updatedUser.Name, updatedUser.Email, id).Scan(&updatedUser.Name, &updatedUser.Email, &updatedUser.ID)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			log.Printf("Update err: %v", err)
			return
		}
		json.NewEncoder(w).Encode(updatedUser)
	case "DELETE":
		_, err := db.GetDB().Exec("DELETE from users WHERE id=$1", id)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}
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
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database: %w", err)
	}
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
