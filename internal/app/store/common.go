package store

import (
	"database/sql"
	"errors"

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

func NewLocal(shr algorithm.Shortener, logger *logrus.Logger) (*LocalStore, error) {
	return &LocalStore{
		shr:    shr,
		logger: logger,
	}, nil
}

func NewDB(shr algorithm.Shortener, logger *logrus.Logger, db *sql.DB) (*DbStore, error) {
	return &DbStore{
		shr:    shr,
		logger: logger,
		db:     db,
	}, nil
}

func NewImproved(shr algorithm.Shortener, logger *logrus.Logger, db *sql.DB) (*ImprovedStore, error) {
	return &ImprovedStore{
		shr:    shr,
		logger: logger,
		db:     db,
	}, nil
}

func New(shr algorithm.Shortener, db *sql.DB, logger *logrus.Logger, storeImpl string) (Store, error) {
	switch storeImpl {

	case "local":
		return NewLocal(shr, logger)

	case "db":
		return NewDB(shr, logger, db)

	case "improved":
		return NewImproved(shr, logger, db)

	case "test":
		return &StubStore{}, nil

	default:
		return nil, errors.New("store: incorrect impl parameter")
	}
}
