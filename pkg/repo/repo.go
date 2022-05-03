// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package repo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/stolostron/multiclusterhub-repo/pkg/config"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

// Create request router with modified index
func (s *Server) SetupRouter() error {
	mux := http.NewServeMux()

	// Add route handlers
	fileServer := http.FileServer(http.Dir(s.Config.ChartDir))
	mux.Handle("/liveness", http.HandlerFunc(livenessHandler))
	mux.Handle("/readiness", http.HandlerFunc(readinessHandler))
	mux.Handle("/charts/index.yaml", loggingMiddleware(http.HandlerFunc(s.indexHandler)))
	mux.Handle("/charts/", loggingMiddleware(http.StripPrefix("/charts/", fileServer)))

	s.Router = mux
	return nil
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

// readIndex builds an index from a flat directory
func createIndex(c *config.Config) ([]byte, error) {
	url := indexURL(c)
	index, err := repo.IndexDirectory(filepath.Clean(c.ChartDir), url)
	if err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(index)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// packageCharts
func packageCharts(c *config.Config) ([]string, error) {
	if _, err := semver.NewVersion(c.Version); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(c.ChartDir)
	if err != nil {
		log.Fatal(err)
	}

	charts := []string{}
	for _, f := range files {
		if f.IsDir() {
			pkgName, err := packageChart(path.Join(c.ChartDir, f.Name()), c.ChartDir, c.Version)
			if err != nil {
				return charts, errors.Wrap(err, fmt.Sprintf("failed to package directory %s", f.Name()))
			}
			charts = append(charts, pkgName)
		}
	}
	return charts, nil
}

func packageChart(src string, dst string, chartVersion string) (string, error) {
	ch, err := loader.LoadDir(src)
	ch.Metadata.Version = chartVersion
	name, err := chartutil.Save(ch, dst)
	if err != nil {
		return "", errors.Wrap(err, "failed to save")
	}
	log.Printf("Packaged chart as %s", name)
	return name, nil
}

func cleanupCharts(c *config.Config, charts []string) error {
	errorMessages := []string{}
	for _, file := range charts {
		err := os.Remove(file)
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	if len(errorMessages) > 0 {
		log.Println(errors.New(fmt.Sprintf(": %s", strings.Join(errorMessages, "; "))))
		return errors.New(fmt.Sprintf(": %s", strings.Join(errorMessages, "; ")))
	}
	return nil
}

// readIndex builds an index from a flat directory
func (s *Server) Reindex() error {
	log.Println("Reindexing")
	s.Lock()
	defer s.Unlock()

	index, err := createIndex(s.Config)
	if err != nil {
		return err
	}

	s.Index = index
	return nil
}

// indexURL returns a formatted URL based on Config parameters
func indexURL(c *config.Config) string {
	if c.Namespace == "" {
		return fmt.Sprintf("http://%s:%s/charts", c.Service, c.Port)
	}
	return fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/charts", c.Service, c.Namespace, c.Port)
}

// indexHandler serves the index.yaml file from in memory
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()
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
