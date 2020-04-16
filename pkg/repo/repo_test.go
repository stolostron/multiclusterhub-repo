package repo

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/open-cluster-management/multicloudhub-repo/pkg/config"
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

func Test_modifyIndex(t *testing.T) {
	type args struct {
		index   []byte
		ns      string
		service string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"Update index",
			args{
				index:   []byte("http://multiclusterhub-repo:3000/charts/application-chart-1.0.0.tgz"),
				ns:      "default",
				service: "multiclusterhub-repo",
			},
			[]byte("http://multiclusterhub-repo.default:3000/charts/application-chart-1.0.0.tgz"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifyIndex(tt.args.index, tt.args.ns, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modifyIndex() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestFileServer(t *testing.T) {
	c := &config.Config{
		ChartDir:  "../../multiclusterhub/charts/",
		Namespace: "test",
		Port:      "8000",
		Service:   "test-service",
	}
	ts := httptest.NewServer(SetupRouter(c))
	defer ts.Close()

	tests := []struct {
		name string
		file string
		want int
	}{
		{"Get index", "index.yaml", http.StatusOK},
		{"Get cert-manager chart", "cert-manager-3.5.0.tgz", http.StatusOK},
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
