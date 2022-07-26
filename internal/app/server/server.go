package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/store"
)

type APIServer struct {
	config           *Config
	router           *mux.Router
	store            store.Store
	logger           *logrus.Logger
	connectionString string
}

func New(config *Config, logger *logrus.Logger) (*APIServer, error) {
	srv := &APIServer{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}

	srv.connectionString = fmt.Sprintf("host=%s port=%d user=%s password=\"%s\" dbname=%s sslmode=disable",
		srv.config.Database.DbHost,
		srv.config.Database.DbPort,
		srv.config.Database.DbUser,
		srv.config.Database.DbPassword,
		srv.config.Database.DbName)

	if err := srv.configureStore(); err != nil {
		return nil, err
	}

	srv.configureRouter()

	return srv, nil
}

func (s *APIServer) Start() error {

	s.logger.Info(fmt.Sprintf("Starting server on port %d", s.config.ServerPort))

	defer func(store store.Store) {
		_ = store.Close()
	}(s.store)

	return http.ListenAndServe(":"+strconv.Itoa(s.config.ServerPort), s.router)
}

func (s *APIServer) Close() error {
	return s.store.Close()
}

func (s *APIServer) configureStore() error {
	s.logger.Info(fmt.Sprintf("Configuring store %s", s.config.StoreImpl))

	db, err := sql.Open("postgres", s.connectionString)

	if err != nil {
		s.logger.Warnf("Couldn't open db connection: %s", err.Error())
	} else {
		s.logger.Info("DB connection set")
	}

	st, err := store.New(db, s.logger, s.config.StoreImpl)

	if err != nil {
		return err
	}

	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	s.logger.Info("Store configured")

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/create", s.handleCreate()).Methods("POST")
	s.router.HandleFunc("/restore", s.handleRestore()).Methods("GET")
}

func (s *APIServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}
