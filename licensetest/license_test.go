// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project


// Package license scans the repo for missing license or copyright headers
package licensetest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// slashScanner is for validating the copyright comment in Go files
var slashScanner = regexp.MustCompile(`// Copyright \(c\) 2020 Red Hat, Inc\.`)

// poundScanner is for validating the copyright comment in shell and Python files
var poundScanner = regexp.MustCompile(`\# Copyright \(c\) 2020 Red Hat, Inc\.`)

var skip = map[string]bool{
	// Build Harness
	"../vbh":                      true,
	"../build-harness-extensions": true,
	"../build-harness":            true,
}

func TestLicense(t *testing.T) {
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if skip[path] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if err != nil {
			return err
		}

		// Capture Go code, Python code, and shell scripts
		if filepath.Ext(path) != ".go" && filepath.Ext(path) != ".sh" && filepath.Ext(path) != ".py" {
			return nil
		}

		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		// Find license
		if filepath.Ext(path) == ".go" {
			if !slashScanner.Match(src) {
				t.Errorf("%v: license header not present", path)
				return nil
			}
		} else {
			if !poundScanner.Match(src) {
				t.Errorf("%v: license header not present", path)
				return nil
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
