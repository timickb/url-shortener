package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_HandleURLCreate(t *testing.T) {
	payloadString := "{\"url\": \"https://example.com\"}"

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/create", strings.NewReader(payloadString))

	s := NewServer(&Config{ServerAddress: ":8080", StoreImpl: "test"})
	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}
