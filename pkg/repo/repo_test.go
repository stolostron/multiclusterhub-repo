// Copyright (c) 2020 Red Hat, Inc.

package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/open-cluster-management/multiclusterhub-repo/pkg/config"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

func TestLiveness(t *testing.T) {

	handler := http.HandlerFunc(livenessHandler)
	req := httptest.NewRequest("GET", "/liveness", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestReadiness(t *testing.T) {
	handler := http.HandlerFunc(readinessHandler)
	req := httptest.NewRequest("GET", "/readiness", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFileServer(t *testing.T) {
	c := &config.Config{
		ChartDir:  "testdata/charts/",
		Namespace: "test",
		Port:      "8000",
		Service:   "test-service",
	}
	s, _ := New(c)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	tests := []struct {
		name string
		file string
		want int
	}{
		{"Get index", "index.yaml", http.StatusOK},
		{"Get nginx chart", "nginx-5.6.0.tgz", http.StatusOK},
		{"Get non-existant chart", "not-found.tgz", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(ts.URL + "/charts/" + tt.file)
			if err != nil {
				t.Errorf("request error: %v", err)
			}
			if status := res.StatusCode; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func Test_indexURL(t *testing.T) {
	type args struct {
		c *config.Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No namespace",
			args: args{
				c: &config.Config{
					Port:    "3000",
					Service: "multiclusterhub-repo",
				}},
			want: "http://multiclusterhub-repo:3000/charts",
		},
		{
			name: "Namespaced",
			args: args{
				c: &config.Config{
					Namespace: "test",
					Port:      "8000",
					Service:   "multiclusterhub-repo",
				}},
			want: "http://multiclusterhub-repo.test.svc.cluster.local:8000/charts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexURL(tt.args.c); got != tt.want {
				t.Errorf("indexURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createIndex(t *testing.T) {
	c := &config.Config{
		ChartDir:  "testdata/charts/",
		Namespace: "test",
		Port:      "3000",
		Service:   "multiclusterhub-repo",
	}
	b, err := createIndex(c)
	if err != nil {
		t.Errorf("createIndex() error = %v, wantErr %v", err, nil)
	}

	i := &repo.IndexFile{}
	if err := yaml.Unmarshal(b, i); err != nil {
		t.Errorf("createIndex() index failed to unmarshal: %s", err)
	}

	numEntries := len(i.Entries)
	if numEntries != 1 {
		t.Errorf("Expected 1 entry in index file but got %d", numEntries)
	}
	nginx, ok := i.Entries["nginx"]
	if !ok || len(nginx) != 1 {
		t.Errorf("Expected 1 nginx entry")
	}

}

func TestReindex(t *testing.T) {
	// Serve charts in tempDir
	tmpDir, err := ioutil.TempDir("testdata", "charts_tmp_")
	if err != nil {
		t.Fatalf("Could not create chart dir")
	}
	defer os.RemoveAll(tmpDir)

	c := &config.Config{
		// ChartDir:  "testdata/" + tmpDir + "/",
		ChartDir:  path.Join(tmpDir),
		Namespace: "test",
		Port:      "8000",
		Service:   "test-service",
	}

	s, _ := New(c)
	s.Start()
	defer s.Stop()

	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	// Empty chart directory
	res, err := http.Get(ts.URL + "/charts/" + "index.yaml")
	if err != nil {
		t.Errorf("request error: %v", err)
	}

	// Check that index has no entries
	i, err := getIndex(res.Body)
	if err != nil {
		t.Errorf("createIndex() index failed to get index: %s", err)
	}
	numEntries := len(i.Entries)
	if numEntries != 0 {
		t.Errorf("Expected 0 entry in index file but got %d", numEntries)
	}

	// Add chart (should trigger reindexing)
	err = copyFile("testdata/charts/nginx-5.6.0.tgz", path.Join(tmpDir, "nginx-5.6.0.tgz"))
	if err != nil {
		t.Errorf("Error adding chart: %v", err)
	}

	// Give time to reindex
	time.Sleep(4 * time.Second)

	// Populated chart directory
	res, err = http.Get(ts.URL + "/charts/" + "index.yaml")
	if err != nil {
		t.Errorf("request error: %v", err)
	}

	// Check that index now includes nginx
	i, err = getIndex(res.Body)
	if err != nil {
		log.Println(i)
		t.Errorf("createIndex() index failed to get index: %s", err)
	}

	nginx, ok := i.Entries["nginx"]
	if !ok || len(nginx) != 1 {
		t.Errorf("Expected 1 nginx entry")
	}
}

// helper for getting an index from an http response
func getIndex(res io.Reader) (*repo.IndexFile, error) {
	b, err := ioutil.ReadAll(res)
	i := &repo.IndexFile{}
	if err != nil {
		return i, err
	}
	if err := yaml.Unmarshal(b, i); err != nil {
		return i, err
	}
	return i, nil
}

// helper for copying charts
func copyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	return nil
}
