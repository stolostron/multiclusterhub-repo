package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// chartsDir is the directory that holds all charts to serve
	chartDir = "./charts"
	// gracePeriod is the duration for which the server will gracefully wait for existing connections to finish
	gracePeriod = time.Second * 15
)

func main() {
	envPort, exists := os.LookupEnv("RHACM_REPO_SERVICE_PORT")
	if !exists {
		log.Println("No port provided, defaulting to port :3000")
		envPort = "3000"
	}
	port := fmt.Sprintf(":%s", envPort)

	mux := http.NewServeMux()

	// Add route handlers
	fileServer := http.FileServer(http.Dir(chartDir))
	mux.Handle("/liveness", loggingMiddleware(http.HandlerFunc(livenessHandler)))
	mux.Handle("/readiness", loggingMiddleware(http.HandlerFunc(readinessHandler)))
	mux.Handle("/charts/", loggingMiddleware(http.StripPrefix("/charts", fileServer)))

	srv := &http.Server{
		Addr: port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      mux, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-sigs

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// loggingMiddleware logs each request sent to the server
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// livenessHandler returns a 200 status as long as the server is running
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readinessHandler returns a 200 status as long as the server is running
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
