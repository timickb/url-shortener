package store

import (
	"errors"
	"io"
)

type Store interface {
	io.Closer

	Open() error

	// CreateLink Creates new URL shorting recording
	CreateLink(string) (string, error)

	// RestoreLink Restores original URL by given shorting
	RestoreLink(string) (string, error)
}

func NewStore(config *Config, storeImpl string) (Store, error) {
	switch storeImpl {
	case "local":
		return &LocalStore{config: config}, nil
	case "db":
		return &DbStore{config: config}, nil
	default:
		return nil, errors.New("store: incorrect impl parameter")
	}
}
