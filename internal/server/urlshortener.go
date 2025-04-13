package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"new_project/internal/database"
	"time"
)

type UrlRequest struct {
	Url string `json:"url"`
}
type UrlResponse struct {
	Key      string `json:"key"`
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

// Url alias for database.Url
type Url = database.Url

func (s *Server) GetShortenedUrl(w http.ResponseWriter, r *http.Request) {
	var request UrlRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Panicf("error handling JSON marshal. Err: %v", err)
	}
	longUrl := request.Url
	shortKey := generateShortKey(8)

	// TODO : add validations for long url

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	urlObj := Url{Key: shortKey, LongUrl: longUrl}
	err = s.db.AddShortenedUrl(&urlObj)

	//err = s.db.AddShortenedUrl(&Url{ // this way of writing struct and dynamically passing address also works
	//	Key:     shortKey,
	//	LongUrl: longUrl,
	//})
	log.Printf("Original url : %v\n", longUrl)
	log.Println("Generating a shortened url " + shortenedURL)

	urlResponse := UrlResponse{
		Key:      shortKey,
		LongUrl:  longUrl,
		ShortUrl: shortenedURL,
	}

	//err = json.NewEncoder(w).Encode(urlResponse)

	jsonResp, err := json.Marshal(urlResponse)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func (s *Server) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
	log.Println("shortkey : ", shortKey)
	urlResponse, err := s.db.GetUrlByKey(shortKey)
	if err != nil {
		log.Panicf("Err: %v", err)
	}
	originalUrl := urlResponse.LongUrl
	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
	//jsonResp, err := json.Marshal(urlResponse)
	//if err != nil {
	//	log.Fatalf("error handling JSON marshal. Err: %v", err)
	//}
	//w.Header().Set("Content-Type", "application/json")
	//_, _ = w.Write(jsonResp)
}

func generateShortKey(keyLength int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//const keyLength = 8

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
