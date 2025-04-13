package database

import (
	"fmt"
	"new_project/internal/models"
)

// AddWorkspace inserts a new workspace into the database.
func (s *service) AddWorkspace(userID string, name string, joinCode string) error {
	_, err := s.db.Exec(`
		INSERT INTO workspace (name, join_code, user_id)
		VALUES ($1, $2, $3)`,
		name, joinCode, userID,
	)
	return err
}

//type Workspace models.Workspace

// GetWorkspaces retrieves all workspaces for a given user_id.
func (s *service) GetWorkspaces(userId string) ([]models.Workspace, error) {
	rows, err := s.db.Query(`
		SELECT id, name, join_code, user_id, created_at, updated_at
		FROM workspace
		WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//var workspaces []models.Workspace
	workspaces := []models.Workspace{}

	for rows.Next() {
		var ws models.Workspace
		if err := rows.Scan(&ws.Id, &ws.Name, &ws.JoinCode, &ws.UserId, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

// GetWorkspaces retrieves all workspaces for a given user_id.
func (s *service) GetWorkspacesById(userId, workspaceId string) (*models.Workspace, error) {
	rows, err := s.db.Query(`
		SELECT id, name, join_code, user_id, created_at, updated_at
		FROM workspace
		WHERE user_id = $1 and id = $2`, userId, workspaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//var workspaces []models.Workspace
	//workspaces := models.Workspace{}

	//for rows.Next() {
	var ws models.Workspace
	// Use rows.Next() to move to the next row
	if rows.Next() {
		if err := rows.Scan(&ws.Id, &ws.Name, &ws.JoinCode, &ws.UserId, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, err
		}
	} else {
		// Return nil and an error if no rows were found
		return nil, fmt.Errorf("no workspace found with user_id: %s and id: %s", userId, workspaceId)
	}

	//workspaces = append(workspaces, ws)
	//}
	return &ws, nil
}
