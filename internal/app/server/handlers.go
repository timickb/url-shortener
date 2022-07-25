package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

func (server *APIServer) error(wr http.ResponseWriter, status int, err error) {
	server.makeResponse(wr, status, map[string]string{"error": err.Error()})
}

func (server *APIServer) makeResponse(wr http.ResponseWriter, status int, data interface{}) {
	wr.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(wr).Encode(data)
	}
}

func (server *APIServer) handleCreate() http.HandlerFunc {
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
			server.error(writer, http.StatusBadRequest, err)
			return
		}

		if req.Url == "" {
			server.error(writer, http.StatusBadRequest, errors.New("cannot parse empty string"))
			return
		}

		go func(wg *sync.WaitGroup) {
			result, err := server.store.CreateLink(req.Url)
			if err != nil {
				server.error(writer, http.StatusBadRequest, err)
				return
			}
			server.makeResponse(writer, http.StatusOK, &createResponse{Hash: result})
			wg.Done()
		}(&wg)
		wg.Wait()
	}
}

func (server *APIServer) handleRestore() http.HandlerFunc {
	type restoreResponse struct {
		Original string `json:"original"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		hash := request.FormValue("hash")

		result, err := server.store.RestoreLink(hash)

		if err != nil {
			server.error(writer, http.StatusBadRequest, err)
			return
		}

		server.makeResponse(writer, http.StatusOK, &restoreResponse{Original: result})
	}
}
