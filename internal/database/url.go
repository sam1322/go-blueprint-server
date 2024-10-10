package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Url struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	LongUrl string `json:"longUrl"`
}

func (s *service) GetUrlByKey(urlKey string) (*Url, error) {
	row := s.db.QueryRowContext(context.Background(), `
		SELECT id, url_key, long_url
		FROM urls
		WHERE url_key = $1`, urlKey)

	var urlResp Url
	err := row.Scan(&urlResp.Id, &urlResp.Key, &urlResp.LongUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("video not found")
		}
		return nil, err
	}

	return &urlResp, nil
}

// AddShortenedUrl - Add shortened url into the database.
func (s *service) AddShortenedUrl(urlResp *Url) error {
	_, err := s.db.ExecContext(context.Background(), "INSERT INTO urls (url_key,long_url) VALUES ($1, $2)",
		urlResp.Key, urlResp.LongUrl)
	return err
}
