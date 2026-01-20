// Copyright 2025 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"ip2region-http/xdb"
)

var (
	searcher *xdb.Searcher
	dbSize   int64
	mu       sync.RWMutex
)

type SearchResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type RegionData struct {
	IP       string `json:"ip"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	District string `json:"district,omitempty"`
	ISP      string `json:"isp,omitempty"`
}

type HealthData struct {
	Status   string `json:"status"`
	DBLoaded bool   `json:"db_loaded"`
	DBPath   string `json:"db_path"`
	DBSize   int64  `json:"db_size,omitempty"`
}

func parseRegion(region string) RegionData {
	parts := strings.Split(region, "|")
	data := RegionData{IP: ""}
	if len(parts) > 0 {
		data.Country = parts[0]
	}
	if len(parts) > 1 {
		data.Province = parts[1]
	}
	if len(parts) > 2 {
		data.City = parts[2]
	}
	if len(parts) > 3 {
		data.ISP = parts[3]
	}
	return data
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	format := r.URL.Query().Get("format")

	if ip == "" {
		writeJSON(w, SearchResponse{Code: 400, Msg: "missing ip parameter"})
		return
	}

	mu.RLock()
	s := searcher
	mu.RUnlock()

	if s == nil {
		writeJSON(w, SearchResponse{Code: 500, Msg: "database not loaded"})
		return
	}

	region, err := s.SearchByStr(ip)
	if err != nil {
		writeJSON(w, SearchResponse{Code: 404, Msg: err.Error()})
		return
	}

	switch format {
	case "text":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(region))
	case "fields":
		data := parseRegion(region)
		data.IP = ip
		writeJSON(w, SearchResponse{Code: 0, Msg: "success", Data: data})
	default:
		data := RegionData{
			IP:     ip,
			Region: region,
		}
		writeJSON(w, SearchResponse{Code: 0, Msg: "success", Data: data})
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	s := searcher
	mu.RUnlock()

	status := "healthy"
	if s == nil {
		status = "unhealthy"
	}

	data := HealthData{
		Status:   status,
		DBLoaded: s != nil,
		DBPath:   os.Getenv("DB_PATH"),
		DBSize:   dbSize,
	}

	writeJSON(w, SearchResponse{Code: 0, Msg: "ok", Data: data})
}

func writeJSON(w http.ResponseWriter, data SearchResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func initDatabase(config *Config) error {
	mu.Lock()
	defer mu.Unlock()

	if !DatabaseExists(config.DBPath) {
		if config.DownloadMode {
			log.Printf("Downloading database from %s...", config.DBUrl)
			if err := DownloadDatabase(config.DBUrl, config.DBPath); err != nil {
				log.Printf("Download failed: %v", err)
				return fmt.Errorf("download database failed: %w", err)
			}
			log.Println("Database downloaded successfully")
		} else {
			return fmt.Errorf("database file not found: %s", config.DBPath)
		}
	}

	fileInfo, err := os.Stat(config.DBPath)
	if err != nil {
		return fmt.Errorf("stat database: %w", err)
	}
	dbSize = fileInfo.Size()

	content, err := xdb.LoadContentFromFile(config.DBPath)
	if err != nil {
		return fmt.Errorf("load database: %w", err)
	}

	searcher, err = xdb.NewWithBuffer(xdb.IPv4, content)
	if err != nil {
		return fmt.Errorf("create searcher: %w", err)
	}
	log.Printf("Database loaded successfully: %s (%d bytes)", config.DBPath, dbSize)

	return nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed: %v", err)
	}

	if err := initDatabase(config); err != nil {
		log.Printf("Warning: %v", err)
	}

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf(":%d", config.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
