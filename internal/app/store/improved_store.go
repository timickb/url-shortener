package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type ImprovedStore struct {
	memoryStorage map[string]string
	db            *sql.DB
	logger        *logrus.Logger
}

func (s *ImprovedStore) Open() error {
	s.memoryStorage = make(map[string]string)
	s.logger.Info("Memory storage intialized")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.db != nil {
		if err := s.db.PingContext(ctx); err != nil {
			s.logger.Errorf("Couldn't ping database: %s\n", err)
		} else {
			s.logger.Info("Database connection established")
		}
	}

	return nil
}

func (s *ImprovedStore) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}

func (s *ImprovedStore) CreateLink(url string) (string, error) {
	hash := algorithm.ComputeShortening(url)

	_, ok := s.memoryStorage[hash]

	if ok {
		return hash, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO recordings (original, shortened) VALUES($1, $2) ON CONFLICT DO NOTHING",
		url, hash)

	if err != nil {
		return "", nil
	}

	if len(s.memoryStorage) < 100 {
		s.memoryStorage[hash] = url
	}

	return hash, nil
}

func (s *ImprovedStore) RestoreLink(hash string) (string, error) {
	value, ok := s.memoryStorage[hash]

	if ok {
		return value, nil
	}

	// if not found - request from database
	var original string
	err := s.db.QueryRow("SELECT original FROM recordings WHERE shortened = $1", hash).Scan(&original)

	if err != nil {
		return "", fmt.Errorf("shortening %s doesn't exist", hash)
	}

	s.memoryStorage[hash] = original

	return original, nil
}
