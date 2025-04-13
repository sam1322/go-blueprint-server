package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"new_project/internal/database"
	"new_project/internal/models"
	"new_project/internal/response"
	"strconv"
)

type animal *models.Animal

func (s *Server) AddAnimals(w http.ResponseWriter, r *http.Request) {

	var req *database.Animal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input (you may want to add more checks here)
	if req.Name == "" || req.FunFact == "" || req.Description == "" {
		http.Error(w, "Name and Fun fact are required", http.StatusBadRequest)
		return
	}

	err := s.db.InsertAnimal(req)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		s.badRequest(w, r, err)
		return

	}

	err = response.JSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "successfully added new animal"})
	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) GetAnimalsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Validate input (you may want to add more checks here)
	if id == "" {
		http.Error(w, "Valid Id is required", http.StatusBadRequest)
		return
	}
	animal, err := s.db.GetAnimalById(id)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		s.badRequest(w, r, err)
		return
	}
	err = response.JSON(w, http.StatusOK, animal)
	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) GetAllAnimals(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	if page == "" || pageSize == "" {
		http.Error(w, "Valid Page number and Page size is required", http.StatusBadRequest)
		return
	}

	page1, err := strconv.Atoi(page)
	if err != nil {
		s.badRequest(w, r, err)
		return
	}
	pageSize1, err := strconv.Atoi(pageSize)
	if err != nil {
		s.badRequest(w, r, err)
		return
	}

	animalResp, err := s.db.GetAllAnimals(page1, pageSize1)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		s.badRequest(w, r, err)
		return
	}
	err = response.JSON(w, http.StatusOK, animalResp)
	if err != nil {
		s.serverError(w, r, err)
	}
}
