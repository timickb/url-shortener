package store

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type LocalStore struct {
	logger        *logrus.Logger
	memoryStorage map[string]string
}

func (s *LocalStore) Open() error {
	s.memoryStorage = make(map[string]string)
	s.logger.Info("In-memory storage initialized")
	return nil
}

func (store *LocalStore) Close() error {
	return nil
}

func (store *LocalStore) CreateLink(url string) (string, error) {
	hash := algorithm.ComputeShortening(url)

	value, ok := store.memoryStorage[hash]

	// If hash exists and its original is different -> collision
	if ok && url != value {
		return "", fmt.Errorf("collision happened with url %s", value)
	}

	// If hash exists and its original is same -> do nothing

	if !ok {
		store.memoryStorage[hash] = url
	}

	return hash, nil
}

func (store *LocalStore) RestoreLink(hash string) (string, error) {
	value, ok := store.memoryStorage[hash]

	if ok {
		return value, nil
	}

	return "", fmt.Errorf("shortening %s doesn't exist", hash)
}
