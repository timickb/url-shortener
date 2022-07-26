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

func NewLocal(logger *logrus.Logger) (*LocalStore, error) {
	return &LocalStore{
		logger:       logger,
		maxUrlLength: 500,
	}, nil
}

func NewDB(logger *logrus.Logger, db *sql.DB) (*DbStore, error) {
	return &DbStore{
		logger:       logger,
		db:           db,
		maxUrlLength: 500,
	}, nil
}

func NewImproved(logger *logrus.Logger, db *sql.DB) (*ImprovedStore, error) {
	return &ImprovedStore{
		logger:       logger,
		db:           db,
		maxUrlLength: 500,
	}, nil
}

func New(db *sql.DB, logger *logrus.Logger, storeImpl string, maxUrlLength int) (Store, error) {
	switch storeImpl {

	case "local":
		return NewLocal(logger)

	case "db":
		return NewDB(logger, db)

	case "improved":
		return NewImproved(logger, db)

	case "test":
		return &StubStore{}, nil

	default:
		return nil, errors.New("store: incorrect impl parameter")
	}
}
