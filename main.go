package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var index []byte

const (
	// chartsDir is the directory that holds all charts to serve
	chartDir = "./multicloudhub/charts/"
	// gracePeriod is the duration for which the server will gracefully wait for existing connections to finish
	gracePeriod = time.Second * 15
	// port the server listens on
	port = ":3000"
)

func main() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	// Hold index file in memory
	index, err := readIndex()
	if err != nil {
		log.Fatal(err)
	}

	// Modify urls in index to reference namespace deployed in
	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		log.Println("POD_NAMESPACE missing, index.yaml will not be updated")
	} else {
		index = modifyIndex(index, ns)
	}

	mux := http.NewServeMux()

	// Add route handlers
	fileServer := http.FileServer(http.Dir(chartDir))
	mux.Handle("/liveness", loggingMiddleware(http.HandlerFunc(livenessHandler)))
	mux.Handle("/readiness", loggingMiddleware(http.HandlerFunc(readinessHandler)))
	mux.Handle("/charts/index.yaml", loggingMiddleware(http.HandlerFunc(indexHandler)))
	mux.Handle("/charts/", loggingMiddleware(http.StripPrefix("/charts/", fileServer)))

	srv := &http.Server{
		Addr: port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      mux,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Serving on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	// Kubernetes sends a SIGTERM, waits for a grace period, and then a SIGKILL
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-sigs
	log.Println("Shutdown signal received")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("Goodbye")
	os.Exit(0)
}

// CustomResponseWriter adds a field an http.ResponseWriter to track status
type CustomResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader populates the status field before calling WriteHeader
func (w *CustomResponseWriter) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	w.ResponseWriter.WriteHeader(status)
}

// loggingMiddleware logs each request sent to the server
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use custom ResponseWriter to track statuscode
		crw := &CustomResponseWriter{ResponseWriter: w}

		startTime := time.Now()
		next.ServeHTTP(crw, r)
		duration := time.Since(startTime)

		log.Printf("%d %3dms %s", crw.status, duration.Milliseconds(), r.RequestURI)
	})
}

func readIndex() ([]byte, error) {
	f, err := ioutil.ReadFile(chartDir + "index.yaml")
	if err != nil {
		return nil, err
	}
	log.Println(len(f))
	return f, nil
}

func modifyIndex(index []byte, ns string) []byte {
	oldURL := []byte("multicloudhub-repo:3000")
	newURL := []byte(fmt.Sprintf("multicloudhub-repo.%s.svc.cluster.local:3000", ns))

	newIndex := bytes.ReplaceAll(index, oldURL, newURL)
	return newIndex
}

// indexHandler serves the index.yaml file from in memory
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Hold index file in memory
	index, err := readIndex()
	if err != nil {
		log.Fatal(err)
	}

	// Modify urls in index to reference namespace deployed in
	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		log.Println("POD_NAMESPACE missing, index.yaml will not be updated")
	} else {
		index = modifyIndex(index, ns)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(index)
}

// livenessHandler returns a 200 status as long as the server is running
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readinessHandler returns a 200 status as long as the server is running
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
