package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type DbStore struct {
	config *Config
	db     *sql.DB
	local  map[string]Recording
}

func (store *DbStore) Open() error {
	db, err := sql.Open("postgres", store.config.ConnectionString)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	store.db = db
	fmt.Println("DB connected")
	return nil
}

func (store *DbStore) Close() error {
	err := store.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (store *DbStore) CreateLink(url string) (string, error) {
	hash := algorithm.ComputeHash(url)
	// db
	return hash, nil
}

func (store *DbStore) RestoreLink(hash string) (string, error) {
	return "https://example.com", nil
}
