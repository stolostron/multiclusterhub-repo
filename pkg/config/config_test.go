// Copyright (c) 2020 Red Hat, Inc.

package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("New config", func(t *testing.T) {
		os.Setenv("POD_NAMESPACE", "default")
		os.Setenv("MCH_CHART_DIR", "./charts/")
		os.Setenv("MCH_REPO_PORT", "3000")
		os.Setenv("MCH_REPO_SERVICE", "repo-svc")

		got := New()
		if got.Namespace != "default" ||
			got.ChartDir != "./charts/" ||
			got.Port != "3000" ||
			got.Service != "repo-svc" {
			t.Errorf("Config not set to env variables")
		}
	})

	t.Run("New config", func(t *testing.T) {
		os.Unsetenv("MCH_CHART_DIR")
		os.Unsetenv("POD_NAMESPACE")
		os.Unsetenv("MCH_REPO_PORT")
		os.Unsetenv("MCH_REPO_SERVICE")

		want := Config{
			ChartDir: defaultChartDir,
			Port:     defaultPort,
			Service:  defaultServiceName,
		}

		got := New()
		if got.ChartDir != defaultChartDir ||
			got.Namespace != "" ||
			got.Port != defaultPort ||
			got.Service != defaultServiceName {
			t.Errorf("New(), got = %v, want %v", got, want)
		}
	})

}

func Test_getEnv(t *testing.T) {
	os.Setenv("PRESENT", "present")

	type args struct {
		key        string
		defaultVal string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Existing env var",
			args{"PRESENT", "absent"},
			"present",
		},
		{
			"Missing env var",
			args{"ABSENT", "absent"},
			"absent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.key, tt.args.defaultVal); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
