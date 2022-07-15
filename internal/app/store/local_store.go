package store

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type LocalStore struct {
	maxUrlLength int
	// key -> shortened url, value -> original url
	memoryStorage map[string]string
}

func (store *LocalStore) Open() error {
	store.memoryStorage = make(map[string]string)
	return nil
}

func (store *LocalStore) Close() error {
	return nil
}

func (store *LocalStore) CreateLink(url string) (string, error) {
	if len(url) > store.maxUrlLength {
		return "", errors.New(fmt.Sprintf("maximum URL length is %d", store.maxUrlLength))
	}

	hash := algorithm.ComputeShortening(url)

	value, ok := store.memoryStorage[hash]

	// If hash exists and its original is different -> collision
	if ok && url != value {
		return "", errors.New(fmt.Sprintf("collision happened with url %s", value))
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

	return "", errors.New(fmt.Sprintf("shortening %s doesn't exist", hash))
}
