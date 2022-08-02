package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timickb/url-shortener/internal/app/store"
)

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig()

	assert.Equal(t, conf.ServerPort, 8080)
	assert.Equal(t, conf.MaxUrlLength, 300)
}

func TestServerGeneralResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	http.NewRequest(http.MethodGet, "/some-strange-endpoint", nil)

	s := New()

	s.makeResponse(rec, 201, nil)
	assert.Equal(t, rec.Code, http.StatusCreated)
}

func TestServerErrorResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	http.NewRequest(http.MethodGet, "/some-strange-endpoint", nil)

	s := New()

	s.error(rec, 400, errors.New("bad request"))
	assert.Equal(t, rec.Code, http.StatusBadRequest)

}

func TestServerRestoreEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/restore?hash=t3st_h4sh_", nil)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s := New(WithStore(&store.StubStore{}))
	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestServerCreateEndpoint(t *testing.T) {
	payloadString := `{"url": "https://example.com"}`

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/create", strings.NewReader(payloadString))

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s := New(WithStore(&store.StubStore{}))
	s.ServeHTTP(rec, req)

	fmt.Println(rec.Body)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServerCreateEndpointInvalidJSON(t *testing.T) {
	payloadString := `invalid json`

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/create", strings.NewReader(payloadString))

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s := New()
	s.ServeHTTP(rec, req)

	assert.NotEqual(t, rec.Code, http.StatusOK)
}

func TestServerCreateEmptyURL(t *testing.T) {
	payloadString := `{"url": ""}`

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/create", strings.NewReader(payloadString))

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	s := New()
	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusBadRequest)
}
