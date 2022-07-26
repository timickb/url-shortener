package store

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/url-shortener/internal/app/algorithm"
)

func TestNewStore(t *testing.T) {
	st1, err1 := New(nil, nil, logrus.StandardLogger(), "local")
	st2, err2 := New(nil, nil, logrus.StandardLogger(), "db")
	st3, err3 := New(nil, nil, logrus.StandardLogger(), "test")

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("error occured while creating store instances")
	}

	assert.Equal(t, reflect.TypeOf(&LocalStore{}), reflect.TypeOf(st1))
	assert.Equal(t, reflect.TypeOf(&DbStore{}), reflect.TypeOf(st2))
	assert.Equal(t, reflect.TypeOf(&StubStore{}), reflect.TypeOf(st3))
}

func TestLocalStoreRestoreLink(t *testing.T) {
	shr := algorithm.StubShortener{}
	st, err := NewLocal(&shr, logrus.StandardLogger())

	if err != nil {
		t.Fatalf("error occurred while creating local store instance")
	}

	st.Open()
	link := "https://abc.xyz"
	hash := shr.ComputeShortening(link)

	_, _ = st.CreateLink(link)
	result, _ := st.RestoreLink(hash)

	assert.Equal(t, link, result)

}

func TestLocalStoreCreateLink(t *testing.T) {
	shr := algorithm.StubShortener{}
	st, err := NewLocal(&shr, logrus.StandardLogger())

	if err != nil {
		t.Fatalf("error occurred while creating local store instance")
	}

	st.Open()
	link := "https://abc.xyz"
	hash := shr.ComputeShortening(link)

	result, _ := st.CreateLink(link)

	assert.Equal(t, hash, result)

}

func TestDbStoreRestoreLinkDoesntExist(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	s, _ := NewDB(algorithm.DefaultShortener{HashSize: 10}, logrus.StandardLogger(), db)

	mock.ExpectQuery("SELECT original FROM recordings")

	if _, err := s.RestoreLink("some-hash"); err == nil {
		t.Errorf("error was expected while restoring link")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDbStoreCreateLink(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	shr := algorithm.DefaultShortener{HashSize: 10}

	s, _ := NewDB(shr, logrus.StandardLogger(), db)

	hash := shr.ComputeShortening("test-link")

	mock.ExpectExec("INSERT INTO recordings").
		WithArgs("test-link", hash).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err = s.CreateLink("test-link"); err != nil {
		t.Errorf("error was not expected when call CreateLink: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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
	s, _ := NewImproved(shr, logrus.StandardLogger(), db)

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
	s, _ := NewImproved(shr, logrus.StandardLogger(), db)

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
