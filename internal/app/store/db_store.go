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
	maxUrlLength     int
	connectionString string
	db               *sql.DB
}

func (store *DbStore) Open() error {
	db, err := sql.Open("postgres", store.connectionString)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	store.db = db
	logrus.Info("Database connection set")
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
	if len(url) > store.maxUrlLength {
		return "", errors.New(fmt.Sprintf("maximum URL length is %s", url))
	}

	hash := algorithm.ComputeHash(url)

	row := store.db.QueryRow("INSERT INTO Recordings (original, shortened) VALUES($1, $2) ON CONFLICT DO NOTHING",
		url, hash)

	if row.Err() != nil {
		return "", row.Err()
	}

	return hash, nil
}

func (store *DbStore) RestoreLink(hash string) (string, error) {
	var original string

	err := store.db.QueryRow("SELECT original FROM Recordings WHERE shortened = $1", hash).Scan(&original)

	if err != nil {
		return "", errors.New(fmt.Sprintf("shortening %s doesn't exist", hash))
	}

	return original, nil
}
