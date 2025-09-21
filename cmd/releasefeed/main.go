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

	configureMiddlewareLogger()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if os.Getenv("RELEASEFEED_MIDDLEWARE_REALIP") == "true" {
		r.Use(middleware.RealIP)
	}

	r.Get("/{product}", handleProduct)
	r.Get("/{product}/{cycle}", handleCycle)
	addr := os.Getenv("RELEASEFEED_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8090"
	}
	http.ListenAndServe(addr, r)
}

func configureMiddlewareLogger() {
	middlewareLogger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: middlewareLogger})
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
