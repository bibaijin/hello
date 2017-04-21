package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// LogFlag 控制日志的前缀
const LogFlag = log.LstdFlags | log.Lmicroseconds | log.Lshortfile

var (
	errLogger  = log.New(os.Stderr, "ERROR ", LogFlag)
	infoLogger = log.New(os.Stdout, "INFO ", LogFlag)
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		infoLogger.Printf("Receive a ping request.")

		if _, err := fmt.Fprintf(w, "OK"); err != nil {
			errLogger.Printf("fmt.Fprintf failed, error: %s.", err)
		}

		infoLogger.Printf("Response is sent.")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("http.ListenAndServe failed, error: %s.", err)
	}
}
