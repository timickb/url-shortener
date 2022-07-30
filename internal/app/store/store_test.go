package store

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

func TestNewStoreDefault(t *testing.T) {
	// create default configured store
	st := New()

	// db must be nil
	assert.Equal(t, (*sql.DB)(nil), st.db)

	// shr must be DefaultShortener
	assert.Equal(t, reflect.TypeOf(algorithm.DefaultShortener{}), reflect.TypeOf(st.shr))

	// logger must be logrus.Logger
	assert.Equal(t, reflect.TypeOf(&logrus.Logger{}), reflect.TypeOf(st.logger))
}

func TestNewStoreWithDB(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}

	st := New(
		WithDB(db),
	)

	assert.Equal(t, db, st.db)
}

func TestNewStoreWithLogger(t *testing.T) {
	logger := logrus.StandardLogger()

	st := New(
		WithLogger(logger),
	)

	assert.Equal(t, logger, st.logger)
}

func TestNewStoreWithShortener(t *testing.T) {
	shr := &algorithm.DefaultShortener{}

	st := New(
		WithShortener(shr),
	)

	assert.Equal(t, shr, st.shr)
}

func TestImprovedStoreCreateLinkFirstTime(t *testing.T) {
	// create mock database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error %s was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	shr := algorithm.DefaultShortener{HashSize: 10}

	// create ImprovedStore instance
	s := New(
		WithDB(db),
		WithLogger(logrus.StandardLogger()),
		WithShortener(shr),
	)

	url := "test-url"
	hash := shr.ComputeShortening(url)

	// expecting insert query
	mock.ExpectExec("INSERT INTO recordings").
		WithArgs(url, hash).WillReturnResult(sqlmock.NewResult(1, 1))

	// opening store
	err = s.Open()
	if err != nil {
		t.Errorf("error was not expected when call Open: %s", err)
	}
	defer s.Close()

	// testing CreateLink method
	if _, err := s.CreateLink(url); err != nil {
		t.Errorf("error was not expected when call CreateLink: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// check whether the url added to in-memory
	_, ok := s.memoryStorage[hash]
	assert.Equal(t, true, ok)
}

func TestImprovedStoreCreateLinkSecondTime(t *testing.T) {
	// create mock database
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error %s was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	shr := algorithm.DefaultShortener{HashSize: 10}

	// create ImprovedStore instance
	s := New(
		WithDB(db),
		WithLogger(logrus.StandardLogger()),
		WithShortener(shr),
	)

	url := "test-url"
	hash := shr.ComputeShortening(url)

	// NOT expecting insert query
	mock.ExpectExec("INSERT INTO recordings").
		WithArgs(url, hash).WillReturnResult(sqlmock.NewResult(1, 1))

	// opening store
	err = s.Open()
	if err != nil {
		t.Errorf("error was not expected when call Open: %s", err)
	}
	defer s.Close()

	// manually add existing url
	s.memoryStorage[hash] = url

	// testing CreateLink method
	if _, err := s.CreateLink(url); err != nil {
		t.Errorf("error was not expected when call CreateLink: %s", err)
	}

	// we don't want the query executed
	if err = mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were fulfilled expectations: %s", err)
	}
}
