package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/store"
)

type Server struct {
	config           *Config
	router           *mux.Router
	store            store.Store
	connectionString string
}

func NewServer(config *Config) *Server {
	srv := &Server{
		config: config,
		router: mux.NewRouter(),
	}

	srv.connectionString = fmt.Sprintf("host=%s port=%d user=%s password=\"%s\" dbname=%s sslmode=disable",
		srv.config.Database.DbHost,
		srv.config.Database.DbPort,
		srv.config.Database.DbUser,
		srv.config.Database.DbPassword,
		srv.config.Database.DbName)

	srv.ConfigureRouter()

	return srv
}

func (server *Server) Start() error {

	logrus.Info(fmt.Sprintf("Starting server on %s", server.config.ServerAddress))

	if err := server.ConfigureStore(); err != nil {
		return err
	}

	defer func(store store.Store) {
		_ = store.Close()
	}(server.store)

	return http.ListenAndServe(server.config.ServerAddress, server.router)
}

func (server *Server) ConfigureStore() error {
	logrus.Info(fmt.Sprintf("Configuring store %s", server.config.StoreImpl))
	st, stErr := store.NewStore(
		server.connectionString,
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
