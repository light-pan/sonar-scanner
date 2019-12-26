package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/light-pan/sonar-scanner/handle"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	var wait time.Duration
	var jsonDir string
	var jarDir string
	var command string
	var logLevel string
	flag.DurationVar(&wait, "graceful-timeout", time.Second*3, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&jsonDir, "json-dir", "/data/sonar-scanner/json", "the scanner json files path")
	flag.StringVar(&jarDir, "jar-dir", "/data/sonar-scanner/jar", "the scanner jar files path")
	flag.StringVar(&command, "command", "/data/sonar-scanner/bin/sonar-scanner", "the scanner app path")
	flag.StringVar(&logLevel, "log-level", "info", "the scanner log level")
	flag.Parse()

	logger := logrus.New()

	logger.SetOutput(&lumberjack.Logger{
		Filename:   "scanner.log",
		MaxSize:    100,
		MaxBackups: 3,
		LocalTime:  true,
	})

	r := mux.NewRouter()

	h := &handle.Handle{
		Logger:   logger,
		JSONDir:  jsonDir,
		JarDir:   jarDir,
		LogLevel: logLevel,
		Command:  command,
	}

	r.HandleFunc("/admin/config", h.SetConfig).Methods(http.MethodPut)
	r.HandleFunc("/admin/config", h.GetConfig).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id:[0-9]+}", h.Scanner).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id:[0-9]+}", h.DeleteTask).Methods(http.MethodDelete)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 3,
		IdleTimeout:  time.Second * 30,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorln(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Infoln("shutting down")
	os.Exit(0)
}
