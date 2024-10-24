package database

import "time"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Fullname  string    `json:"fullname"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) CountUser(username string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := s.db.QueryRow(query, username).Scan(&count)
	return count, err
}

func (s *service) InsertUserByUsernameAndPassword(username, hashedPassword string) (string, error) {
	// Insert the new user into the database
	fullName := username
	role := "USER"
	var userID string
	err := s.db.QueryRow("INSERT INTO users (username, password, fullname, role) VALUES ($1, $2, $3, $4) RETURNING id", username, hashedPassword, fullName, role).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// GetHashedPassword retrieves the hashed password for a given username from the database
func (s *service) GetHashedPassword(username string) (string, string, error) {
	var hashedPassword string
	var userID string
	err := s.db.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", "", err
	}
	return userID, hashedPassword, nil
}

func (s *service) GetUserById(userId string) (*User, error) {
	var user User
	err := s.db.QueryRow("SELECT id, username, fullname,role,created_at,updated_at FROM users WHERE id = $1", userId).Scan(&user.Id, &user.Username, &user.Fullname, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
