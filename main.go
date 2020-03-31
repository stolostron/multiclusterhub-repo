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

var (
	chartDir    = "./multiclusterhub/charts/"
	ns          = os.Getenv("POD_NAMESPACE")
	port        = "3000"
	serviceName = "multiclusterhub-repo"
)

// Server holds an index.yaml
type Server struct {
	Index []byte
}

// Create request router with modified index
func setupRouter() *http.ServeMux {
	// Hold index file in memory
	index, err := readIndex()
	if err != nil {
		log.Fatal(err)
	}

	// Modify urls in index to reference namespace deployed in
	if ns != "" {
		log.Printf("Updating index with namespace '%s'", ns)
		index = modifyIndex(index, ns)
	}

	s := &Server{Index: index}
	mux := http.NewServeMux()

	// Add route handlers
	fileServer := http.FileServer(http.Dir(chartDir))
	mux.Handle("/liveness", http.HandlerFunc(livenessHandler))
	mux.Handle("/readiness", http.HandlerFunc(readinessHandler))
	mux.Handle("/charts/index.yaml", http.HandlerFunc(s.indexHandler))
	mux.Handle("/charts/", loggingMiddleware(http.StripPrefix("/charts/", fileServer)))

	return mux
}

func main() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	mux := setupRouter()
	srv := &http.Server{
		Addr: ":" + port,
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
	sig := <-sigs
	log.Printf("Received signal: %s", sig.String())

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

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
	return f, nil
}

// modifyIndex makes urls namespace-specific
func modifyIndex(index []byte, ns string) []byte {
	oldURL := []byte(serviceName)
	newURL := []byte(fmt.Sprintf("%s.%s", serviceName, ns))

	newIndex := bytes.ReplaceAll(index, oldURL, newURL)
	return newIndex
}

// indexHandler serves the index.yaml file from in memory
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(s.Index); err != nil {
		log.Println(err)
	}
}

// livenessHandler returns a 200 status as long as the server is running
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readinessHandler returns a 200 status as long as the server is running
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
