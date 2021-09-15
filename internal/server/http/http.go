package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gpiecyk/data-warehouse/internal/api"
)

const gracefulTimeout time.Duration = time.Second * 15

type Handlers struct {
	api *api.API
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	healthData, err := h.api.Health()
	if err != nil {
		http.Error(w, "Service unavailable: cannot check current health of the application", http.StatusServiceUnavailable)
		return
	}

	data, err := json.Marshal(healthData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handlers) routes(graphQL *handler.Server) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/-/health", h.Health)
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", h.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id:[0-9]+}", h.GetUserById).Methods("GET")
	router.HandleFunc("/users", h.FindUsers).Methods("GET").Queries("limit", "{limit:[0-9]+}")

	// graphQL
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", graphQL)
	return router
}

type HTTP struct {
	server *http.Server
	config *Config
}

func (http *HTTP) Start() {
	log.Printf("Listening on port %s", http.config.Port)
	go func() {
		if err := http.server.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe %s: %v", http.config.Host, err)
		}
	}()
	http.gracefullyShutdownServerWithTimeout()
}

func (http *HTTP) gracefullyShutdownServerWithTimeout() {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt)

	// block until the shutdown signal is received
	<-shutdownChannel

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel() // releases resources if Shutdown completes before timeout elapses

	http.server.Shutdown(ctx)
	os.Exit(0)
}

type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewService(config *Config, api *api.API, graphQL *handler.Server) (*HTTP, error) {
	handler := Handlers{
		api: api,
	}
	router := handler.routes(graphQL)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler:           router,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.ReadTimeout * 2,
	}

	return &HTTP{
		server: httpServer,
		config: config,
	}, nil
}
