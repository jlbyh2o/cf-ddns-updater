/*
Cloudflare Dynamic DNS Updater
Copyright (C) 2025

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	Version = "0.3.0"
	AppName = "Cloudflare DDNS Updater"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "cf-ddns.conf", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	logFile := flag.String("log", "", "Log file path (optional, logs to stdout if not specified)")
	runOnce := flag.Bool("once", false, "Run once and exit (ignore interval setting)")
	versionFlag := flag.Bool("version", false, "Show version information and exit")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("%s v%s\n", AppName, Version)
		os.Exit(0)
	}

	fmt.Printf("%s v%s\n", AppName, Version)

	// Load configuration first to get log file setting
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Determine log file: command line flag takes precedence over config file
	logFileToUse := *logFile
	if logFileToUse == "" && config.LogFile != "" {
		logFileToUse = config.LogFile
	}

	// Setup logging
	if err := setupLogging(*verbose, logFileToUse); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	// Initialize updater
	updater := NewDDNSUpdater(config, *verbose)

	// Run update loop
	if *runOnce || config.Interval <= 0 {
		// Run once
		if err := updater.Update(); err != nil {
			log.Fatalf("Failed to update DNS records: %v", err)
		}
		log.Println("DNS records updated successfully")
	} else {
		// Run continuously with interval
		log.Printf("Starting continuous mode with %d second interval", config.Interval)
		for {
			if err := updater.Update(); err != nil {
				log.Printf("Failed to update DNS records: %v", err)
			} else {
				log.Println("DNS records updated successfully")
			}

			log.Printf("Waiting %d seconds before next update...", config.Interval)
			time.Sleep(time.Duration(config.Interval) * time.Second)
		}
	}
}

func loadConfig(filename string) (*Config, error) {
	// Try to find config file in multiple locations
	configPath, err := findConfigFile(filename)
	if err != nil {
		return nil, err
	}

	// Load configuration using TOML parser
	config, err := LoadConfigFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	log.Printf("Loaded configuration from: %s", configPath)
	return config, nil
}

// findConfigFile searches for the config file in multiple locations
func findConfigFile(filename string) (string, error) {
	// If filename is an absolute path, use it directly
	if filepath.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}
		return "", fmt.Errorf("config file not found: %s", filename)
	}

	// Generate alternative filenames (.conf versions)
	filenames := []string{filename}
	if !strings.HasSuffix(strings.ToLower(filename), ".conf") {
		// If no .conf extension specified, try adding .conf
		filenames = append(filenames, filename+".conf")
	}

	// Search paths in order of preference
	var searchPaths []string

	// On Linux, check system config directory first
	if runtime.GOOS == "linux" {
		for _, fname := range filenames {
			searchPaths = append(searchPaths, filepath.Join("/etc/cf-ddns", fname))
		}
	}

	// Check current working directory
	searchPaths = append(searchPaths, filenames...)

	// Check relative to executable
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		for _, fname := range filenames {
			searchPaths = append(searchPaths, filepath.Join(execDir, fname))
		}
	}

	// Try each path
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("config file not found: %s", filename)
}

// setupLogging configures logging based on the provided options
func setupLogging(verbose bool, logFile string) error {
	// Set log flags
	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Setup log output
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		log.SetOutput(file)
		log.Printf("Logging to file: %s", logFile)
	}

	return nil
}
