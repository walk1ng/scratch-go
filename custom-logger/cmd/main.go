package main

import (
	"custom-logger-demo/jsonlog"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8081),
		Handler:      http.DefaultServeMux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(logger, "http_srv", 0),
	}

	logger.PrintInfo("starting server", map[string]interface{}{
		"addr": srv.Addr,
		"pid":  os.Getpid(),
	})

	if err := srv.ListenAndServe(); err != nil {
		logger.PrintFatal(err, nil)
	}

}
