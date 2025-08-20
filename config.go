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
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	// Cloudflare API credentials
	Cloudflare CloudflareConfig `toml:"cloudflare"`
	
	// Domains to update
	Domains []DomainConfig `toml:"domains"`
	
	// Update interval in seconds (0 = run once)
	Interval int `toml:"interval,omitempty"`
	
	// Logging configuration
	Verbose bool `toml:"verbose,omitempty"`
	
	// Optional log file path
	LogFile string `toml:"log_file,omitempty"`
}

// CloudflareConfig contains Cloudflare API settings
type CloudflareConfig struct {
	// API Token (recommended) or API Key + Email
	APIToken string `toml:"api_token,omitempty"`
	APIKey   string `toml:"api_key,omitempty"`
	Email    string `toml:"email,omitempty"`
	
	// Zone ID (optional, will be auto-detected if not provided)
	ZoneID string `toml:"zone_id,omitempty"`
}

// DomainConfig represents a domain to update
type DomainConfig struct {
	// Domain name (e.g., "example.com" or "subdomain.example.com")
	Name string `toml:"name"`
	
	// Record types to update: "A", "AAAA", or "both"
	RecordTypes string `toml:"record_types"`
	
	// TTL for DNS records (default: 300)
	TTL int `toml:"ttl,omitempty"`
	
	// Proxied through Cloudflare (default: false)
	Proxied bool `toml:"proxied,omitempty"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate Cloudflare credentials
	if c.Cloudflare.APIToken == "" && (c.Cloudflare.APIKey == "" || c.Cloudflare.Email == "") {
		return fmt.Errorf("either api_token or both api_key and email must be provided")
	}
	
	// Validate domains
	if len(c.Domains) == 0 {
		return fmt.Errorf("at least one domain must be configured")
	}
	
	for i, domain := range c.Domains {
		if domain.Name == "" {
			return fmt.Errorf("domain[%d]: name is required", i)
		}
		
		// Validate record types
		recordTypes := strings.ToLower(domain.RecordTypes)
		if recordTypes == "" {
			recordTypes = "both" // default
			c.Domains[i].RecordTypes = "both"
		}
		
		if recordTypes != "a" && recordTypes != "aaaa" && recordTypes != "both" {
			return fmt.Errorf("domain[%d]: record_types must be 'A', 'AAAA', or 'both'", i)
		}
		
		// Set default TTL
		if domain.TTL == 0 {
			c.Domains[i].TTL = 300
		}
	}
	
	return nil
}

// ShouldUpdateA returns true if A records should be updated for this domain
func (d *DomainConfig) ShouldUpdateA() bool {
	recordTypes := strings.ToLower(d.RecordTypes)
	return recordTypes == "a" || recordTypes == "both"
}

// ShouldUpdateAAAA returns true if AAAA records should be updated for this domain
func (d *DomainConfig) ShouldUpdateAAAA() bool {
	recordTypes := strings.ToLower(d.RecordTypes)
	return recordTypes == "aaaa" || recordTypes == "both"
}

// LoadConfigFromFile loads configuration from TOML file
func LoadConfigFromFile(filename string) (*Config, error) {
	var config Config
	
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	
	return &config, config.Validate()
}