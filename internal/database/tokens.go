package database

import (
	"time"
)

// InsertToken stores a token in the database for the user
func (s *service) InsertToken(userID int, token string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userID, token, expiresAt)
	return err
}

// InvalidateOldestToken invalidates the oldest token if the user has more than 2 valid tokens
func (s *service) InvalidateOldestToken(userID int) error {
	_, err := s.db.Exec(`
		UPDATE tokens
		SET is_valid = FALSE
		WHERE id = (
			SELECT id FROM tokens
			WHERE user_id = $1 AND is_valid = TRUE
			ORDER BY created_at ASC
			LIMIT 1
		)
	`, userID)
	//,  updated_at = CURRENT_TIMESTAMP
	return err
}

// InvalidateOldestToken invalidates the oldest token if the user has more than 2 valid tokens
func (s *service) InvalidateToken(token string) error {
	_, err := s.db.Exec(`
		UPDATE tokens
		SET is_valid = FALSE
		WHERE token=$1
	`, token)
	return err
}

// GetValidTokenCount returns the number of valid tokens for a user
func (s *service) GetValidTokenCount(userID int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tokens WHERE user_id = $1 AND is_valid = TRUE", userID).Scan(&count)
	return count, err
}

func (s *service) IsTokenValid(token string) (bool, error) {
	var valid bool
	//err := s.db.QueryRow("SELECT EXISTS(SELECT is_valid FROM tokens WHERE token = $1 AND expires_at > NOW())", token).Scan(&valid)
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tokens WHERE token = $1 AND is_valid=true)", token).Scan(&valid)
	return valid, err
}
