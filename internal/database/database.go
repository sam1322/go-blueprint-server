package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	GetUrlByKey(urlKey string) (*Url, error)
	AddShortenedUrl(urlResp *Url) error
	CountUser(username string) (int, error)
	InsertUserByUsernameAndPassword(username, hashedPassword string) (string, error)
	UpdateUserImageById(userImage, userId string) error
	GetHashedPassword(username string) (string, string, error)
	GetUserById(userId string) (*User, error)
	GetValidTokenCount(userID string) (int, error)
	InvalidateOldestToken(userID string) error
	InvalidateToken(token string) error
	InsertToken(userID string, token string, expiresAt time.Time) error
	IsTokenValid(token string) (bool, error)
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	//schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	//connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	//Perform a simple "PING" query to verify connection
	//err = db.Ping()
	//if err != nil {
	//	log.Fatalf("Failed to ping database: %v", err)
	//}
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	log.Println("Connected to the database ")
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) Close1() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}
