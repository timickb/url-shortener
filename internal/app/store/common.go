package store

import (
	"errors"
)

type Store interface {
	Open() error
	Close() error

	// CreateLink Creates new URL shorting recording
	CreateLink(string) (string, error)

	// RestoreLink Restores original URL by given shorting
	RestoreLink(string) (string, error)
}

func NewStore(connectionString string, storeImpl string, maxUrlLength int) (Store, error) {
	switch storeImpl {

	case "local":
		return &LocalStore{maxUrlLength: maxUrlLength}, nil

	case "db":
		return &DbStore{
			connectionString: connectionString,
			maxUrlLength:     maxUrlLength}, nil

	case "test":
		return &TestStore{}, nil

	default:
		return nil, errors.New("store: incorrect impl parameter")
	}
}
