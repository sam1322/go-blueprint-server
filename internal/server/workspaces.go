package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"new_project/internal/models"
	"new_project/internal/response"
)

func (s *Server) AddWorkspace(w http.ResponseWriter, r *http.Request) {

	var req *models.Workspace
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.badRequest(w, r, err)
		return
	}

	// Validate input (you may want to add more checks here)
	if req.Name == "" {
		//http.Error(w, "Name is required", http.StatusBadRequest)
		s.badRequest(w, r, fmt.Errorf("name is required"))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	//userId := fmt.Sprintf("protected area. hi %v", claims["user_id"])
	userId := claims["user_id"].(string)

	joinCode := generateShortKey(8)

	err := s.db.AddWorkspace(userId, req.Name, joinCode)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		s.badRequest(w, r, err)
		return

	}

	err = response.JSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "successfully added new workspace"})
	if err != nil {
		s.serverError(w, r, err)
	}
}

type WorkspaceResp struct {
	Workspace []models.Workspace `json:"workspace"`
}

func (s *Server) GetAllWorkspace(w http.ResponseWriter, r *http.Request) {

	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["user_id"].(string)

	var workspaceResp WorkspaceResp

	workspaces, err := s.db.GetWorkspaces(userId)
	if err != nil {
		s.badRequest(w, r, err)
		return

	}

	workspaceResp.Workspace = workspaces

	err = response.JSON(w, http.StatusOK, workspaceResp)
	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) GetWorkspaceById(w http.ResponseWriter, r *http.Request) {

	workspaceId := chi.URLParam(r, "workspaceId")

	// Validate input (you may want to add more checks here)
	if workspaceId == "" {
		//http.Error(w, "Name is required", http.StatusBadRequest)
		s.badRequest(w, r, fmt.Errorf("id is required"))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["user_id"].(string)

	workspace, err := s.db.GetWorkspacesById(userId, workspaceId)
	if err != nil {
		s.badRequest(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, workspace)
	if err != nil {
		s.serverError(w, r, err)
	}
}
