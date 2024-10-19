package database

func (s *service) CountUser(username string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := s.db.QueryRow(query, username).Scan(&count)
	return count, err
}

func (s *service) InsertUserByUsernameAndPassword(username, hashedPassword string) error {
	// Insert the new user into the database
	fullName := username
	role := "USER"
	_, err := s.db.Exec("INSERT INTO users (username, password,fullname,role) VALUES ($1, $2,$3,$4)", username, hashedPassword, fullName, role)
	return err
}

// GetHashedPassword retrieves the hashed password for a given username from the database
func (s *service) GetHashedPassword(username string) (string, error) {
	var hashedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}
