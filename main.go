package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap.NewProduction() failed, error: %v.", err)
	}
	defer logger.Sync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handle(ping, logger))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("server.ListenAndServe() failed.",
				zap.String("Addr", ":8080"),
				zap.Error(err),
			)
		}
	}()
	logger.Info("server.ListenAndServe()...", zap.String("Addr", ":8080"))
	defer func() {
		if err := server.Close(); err != nil {
			logger.Error("server.Close() failed.", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Shutting down...")
	time.Sleep(1 * time.Second)
}

func ping(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
	if _, err := w.Write([]byte("OK")); err != nil {
		logger.Error("w.Write([]byte(\"OK\")) failed.", zap.Error(err))
	}
}

func handle(f func(w http.ResponseWriter, r *http.Request, logger *zap.Logger), logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newLogger := logger.With(zap.String("RequestID", newRequestID()))
		newLogger.Info("Receive a new request.",
			zap.String("RemoteAddr", r.RemoteAddr),
			zap.String("Method", r.Method),
			zap.String("URL", r.URL.String()),
			zap.Any("Header", r.Header),
		)
		f(w, r, newLogger)
		newLogger.Info("Response is sent.",
			zap.String("RemoteAddr", r.RemoteAddr),
			zap.String("Method", r.Method),
			zap.String("URL", r.URL.String()),
			zap.Any("Header", r.Header),
		)
	}
}

func newRequestID() string {
	bs := make([]byte, 16)
	if _, err := rand.Read(bs); err != nil {
		return "0"
	}

	return hex.EncodeToString(bs)
}
