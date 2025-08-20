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
	"io"
	"net/http"
	"strings"
	"time"
)

// IPDetector handles IP address detection
type IPDetector struct {
	client *http.Client
}

// NewIPDetector creates a new IP detector
func NewIPDetector() *IPDetector {
	return &IPDetector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetIPv4 retrieves the current public IPv4 address
func (d *IPDetector) GetIPv4() (string, error) {
	// Try multiple services for reliability
	services := []string{
		"https://ipv4.icanhazip.com",
		"https://api.ipify.org",
		"https://ipv4.ident.me",
		"https://v4.ident.me",
	}
	
	for _, service := range services {
		ip, err := d.getIPFromService(service)
		if err == nil && isValidIPv4(ip) {
			return ip, nil
		}
	}
	
	return "", fmt.Errorf("failed to get IPv4 address from all services")
}

// GetIPv6 retrieves the current public IPv6 address
func (d *IPDetector) GetIPv6() (string, error) {
	// Try multiple services for reliability
	services := []string{
		"https://ipv6.icanhazip.com",
		"https://api6.ipify.org",
		"https://ipv6.ident.me",
		"https://v6.ident.me",
	}
	
	for _, service := range services {
		ip, err := d.getIPFromService(service)
		if err == nil && isValidIPv6(ip) {
			return ip, nil
		}
	}
	
	return "", fmt.Errorf("failed to get IPv6 address from all services")
}

// getIPFromService fetches IP from a specific service
func (d *IPDetector) getIPFromService(url string) (string, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	ip := strings.TrimSpace(string(body))
	return ip, nil
}

// isValidIPv4 checks if the string is a valid IPv4 address
func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		
		num := 0
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
			num = num*10 + int(char-'0')
		}
		
		if num > 255 {
			return false
		}
	}
	
	return true
}

// isValidIPv6 checks if the string is a valid IPv6 address
func isValidIPv6(ip string) bool {
	// Basic IPv6 validation - contains colons and valid hex characters
	if !strings.Contains(ip, ":") {
		return false
	}
	
	// Split by colons
	parts := strings.Split(ip, ":")
	if len(parts) < 3 || len(parts) > 8 {
		return false
	}
	
	// Check for valid hex characters
	for _, part := range parts {
		if len(part) > 4 {
			return false
		}
		
		for _, char := range part {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}
	
	return true
}