// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package repo

import (
	"fmt"
	"log"
)

// Start starts the watch on the chart directory files.
func (s *Server) Start() error {
	if err := s.watcher.Add(s.Config.RepoDir); err != nil {
		return err
	}

	go s.Watch()

	fmt.Println("Starting file watcher")
	return nil
}

// Stop closes the server's file watcher.
func (s *Server) Stop() error {
	fmt.Println("Stopping file watcher")
	return s.watcher.Close()
}

// Watch reads events from the watcher's channel and reacts to changes.
func (s *Server) Watch() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			// Channel is closed.
			if !ok {
				return
			}

			log.Println("File change detected: ", event.String())
			if err := s.Reindex(); err != nil {
				log.Println(err)
			}

		case err, ok := <-s.watcher.Errors:
			// Channel is closed.
			if !ok {
				return
			}

			log.Println(err)
		}
	}
}
