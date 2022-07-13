package urlapi

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/store"
)

type Server struct {
	config *Config
	router *mux.Router
	logger *logrus.Logger
	store  *store.Store
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (server *Server) Start(storeImpl string) error {
	if err := server.ConfigureLogger(); err != nil {
		return err
	}

	server.logger.Info("Starting server")
	server.ConfigureRouter()

	if err := server.ConfigureStore(storeImpl); err != nil {
		return err
	}

	return http.ListenAndServe(server.config.ServerAddress, server.router)
}

func (server *Server) ConfigureStore(storeImpl string) error {
	st, stErr := store.NewStore(server.config.Store, storeImpl)

	if stErr != nil {
		return stErr
	}

	if err := st.Open(); err != nil {
		return err
	}

	defer func(st store.Store) {
		err := st.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(st)

	server.store = &st
	server.logger.Info("Store created")

	return nil
}

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "It works")
		if err != nil {
			return
		}
	})
}

func (server *Server) ConfigureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)

	if err != nil {
		return err
	}

	server.logger.SetLevel(level)

	return nil
}
