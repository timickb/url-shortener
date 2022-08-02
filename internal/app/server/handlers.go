package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
)

func (s *APIServer) validateURL(url string) bool {
	// no spaces allowed
	parts := strings.Split(url, " ")
	if len(parts) > 1 {
		return false
	}

	// at least one dot must be
	parts = strings.Split(url, ".")
	if len(parts) < 2 {
		return false
	}

	if s.config.MaxUrlLength <= 0 {
		return url != "" && len(url) <= 500
	}
	return url != "" && len(url) <= s.config.MaxUrlLength
}

func (s *APIServer) error(wr http.ResponseWriter, status int, err error) {
	s.makeResponse(wr, status, map[string]string{"error": err.Error()})
}

func (s *APIServer) makeResponse(wr http.ResponseWriter, status int, data interface{}) {
	wr.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(wr).Encode(data)
	}
}

func (s *APIServer) handleCreate() http.HandlerFunc {
	type createRequest struct {
		Url string `json:"url"`
	}
	type createResponse struct {
		Hash string `json:"hash"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)

		writer.Header().Set("Content-Type", "application/json")

		req := &createRequest{}
		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			s.error(writer, http.StatusBadRequest, err)
			return
		}

		if !s.validateURL(req.Url) {
			s.error(writer, http.StatusBadRequest, errors.New("cannot parse empty string"))
			return
		}

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			result, err := s.store.CreateLink(strings.Trim(req.Url, " "))

			if err != nil {
				s.error(writer, http.StatusBadRequest, err)
				return
			}

			s.makeResponse(writer, http.StatusOK, &createResponse{Hash: result})
		}(&wg)

		wg.Wait()
	}
}

func (s *APIServer) handleRestore() http.HandlerFunc {
	type restoreResponse struct {
		Original string `json:"original"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		hash := request.FormValue("hash")

		result, err := s.store.RestoreLink(hash)

		if err != nil {
			s.error(writer, http.StatusBadRequest, err)
			return
		}

		s.makeResponse(writer, http.StatusOK, &restoreResponse{Original: result})
	}
}
