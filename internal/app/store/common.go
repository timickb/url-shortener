package store

import (
	"database/sql"
	"errors"
	"fmt"
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

	case "improved":
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			fmt.Println("couldn't create database connection")
		}
		return &ImprovedStore{
			db:           db,
			maxUrlLength: 300,
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
