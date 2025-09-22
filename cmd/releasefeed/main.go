package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"poschi3/releasefeed/internal/endoflife"
	"poschi3/releasefeed/internal/feed"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	configureLogger()
	slog.Info("Releasefeed started")

	s := CreateNewServer()
	s.MountHandlers()

	http.ListenAndServe(s.Addr, s.Router)
}

type Server struct {
	Router *chi.Mux
	Addr   string
}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()

	s.Addr = os.Getenv("RELEASEFEED_LISTEN_ADDR")
	if s.Addr == "" {
		s.Addr = "127.0.0.1:8090"
	}
	return s
}

func (s *Server) MountHandlers() {
	// Mount all Middleware here
	if os.Getenv("RELEASEFEED_MIDDLEWARE_REALIP") == "true" {
		s.Router.Use(middleware.RealIP)
	}
	middlewareLogger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: middlewareLogger})
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Compress(5, "application/atom+xml"))

	// Mount all handlers here
	s.Router.Get("/{product}", handleProduct)
	s.Router.Get("/{product}/{cycle}", handleCycle)
}

func configureLogger() {
	logLevelVar := os.Getenv("RELEASEFEED_LOG_LEVEL")
	if logLevelVar == "" {
		slog.SetLogLoggerLevel(slog.LevelInfo)
		return
	}

	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(logLevelVar))
	if err != nil {
		slog.SetLogLoggerLevel(slog.LevelInfo)
		slog.Error("Invalid log level", "RELEASEFEED_LOG_LEVEL", logLevelVar, "error", err)
		return
	}
	slog.SetLogLoggerLevel(logLevel)
}

func handleProduct(w http.ResponseWriter, req *http.Request) {
	productName := chi.URLParam(req, "product")
	product, err := endoflife.GetProduct(productName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	feed, err := feed.FeedProduct(strings.Split(req.Host, ":")[0], productName, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeFeed(w, feed)
}

func handleCycle(w http.ResponseWriter, req *http.Request) {
	productName := chi.URLParam(req, "product")
	cycleName := chi.URLParam(req, "cycle")
	cycle, err := endoflife.GetCycle(productName, cycleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	feed, err := feed.FeedCycle(strings.Split(req.Host, ":")[0], productName, cycle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeFeed(w, feed)
}

func writeFeed(w http.ResponseWriter, feed string) {
	w.Header().Add("content-type", "application/atom+xml; charset=UTF-8")
	fmt.Fprintln(w, feed)
}
