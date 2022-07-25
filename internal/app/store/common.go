package store

import (
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

type Store interface {
	Open() error
	Close() error

	// CreateLink Creates new URL shorting recording
	CreateLink(string) (string, error)

	// RestoreLink Restores original URL by given shorting
	RestoreLink(string) (string, error)
}

func New(db *sql.DB, logger *logrus.Logger, storeImpl string, maxUrlLength int) (Store, error) {
	switch storeImpl {

	case "local":
		return &LocalStore{maxUrlLength: maxUrlLength}, nil

	case "db":
		return &DbStore{
			db:           db,
			maxUrlLength: maxUrlLength}, nil

	case "improved":
		return &ImprovedStore{
			db:           db,
			maxUrlLength: 300,
			logger:       logger,
		}, nil

	case "test":
		return &MockStore{}, nil

	default:
		return nil, errors.New("store: incorrect impl parameter")
	}
}

func NewDBStore(db *sql.DB) *DbStore {
	return &DbStore{
		db:           db,
		maxUrlLength: 300,
	}
}

func NewImprovedStore(db *sql.DB) *ImprovedStore {
	return &ImprovedStore{
		db:           db,
		maxUrlLength: 300,
	}
}
