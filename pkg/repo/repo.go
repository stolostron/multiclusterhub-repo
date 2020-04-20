package repo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/open-cluster-management/multicloudhub-repo/pkg/config"
)

// Server holds an index.yaml
type Server struct {
	Index []byte
}

// Create request router with modified index
func SetupRouter(c *config.Config) *http.ServeMux {
	// Hold index file in memory
	index, err := readIndex(c.ChartDir)
	if err != nil {
		log.Fatal(err)
	}

	// Modify urls in index to reference namespace deployed in
	if ns := c.Namespace; ns != "" {
		log.Printf("Updating index with namespace '%s'", ns)
		index = modifyIndex(index, ns, c.Service)
	}

	s := &Server{Index: index}
	mux := http.NewServeMux()

	// Add route handlers
	fileServer := http.FileServer(http.Dir(c.ChartDir))
	mux.Handle("/liveness", http.HandlerFunc(livenessHandler))
	mux.Handle("/readiness", http.HandlerFunc(readinessHandler))
	mux.Handle("/charts/index.yaml", loggingMiddleware(http.HandlerFunc(s.indexHandler)))
	mux.Handle("/charts/", loggingMiddleware(http.StripPrefix("/charts/", fileServer)))

	return mux
}

// StatusWriter adds a field an http.ResponseWriter to track status
type StatusWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader populates the status field before calling WriteHeader
func (w *StatusWriter) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	w.ResponseWriter.WriteHeader(status)
}

// loggingMiddleware logs each request sent to the server
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use custom ResponseWriter to track statuscode
		crw := &StatusWriter{ResponseWriter: w}

		startTime := time.Now()
		next.ServeHTTP(crw, r)
		duration := time.Since(startTime)

		log.Printf("%d %3dms %s", crw.status, duration.Milliseconds(), r.RequestURI)
	})
}

func readIndex(dir string) ([]byte, error) {
	filePath := filepath.Join(filepath.Clean(dir), "index.yaml")
	f, err := ioutil.ReadFile(filePath) // #nosec G304 (index path not configurable by user)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// modifyIndex makes urls namespace-specific
func modifyIndex(index []byte, ns string, service string) []byte {
	oldURL := []byte(service)
	newURL := []byte(fmt.Sprintf("%s.%s", service, ns))

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
