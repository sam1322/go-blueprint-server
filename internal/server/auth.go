package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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
	userID, err := s.db.InsertUserByUsernameAndPassword(req.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	//claims := map[string]interface{}{
	//	"user_id": req.Username, // or a real user ID
	//}

	// Set expiry and issued time for the token
	//jwtauth.SetExpiryIn(claims, time.Hour*15) // Token valid for 15 hours
	//jwtauth.SetIssuedNow(claims)
	//
	//// Encode the claims into a signed JWT token
	//_, tokenString, err := tokenAuth.Encode(claims)
	//if err != nil {
	//	http.Error(w, "Error generating token", http.StatusInternalServerError)
	//	return
	//}

	fmt.Println("userID ", userID)
	tokenString, err := s.createToken(w, userID)
	fmt.Println("token", tokenString)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}
	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "token": tokenString})
	//err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
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
	userID, hashedPassword, err := s.db.GetHashedPassword(req.Username)
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
	//claims := map[string]interface{}{
	//	"user_id": req.Username, // or a real user ID
	//}
	//
	//// Set expiry and issued time for the token
	//jwtauth.SetExpiryIn(claims, time.Hour*15) // Token valid for 15 hours
	//jwtauth.SetIssuedNow(claims)
	//
	//// Encode the claims into a signed JWT token
	//_, tokenString, err := tokenAuth.Encode(claims)
	//if err != nil {
	//	http.Error(w, "Error generating token", http.StatusInternalServerError)
	//	return
	//}
	//
	//// Check if user already has 2 valid tokens
	//tokenCount, err := s.db.GetValidTokenCount(userID)
	//if err != nil {
	//	http.Error(w, "Error checking token count", http.StatusInternalServerError)
	//	return
	//}
	//
	//// If user has 2 tokens, invalidate the oldest one
	//if tokenCount >= 2 {
	//	err = s.db.InvalidateOldestToken(userID)
	//	if err != nil {
	//		http.Error(w, "Error invalidating oldest token", http.StatusInternalServerError)
	//		return
	//	}
	//}
	//
	//// Insert the new token into the database
	//expiresAt := time.Now().Add(time.Hour * 15)
	//err = s.db.InsertToken(userID, tokenString, expiresAt)
	fmt.Println("userId ", userID)
	tokenString, err := s.createToken(w, userID)
	fmt.Println("token ", tokenString)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
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

// Logout handles the user logout process
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := extractTokenFromHeader(r)

	err := s.db.InvalidateToken(tokenString)

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Logout is successful"})
	if err != nil {
		http.Error(w, "Error while logging out", http.StatusBadRequest)
	}

}

func (s *Server) createToken(w http.ResponseWriter, userId int) (string, error) {
	// Create JWT claims (you can add more claims as necessary)
	claims := map[string]interface{}{
		//"user_id": username, // or a real user ID
		"user_id": userId,
	}

	// Set expiry and issued time for the token
	jwtauth.SetExpiryIn(claims, time.Hour*15) // Token valid for 15 hours
	jwtauth.SetIssuedNow(claims)

	// Encode the claims into a signed JWT token
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return "", err
	}

	// Check if user already has 2 valid tokens
	tokenCount, err := s.db.GetValidTokenCount(userId)
	if err != nil {
		http.Error(w, "Error checking token count", http.StatusInternalServerError)
		return "", err
	}

	// If user has 2 tokens, invalidate the oldest one
	if tokenCount >= 2 {
		err = s.db.InvalidateOldestToken(userId)
		if err != nil {
			http.Error(w, "Error invalidating oldest token", http.StatusInternalServerError)
			return "", err
		}
	}

	// Insert the new token into the database
	expiresAt := time.Now().Add(time.Hour * 15)
	err = s.db.InsertToken(userId, tokenString, expiresAt)
	if err != nil {
		http.Error(w, "Error saving token", http.StatusInternalServerError)
		return "", err
	}

	return tokenString, nil

}

// r.Use(s.authenticator(tokenAuth))
//func (s *Server) authenticator(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
//}

func (s *Server) authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Extract token string from Authorization header
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				http.Error(w, "Missing or invalid authorization token", http.StatusUnauthorized)
				return
			}

			// Check if the token is valid in the database
			valid, err := s.db.IsTokenValid(tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error validating token %v", err), http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
func extractTokenFromHeader(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
