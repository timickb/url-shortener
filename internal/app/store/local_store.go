package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/timickb/url-shortener/internal/app/algorithm"
	"time"

	_ "github.com/lib/pq"
)

type Recording struct {
	Original  string
	Shortened string
	Created   time.Time
}

type LocalStore struct {
	config *Config
	db     *sql.DB
	local  map[string]Recording
}

func (store *LocalStore) Open() error {
	store.local = make(map[string]Recording)
	return nil
}

func (store *LocalStore) Close() error {
	return nil
}

func (store *LocalStore) CreateLink(url string) (string, error) {
	hash := algorithm.ComputeHash(url)
	value, ok := store.local[hash]

	if ok {
		return "", errors.New(
			fmt.Sprintf("local_store: such shortening already exists for url %s", value))
	}

	store.local[hash] = Recording{
		Original:  url,
		Shortened: hash,
		Created:   time.Now(),
	}

	return hash, nil
}

func (store *LocalStore) RestoreLink(hash string) (string, error) {
	value, ok := store.local[hash]

	if ok {
		return value.Original, nil
	}

	return "", errors.New("local_store: shortening doesn't exist")
}
