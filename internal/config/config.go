// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"math"

	"github.com/fsouza/fake-gcs-server/fakestorage"
)

const (
	filesystemBackend = "filesystem"
	memoryBackend     = "memory"
)

type Config struct {
	Seed        string
	publicHost  string
	externalURL string
	host        string
	port        uint
	backend     string
	fsRoot      string
}

// Load parses the given arguments list and return a config object (and/or an
// error in case of failures).
func Load(args []string) (Config, error) {
	var cfg Config
	fs := flag.NewFlagSet("gcs-emulator", flag.ContinueOnError)
	fs.StringVar(&cfg.backend, "backend", filesystemBackend, "storage backend (memory or filesystem)")
	fs.StringVar(&cfg.fsRoot, "filesystem-root", "/storage", "filesystem root (required for the filesystem backend). folder will be created if it doesn't exist")
	fs.StringVar(&cfg.publicHost, "public-host", "storage.googleapis.com", "Optional URL for public host")
	fs.StringVar(&cfg.externalURL, "external-url", "", "optional external URL, returned in the Location header for uploads. Defaults to the address where the server is running")
	fs.StringVar(&cfg.host, "host", "0.0.0.0", "host to bind to")
	fs.StringVar(&cfg.Seed, "data", "/data", "where to load data from (provided that the directory exists)")
	fs.UintVar(&cfg.port, "port", 4443, "port to bind to")
	err := fs.Parse(args)
	if err != nil {
		return cfg, err
	}
	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.backend != memoryBackend && c.backend != filesystemBackend {
		return fmt.Errorf(`invalid backend %q, must be either "memory" or "filesystem"`, c.backend)
	}
	if c.backend == filesystemBackend && c.fsRoot == "" {
		return fmt.Errorf("backend %q requires the filesystem-root to be defined", c.backend)
	}
	if c.port > math.MaxUint16 {
		return fmt.Errorf("port %d is too high, maximum value is %d", c.port, math.MaxUint16)
	}
	return nil
}

func (c *Config) ToFakeGcsOptions() fakestorage.Options {
	storageRoot := c.fsRoot
	if c.backend == memoryBackend {
		storageRoot = ""
	}
	return fakestorage.Options{
		StorageRoot: storageRoot,
		Host:        c.host,
		Port:        uint16(c.port),
		PublicHost:  c.publicHost,
		ExternalURL: c.externalURL,
	}
}
