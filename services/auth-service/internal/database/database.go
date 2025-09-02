package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/redis/go-redis/v9"
)

// Connect establishes a connection to PostgreSQL
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0) // No limit

	return db, nil
}

// ConnectRedis establishes a connection to Redis
func ConnectRedis(redisURL string) (*redis.Client, error) {
	// Parse Redis URL
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Extract password
	password := ""
	if parsedURL.User != nil {
		password, _ = parsedURL.User.Password()
	}

	// Extract database number (default to 0)
	db := 0
	if parsedURL.Path != "" && len(parsedURL.Path) > 1 {
		fmt.Sscanf(parsedURL.Path[1:], "%d", &db)
	}

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     parsedURL.Host,
		Password: password,
		DB:       db,
	})

	return client, nil
}
