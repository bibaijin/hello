package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap.NewProduction() failed, error: %v.", err)
	}
	defer logger.Sync()

	http.HandleFunc("/test", handle(test, logger))
	http.HandleFunc("/ping", handle(ping, logger))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("http.ListenAndServe() failed.",
			zap.String("address", ":8080"),
			zap.Error(err),
		)
	}
}

func test(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
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
