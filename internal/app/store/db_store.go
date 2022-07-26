package store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type DbStore struct {
	db     *sql.DB
	logger *logrus.Logger
}

func (s *DbStore) Open() error {
	if s.db == nil {
		return errors.New("database is nil")
	}

	if err := s.db.Ping(); err != nil {
		return err
	}
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
	hash := algorithm.ComputeShortening(url)

	_, err := store.db.Exec("INSERT INTO recordings (original, shortened) VALUES($1, $2) ON CONFLICT DO NOTHING",
		url, hash)

	if err != nil {
		return "", err
	}

	return hash, nil
}

func (store *DbStore) RestoreLink(hash string) (string, error) {
	var original string

	err := store.db.QueryRow("SELECT original FROM recordings WHERE shortened = $1", hash).Scan(&original)

	if err != nil {
		return "", fmt.Errorf("shortening %s doesn't exist", hash)
	}

	return original, nil
}
