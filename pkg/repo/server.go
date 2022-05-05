// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package repo

import (
	"net/http"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/stolostron/multiclusterhub-repo/pkg/config"
)

// Server holds an index.yaml
type Server struct {
	sync.Mutex
	Index   []byte
	Config  *config.Config
	Router  *http.ServeMux
	watcher *fsnotify.Watcher
}

// New returns a new Server
func New(c *config.Config) (*Server, error) {
	var err error

	// Generate packages from charts
	_, err = packageCharts(c)
	if err != nil {
		return nil, err
	}

	// Hold index file in memory
	index, err := createIndex(c)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Index:  index,
		Config: c,
	}

	err = s.SetupRouter()
	if err != nil {
		return nil, err
	}

	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return s, nil
}
