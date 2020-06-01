// Copyright (c) 2020 Red Hat, Inc.

package repo

import (
	"fmt"
	"log"
)

// Start starts the watch on the chart directory files.
func (s *Server) Start() error {
	if err := s.watcher.Add(s.Config.ChartDir); err != nil {
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

			fmt.Println("File change detected: ", event.String())
			s.Reindex()

		case err, ok := <-s.watcher.Errors:
			// Channel is closed.
			if !ok {
				return
			}

			log.Println(err)
		}
	}
}
