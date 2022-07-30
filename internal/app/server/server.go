package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/store"
)

type APIServer struct {
	config *Config
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func WithLogger(logger *logrus.Logger) func(*APIServer) {
	return func(s *APIServer) {
		s.logger = logger
	}
}

func WithConfig(config *Config) func(*APIServer) {
	return func(s *APIServer) {
		s.config = config
	}
}

func WithStore(store store.Store) func(*APIServer) {
	return func(s *APIServer) {
		s.store = store
	}
}

func New(options ...func(s *APIServer)) *APIServer {
	srv := &APIServer{}

	for _, applyOpt := range options {
		applyOpt(srv)
	}

	if srv.logger == nil {
		fmt.Println("Logger not specidied. Using default")
		srv.logger = logrus.New()
	}

	if srv.store == nil {
		fmt.Println("Store not specified. Using default")
		srv.store = store.New()
		srv.store.Open()
	}

	if srv.config == nil {
		fmt.Println("Config not specified. Using default")
		srv.config = DefaultConfig()
	}

	srv.router = mux.NewRouter()
	srv.router.HandleFunc("/create", srv.handleCreate()).Methods("POST")
	srv.router.HandleFunc("/restore", srv.handleRestore()).Methods("GET")

	return srv
}

func (s *APIServer) Start() error {
	s.logger.Info(fmt.Sprintf("Starting server on port %d", s.config.ServerPort))
	return http.ListenAndServe(":"+strconv.Itoa(s.config.ServerPort), s.router)
}

func (s *APIServer) Close() error {
	return s.store.Close()
}

func (s *APIServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}
