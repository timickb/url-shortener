package store

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

type Store interface {
	Open() error
	Close() error

	// CreateLink Creates new URL shorting recording
	CreateLink(string) (string, error)

	// RestoreLink Restores original URL by given shorting
	RestoreLink(string) (string, error)
}

func New(options ...func(store *ImprovedStore)) *ImprovedStore {
	store := &ImprovedStore{}

	for _, applyOpt := range options {
		applyOpt(store)
	}

	if store.logger == nil {
		store.logger = logrus.StandardLogger()
	}

	if store.shr == nil {
		store.shr = algorithm.DefaultShortener{HashSize: 10}
	}

	if store.maxURLLen <= 0 {
		store.maxURLLen = 500
	}

	return store
}

func WithLogger(logger *logrus.Logger) func(*ImprovedStore) {
	return func(store *ImprovedStore) {
		store.logger = logger
	}
}

func WithDB(db *sql.DB) func(*ImprovedStore) {
	return func(store *ImprovedStore) {
		store.db = db
	}
}

func WithShortener(shr algorithm.Shortener) func(*ImprovedStore) {
	return func(store *ImprovedStore) {
		store.shr = shr
	}
}

func WithMaxURLLen(maxLen int) func(*ImprovedStore) {
	return func(store *ImprovedStore) {
		store.maxURLLen = maxLen
	}
}
