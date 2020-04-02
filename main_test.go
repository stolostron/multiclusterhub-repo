package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
	ns := "default"

	type args struct {
		index []byte
		ns    string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"Update index",
			args{index: []byte("http://" + serviceName + "/index.yaml"), ns: ns},
			[]byte("http://" + serviceName + "." + ns + "/index.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifyIndex(tt.args.index, tt.args.ns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("modifyIndex() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestFileServer(t *testing.T) {
	ts := httptest.NewServer(setupRouter())
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
