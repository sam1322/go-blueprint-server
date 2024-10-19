package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// RegisterRequest represents the data needed for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles the user registration process
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input (you may want to add more checks here)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check if the username already exists in the database
	count, err := s.db.CountUser(req.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking username availability: %v", err), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error hashing password: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	err = s.db.InsertUserByUsernameAndPassword(req.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewLogin handles the user login process
func (s *Server) NewLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input (you may want to add more checks here)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve the stored hashed password for the username
	hashedPassword, err := s.db.GetHashedPassword(req.Username)
	if err != nil {
		//http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("User not found"), http.StatusNotFound)
		return
	}

	// If no user found
	if hashedPassword == "" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		// Password does not match
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Password matches, login is successful
	w.WriteHeader(http.StatusOK)
	// Create JWT claims (you can add more claims as necessary)
	claims := map[string]interface{}{
		"user_id": req.Username, // or a real user ID
	}

	// Set expiry and issued time for the token
	jwtauth.SetExpiryIn(claims, time.Hour*15) // Token valid for 15 hours
	jwtauth.SetIssuedNow(claims)

	// Encode the claims into a signed JWT token
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

	//
	//err = json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
}
