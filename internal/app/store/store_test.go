package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/url-shortener/internal/app/algorithm"
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	st1, err1 := NewStore("", "local", 300)
	st2, err2 := NewStore("", "db", 300)
	st3, err3 := NewStore("", "test", 300)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("error occured while creating store instances")
	}

	assert.Equal(t, reflect.TypeOf(&LocalStore{}), reflect.TypeOf(st1))
	assert.Equal(t, reflect.TypeOf(&DbStore{}), reflect.TypeOf(st2))
	assert.Equal(t, reflect.TypeOf(&MockStore{}), reflect.TypeOf(st3))
}

func TestLocalStore_RestoreLink(t *testing.T) {
	st, err := NewStore("", "local", 300)

	if err != nil {
		t.Fatalf("error occurred while creating local store instance")
	}

	st.Open()
	link := "https://abc.xyz"
	hash := algorithm.ComputeShortening(link)

	_, _ = st.CreateLink(link)
	result, _ := st.RestoreLink(hash)

	assert.Equal(t, link, result)

}

func TestLocalStore_CreateLink(t *testing.T) {
	st, err := NewStore("", "local", 300)

	if err != nil {
		t.Fatalf("error occurred while creating local store instance")
	}

	st.Open()
	link := "https://abc.xyz"
	hash := algorithm.ComputeShortening(link)

	result, _ := st.CreateLink(link)

	assert.Equal(t, hash, result)

}

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
