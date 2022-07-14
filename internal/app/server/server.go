package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/store"
)

type Server struct {
	config *Config
	router *mux.Router
	store  store.Store
}

func NewServer(config *Config) *Server {
	serv := &Server{
		config: config,
		router: mux.NewRouter(),
	}

	serv.ConfigureRouter()
	return serv
}

func (server *Server) Start() error {
	defer func(store store.Store) {
		err := store.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server.store)

	logrus.Info("Starting server")

	if err := server.ConfigureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(server.config.ServerAddress, server.router)
}

func (server *Server) ConfigureStore() error {
	st, stErr := store.NewStore(
		server.config.ConnectionString,
		server.config.StoreImpl,
		server.config.MaxUrlLength)

	if stErr != nil {
		return stErr
	}

	if err := st.Open(); err != nil {
		return err
	}
	server.store = st
	logrus.Info("Store configured")

	return nil
}

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/create", server.handleCreate()).Methods("POST")
	server.router.HandleFunc("/restore", server.handleRestore()).Methods("GET")
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.router.ServeHTTP(writer, request)
}
