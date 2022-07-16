package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/timickb/url-shortener/internal/app/algorithm"
	"testing"
)

func TestDbStore_RestoreLink_DoesntExist(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewDBStore(db)

	mock.ExpectQuery("SELECT original FROM recordings")

	if _, err := s.RestoreLink("some-hash"); err == nil {
		t.Errorf("error was expected while restoring link")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDbStore_CreateLink(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewDBStore(db)

	hash := algorithm.ComputeShortening("test-link")

	mock.ExpectExec("INSERT INTO recordings").
		WithArgs("test-link", hash).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err = s.CreateLink("test-link"); err != nil {
		t.Errorf("error was not expected while creating link: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
