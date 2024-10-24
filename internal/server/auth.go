package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"new_project/internal/response"
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

func (s *Server) createToken(w http.ResponseWriter, userId string) (string, error) {
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

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("error in user authentication %v : %v", provider, err)
		fmt.Fprintln(w, err)
		return
	}
	postJsonBytes, err := json.MarshalIndent(user, "", "    ")
	// postJsonBytes, err := JSONMarshal(postJson, "", "    ")
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Println(string(postJsonBytes))
	//http.Redirect(w, r, "http://localhost:3000/movies/dashboard", http.StatusFound)

	tokenString, err := s.RegisterOrLogin(w, user.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating token: %v", err), http.StatusInternalServerError)
		return
	}

	// Create and set the cookie
	cookie := &http.Cookie{
		Name:   "token",
		Value:  tokenString,
		Path:   "/",
		MaxAge: 86400, // 24 hours in seconds
		//HttpOnly: true,  // Cannot be accessed by JavaScript (more secure)
		HttpOnly: false, // Cannot be accessed by JavaScript (more secure)
		Secure:   true,  // Only sent over HTTPS
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost", // Change this in production
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
}

func (s *Server) beginAuthProvideCallback(w http.ResponseWriter, r *http.Request) {

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) logOutProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.Logout(w, r)
	w.Header().Set("Location", "http://localhost:3000/movies/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) RegisterOrLogin(w http.ResponseWriter, username string) (string, error) {
	count, err := s.db.CountUser(username)
	if err != nil {
		//http.Error(w, fmt.Sprintf("Error checking username availability: %v", err), http.StatusInternalServerError)
		return "", err
	}

	userID := ""
	if count > 0 {
		// signin in the user
		// Retrieve the stored hashed password for the username
		userID, _, err = s.db.GetHashedPassword(username)
		if err != nil {
			return "", err
		}

	} else { // registering the user for the first time

		randomPassword := generateSecureRandomString(12)
		if err != nil {
			return "", err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		// Insert the new user into the database
		userID, err = s.db.InsertUserByUsernameAndPassword(username, string(hashedPassword))
		if err != nil {
			return "", err
		}

	}
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("Invalid username or password")
	}

	// Return success response
	fmt.Println("userID ", userID)
	tokenString, err := s.createToken(w, userID)
	fmt.Println("token", tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Secure method using crypto/rand
func generateSecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func (s *Server) GetUserDetailsByUserId(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	//message := fmt.Sprintf("protected area. hi %v", claims["user_id"])
	userId := claims["user_id"].(string)

	user, err := s.db.GetUserById(userId)
	err = response.JSON(w, http.StatusOK, user)
	if err != nil {
		s.serverError(w, r, err)
	}
}
