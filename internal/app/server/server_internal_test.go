package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerRestoreEndpoint(t *testing.T) {
	t.Helper()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/restore?hash=t3st_h4sh_", nil)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s, err := NewServer(&Config{ServerPort: 8080, StoreImpl: "test"})
	if err != nil {
		t.Fatalf("an error occured while creating server instance")
	}
	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestServerCreateEndpoint(t *testing.T) {
	t.Helper()

	payloadString := `{"url": "https://example.com"}`

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/create", strings.NewReader(payloadString))

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s, err := NewServer(&Config{ServerPort: 8080, StoreImpl: "test"})
	if err != nil {
		t.Fatalf("an error occured while creating server instance")
	}
	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}
