// Copyright (c) 2020 Red Hat, Inc.

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

	// Hold index file in memory
	index, err := createIndex(c)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Index:  index,
		Config: c,
	}

	s.SetupRouter()

	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return s, nil
}
