package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type config struct {
	WorkerID      string `json:"worker_id"`
	SiteID        string `json:"site_id"`
	SiteLocation  string `json:"site_location"`
	SiteIsp       string `json:"site_isp"`
	Elasticsearch string `json:"elasticsearch"`
}

var con = &config{}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*3, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/admin/config", setConfig).Methods(http.MethodPut)
	r.HandleFunc("/admin/config", getConfig).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id:[0-9]+}", scanner).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id:[0-9]+}", deleteTask).Methods(http.MethodDelete)
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 3,
		IdleTimeout:  time.Second * 30,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
