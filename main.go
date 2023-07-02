package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
)

// Define the Casbin enforcer as a global variable
// var enforcer *casbin.Enforcer

// Define the user struct
type User struct {
	ID       string
	Username string
	Role     string
}

// Sample user data
var users = []User{
	{ID: "1", Username: "admin", Role: "admin"},
	{ID: "2", Username: "user1", Role: "user"},
	{ID: "3", Username: "user2", Role: "user"},
}

func main() {
	// Initialize Casbin enforcer
	enforcer, err := initializeCasbinEnforcer()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Chi router
	r := chi.NewRouter()

	// Public routes accessible by all users (read-only)
	r.Group(func(r chi.Router) {
		// Apply Casbin middleware to protect the routes
		r.Use(casbinMiddleware("user", enforcer))

		// Public routes for read-only access to users
		r.Get("/users", getAllUsers)
		r.Get("/users/{id}", getUserByID)
	})

	// Protected routes accessible only by admins
	r.Group(func(r chi.Router) {
		// Apply Casbin middleware to protect the routes
		r.Use(casbinMiddleware("admin", enforcer))

		// CRUD operations for users
		r.Get("/users", getAllUsers)
		r.Get("/users/{id}", getUserByID)
		r.Post("/users", createUser)
		r.Put("/users/{id}", updateUser)
		r.Delete("/users/{id}", deleteUser)
	})

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}

// Initialize Casbin enforcer
func initializeCasbinEnforcer() (*casbin.Enforcer, error) {
	modelPath := "./model.conf"
	policyPath := "./policy.csv"
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}

// Casbin middleware to enforce role-based access control
func casbinMiddleware(role string, enforcer *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Perform authorization check
			subject := role // Assuming the role is the subject in Casbin
			resource := r.URL.Path
			action := r.Method

			allowed, err := enforcer.Enforce(subject, resource, action)
			if err != nil {
				http.Error(w, "InternalServerError", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Handler functions
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve all users (dummy data)
	fmt.Fprintf(w, "List of all users:\n")
	for _, user := range users {
		fmt.Fprintf(w, "ID: %s, Username: %s, Role: %s\n", user.ID, user.Username, user.Role)
	}
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user by ID (dummy data)
	id := chi.URLParam(r, "id")
	for _, user := range users {
		if user.ID == id {
			fmt.Fprintf(w, "User found:\n")
			fmt.Fprintf(w, "ID: %s, Username: %s, Role: %s\n", user.ID, user.Username, user.Role)
			return
		}
	}

	http.NotFound(w, r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Create a new user (dummy data)
	// ...

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Update a user (dummy data)
	// ...

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Delete a user (dummy data)
	// ...

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
