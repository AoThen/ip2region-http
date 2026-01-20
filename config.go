// Copyright 2025 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ServerPort   int
	DBUrl        string
	DBPath       string
	DownloadMode bool
}

func LoadConfig() (*Config, error) {
	config := &Config{
		ServerPort:   8080,
		DBPath:       "data/ip2region_n.xdb",
		DownloadMode: true,
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		if n, err := fmt.Sscanf(port, "%d", &config.ServerPort); n != 1 || err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %s", port)
		}
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		config.DBPath = dbPath
	}

	if dbUrl := os.Getenv("DB_URL"); dbUrl != "" {
		config.DBUrl = dbUrl
	} else {
		config.DBUrl = "https://github.com/hel2o/ip2region_update/releases/download/250820/ip2region_n.xdb"
	}

	if downloadMode := os.Getenv("DOWNLOAD_MODE"); downloadMode != "" {
		config.DownloadMode = strings.ToLower(downloadMode) == "true"
	}

	return config, nil
}
