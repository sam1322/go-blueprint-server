package server

import (
	"encoding/json"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
	"new_project/internal/response"

	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/coder/websocket"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(response.JSONErrorMiddleware)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(response.JSONErrorMiddleware(s.logger))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Route does not exist",
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method is not valid",
		})
	})

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	//r.Post("/login", s.Login)
	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		//r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(s.authenticator())

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			message := fmt.Sprintf("protected area. hi %v", claims["user_id"])

			resp := make(map[string]string)
			resp["message"] = message
			err := response.JSON(w, http.StatusOK, resp)
			if err != nil {
				s.serverError(w, r, err)
			}
		})

		r.Route("/api/p/v1", func(r chi.Router) {
			r.Get("/logout", s.Logout)
			r.Get("/user", s.GetUserDetailsByUserId)

			r.Post("/workspace", s.AddWorkspace)
			r.Get("/workspace", s.GetAllWorkspace)
			r.Get("/workspace/{workspaceId}", s.GetWorkspaceById)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/url", s.GetShortenedUrl)
		r.Post("/login", s.NewLogin)
		r.Post("/register", s.Register)

		r.Post("/animal", s.AddAnimals)
		r.Get("/animal/{id}", s.GetAnimalsById)
		r.Get("/animal", s.GetAllAnimals)
	})

	//r.Route("/api", func(r chi.Router) {
	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/{provider}", s.beginAuthProvideCallback)
	r.Get("/logout/{provider}", s.logOutProvider)
	//})

	r.Get("/short/{shortKey}", s.HandleRedirect)

	r.Get("/websocket", s.websocketHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

// TODO: deprecated function
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here you should verify the username and password against your database
	// For this example, we'll just check for a hardcoded value
	if req.Username != "admin" || req.Password != "123456" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the token
	claims := map[string]interface{}{
		"user_id": req.Username,
	}
	//jwtauth.SetExpiryIn(claims, time.Minute*15)
	jwtauth.SetExpiryIn(claims, time.Hour*15)
	jwtauth.SetIssuedNow(claims)

	_, tokenString, _ := tokenAuth.Encode(claims)

	// Return the token
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	err := response.JSON(w, http.StatusOK, map[string]string{"token": tokenString})
	if err != nil {
		s.serverError(w, r, err)
	}
}
