// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/stolostron/multiclusterhub-repo/pkg/config"
	"github.com/stolostron/multiclusterhub-repo/pkg/repo"
)

func main() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	c := config.New()
	server, err := repo.New(c)
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr: ":" + c.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      time.Second * 30,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Second * 30,
		Handler:           server.Router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Serving on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	// Kubernetes sends a SIGTERM, waits for a grace period, and then a SIGKILL
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	sig := <-sigs
	log.Printf("Received signal: %s", sig.String())

	// Stop file watcher
	err = server.Stop()
	if err != nil {
		panic(err)
	}

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Goodbye")
	os.Exit(0)
}
