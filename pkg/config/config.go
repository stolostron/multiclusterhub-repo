package config

import (
	"log"
	"os"
)

var (
	defaultChartDir    = "./multiclusterhub/charts/"
	defaultPort        = "3000"
	defaultServiceName = "multiclusterhub-repo"
)

type Config struct {
	// Directory to serve files from
	ChartDir string
	// Namespace charts are served from
	Namespace string
	// Port to serve on
	Port string
	// Kubernetes service exposing app
	Service string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		ChartDir:  getEnv("MCH_CHART_DIR", defaultChartDir),
		Namespace: getEnv("POD_NAMESPACE", ""),
		Port:      getEnv("MCH_REPO_PORT", defaultPort),
		Service:   getEnv("MCH_REPO_SERVICE", defaultServiceName),
	}
}

// Helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("%s not set, defaulting to '%s'", key, defaultVal)
	return defaultVal
}
